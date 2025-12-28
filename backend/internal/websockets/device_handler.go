package websockets

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/octopuslowtech/tinghook-project/backend/internal/models"
	"github.com/octopuslowtech/tinghook-project/backend/internal/services"
)

type DeviceHandler struct {
	hub           *Hub
	userService   services.UserService
	deviceService services.DeviceService
	logService    services.LogService
	ruleService   services.RuleService
}

func NewDeviceHandler(
	hub *Hub,
	userService services.UserService,
	deviceService services.DeviceService,
	logService services.LogService,
	ruleService services.RuleService,
) *DeviceHandler {
	return &DeviceHandler{
		hub:           hub,
		userService:   userService,
		deviceService: deviceService,
		logService:    logService,
		ruleService:   ruleService,
	}
}

func (h *DeviceHandler) HandleMessage(conn *DeviceConnection, msg *Message) {
	switch msg.Type {
	case MsgTypePing:
		h.handlePing(conn, msg)
	case MsgTypeSMSReceived:
		h.handleSMSReceived(conn, msg)
	case MsgTypeNotifReceived:
		h.handleNotificationReceived(conn, msg)
	case MsgTypeSMSSent:
		h.handleSMSSent(conn, msg)
	case MsgTypeSMSFailed:
		h.handleSMSFailed(conn, msg)
	default:
		log.Printf("unknown message type: %s", msg.Type)
	}
}

func (h *DeviceHandler) handlePing(conn *DeviceConnection, msg *Message) {
	var data PingData
	if err := msg.UnmarshalData(&data); err != nil {
		log.Printf("failed to unmarshal ping data: %v", err)
		return
	}

	go func() {
		if err := h.deviceService.SetOnline(conn.DeviceID, data.Battery); err != nil {
			log.Printf("failed to update device status: %v", err)
		}
	}()

	pongMsg, err := NewMessage(MsgTypePong, &PongData{
		Timestamp: time.Now(),
	})
	if err != nil {
		log.Printf("failed to create pong message: %v", err)
		return
	}

	if err := conn.SendMessage(pongMsg); err != nil {
		log.Printf("failed to send pong: %v", err)
	}
}

func (h *DeviceHandler) handleSMSReceived(conn *DeviceConnection, msg *Message) {
	var data SMSReceivedData
	if err := msg.UnmarshalData(&data); err != nil {
		log.Printf("failed to unmarshal sms received data: %v", err)
		return
	}

	go func() {
		msgLog, err := h.logService.Create(
			conn.UserID,
			&conn.DeviceID,
			models.DirectionInbound,
			data.Sender,
			"",
			data.Content,
			data.SimSlot,
		)
		if err != nil {
			log.Printf("failed to create message log: %v", err)
			return
		}

		h.matchAndDispatch(conn.DeviceID, "sms", data.Sender, data.Content, msgLog.ID)
	}()
}

func (h *DeviceHandler) handleNotificationReceived(conn *DeviceConnection, msg *Message) {
	var data NotificationReceivedData
	if err := msg.UnmarshalData(&data); err != nil {
		log.Printf("failed to unmarshal notification data: %v", err)
		return
	}

	go func() {
		content := data.Title + "\n" + data.Content
		h.matchAndDispatch(conn.DeviceID, "notification", data.PackageName, content, 0)
	}()
}

func (h *DeviceHandler) handleSMSSent(conn *DeviceConnection, msg *Message) {
	var data SMSSentData
	if err := msg.UnmarshalData(&data); err != nil {
		log.Printf("failed to unmarshal sms sent data: %v", err)
		return
	}

	go func() {
		logID, err := parseLogID(data.RequestID)
		if err != nil {
			log.Printf("failed to parse request id: %v", err)
			return
		}

		if err := h.logService.UpdateStatus(logID, models.StatusSent, ""); err != nil {
			log.Printf("failed to update log status: %v", err)
		}
	}()
}

func (h *DeviceHandler) handleSMSFailed(conn *DeviceConnection, msg *Message) {
	var data SMSFailedData
	if err := msg.UnmarshalData(&data); err != nil {
		log.Printf("failed to unmarshal sms failed data: %v", err)
		return
	}

	go func() {
		logID, err := parseLogID(data.RequestID)
		if err != nil {
			log.Printf("failed to parse request id: %v", err)
			return
		}

		if err := h.logService.UpdateStatus(logID, models.StatusFailed, data.Error); err != nil {
			log.Printf("failed to update log status: %v", err)
		}
	}()
}

func (h *DeviceHandler) matchAndDispatch(deviceID uuid.UUID, triggerType, sender, content string, logID uint) {
	rules, err := h.ruleService.MatchRules(deviceID, triggerType, sender, content)
	if err != nil {
		log.Printf("failed to match rules: %v", err)
		return
	}

	for _, rule := range rules {
		log.Printf("matched rule %d, dispatching webhook to %s", rule.ID, rule.WebhookURL)
		// TODO: Dispatch via Asynq worker when webhook dispatcher is implemented
		// h.dispatcher.Dispatch(rule, sender, content, logID)
	}
}

func parseLogID(requestID string) (uint, error) {
	var id uint
	_, err := fmt.Sscanf(requestID, "%d", &id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
