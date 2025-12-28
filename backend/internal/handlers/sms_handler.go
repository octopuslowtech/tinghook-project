package handlers

import (
	"regexp"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/octopuslowtech/tinghook-project/backend/internal/handlers/dto"
	"github.com/octopuslowtech/tinghook-project/backend/internal/models"
	"github.com/octopuslowtech/tinghook-project/backend/internal/services"
	"github.com/octopuslowtech/tinghook-project/backend/internal/websockets"
)

const (
	maxSMSContentLength = 1600
	freePlanSignature   = "\n\n- Sent via TingHook"
)

type SMSHandler struct {
	hub           *websockets.Hub
	userService   services.UserService
	deviceService services.DeviceService
	logService    services.LogService
}

func NewSMSHandler(
	hub *websockets.Hub,
	userService services.UserService,
	deviceService services.DeviceService,
	logService services.LogService,
) *SMSHandler {
	return &SMSHandler{
		hub:           hub,
		userService:   userService,
		deviceService: deviceService,
		logService:    logService,
	}
}

func (h *SMSHandler) SendSMS(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	var req dto.SendSMSRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "invalid request body"})
	}

	if req.Phone == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "phone is required"})
	}

	if req.Content == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "content is required"})
	}

	if len(req.Content) > maxSMSContentLength {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "content exceeds maximum length of 1600 characters"})
	}

	if !isValidPhone(req.Phone) {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "invalid phone number format"})
	}

	var targetDeviceID uuid.UUID
	var err error

	if req.DeviceID != nil && *req.DeviceID != "" {
		targetDeviceID, err = uuid.Parse(*req.DeviceID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "invalid device_id format"})
		}

		device, err := h.deviceService.GetByID(targetDeviceID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{Error: "device not found"})
		}

		if device.UserID != user.ID {
			return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{Error: "device not found"})
		}

		if !h.hub.GetDeviceStatus(targetDeviceID) {
			return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{Error: "device is offline"})
		}
	} else {
		onlineDevices := h.hub.GetOnlineDevices(user.ID)
		if len(onlineDevices) == 0 {
			return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{Error: "no online device available"})
		}
		targetDeviceID = onlineDevices[0]
	}

	content := req.Content
	if user.SubscriptionPlan == "free" {
		content = req.Content + freePlanSignature
	}

	requestID := uuid.New()

	log, err := h.logService.Create(
		user.ID,
		&targetDeviceID,
		models.DirectionOutbound,
		"",
		req.Phone,
		content,
		req.SimSlot,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Error: "failed to create message log"})
	}

	sendMsg, err := websockets.NewMessage(websockets.MsgTypeSendSMS, websockets.SendSMSData{
		RequestID: requestID.String(),
		Phone:     req.Phone,
		Content:   content,
		SimSlot:   req.SimSlot,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Error: "failed to create message"})
	}

	if err := h.hub.SendToDevice(targetDeviceID, sendMsg); err != nil {
		_ = h.logService.UpdateStatus(log.ID, models.StatusFailed, "failed to send to device")
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Error: "failed to send SMS to device"})
	}

	return c.Status(fiber.StatusAccepted).JSON(dto.SendSMSResponse{
		RequestID: requestID.String(),
		Status:    "queued",
		DeviceID:  targetDeviceID.String(),
	})
}

func (h *SMSHandler) GetDevicesStatus(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	devices, err := h.deviceService.ListByUser(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Error: "failed to get devices"})
	}

	deviceDTOs := make([]dto.DeviceStatusDTO, 0, len(devices))
	for _, device := range devices {
		status := device.Status
		if h.hub.GetDeviceStatus(device.ID) {
			status = models.DeviceStatusOnline
		} else {
			status = models.DeviceStatusOffline
		}

		var lastSeenAt *string
		if device.LastSeenAt != nil {
			t := device.LastSeenAt.Format("2006-01-02T15:04:05Z07:00")
			lastSeenAt = &t
		}

		deviceDTOs = append(deviceDTOs, dto.DeviceStatusDTO{
			ID:         device.ID.String(),
			Name:       device.Name,
			Status:     status,
			Battery:    device.BatteryLevel,
			LastSeenAt: lastSeenAt,
		})
	}

	return c.JSON(dto.DevicesStatusResponse{
		Devices: deviceDTOs,
	})
}

func isValidPhone(phone string) bool {
	phoneRegex := regexp.MustCompile(`^\+?[0-9]{7,15}$`)
	return phoneRegex.MatchString(phone)
}
