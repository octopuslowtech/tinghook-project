package handlers

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/octopuslowtech/tinghook-project/backend/internal/services"
	ws "github.com/octopuslowtech/tinghook-project/backend/internal/websockets"
)

const (
	authTimeout = 30 * time.Second
)

type WSHandler struct {
	hub           *ws.Hub
	userService   services.UserService
	deviceService services.DeviceService
	logService    services.LogService
	ruleService   services.RuleService
	deviceHandler *ws.DeviceHandler
}

func NewWSHandler(
	hub *ws.Hub,
	userService services.UserService,
	deviceService services.DeviceService,
	logService services.LogService,
	ruleService services.RuleService,
) *WSHandler {
	deviceHandler := ws.NewDeviceHandler(hub, userService, deviceService, logService, ruleService)
	return &WSHandler{
		hub:           hub,
		userService:   userService,
		deviceService: deviceService,
		logService:    logService,
		ruleService:   ruleService,
		deviceHandler: deviceHandler,
	}
}

func (h *WSHandler) HandleDeviceWS(c *websocket.Conn) {
	defer c.Close()

	authChan := make(chan *ws.DeviceConnection)
	errChan := make(chan error)

	go func() {
		conn, err := h.waitForAuth(c)
		if err != nil {
			errChan <- err
			return
		}
		authChan <- conn
	}()

	select {
	case conn := <-authChan:
		h.hub.RegisterDevice(conn)

		go conn.WritePump()
		conn.ReadPump(h.deviceHandler.HandleMessage)

		h.handleDisconnect(conn)

	case err := <-errChan:
		log.Printf("auth failed: %v", err)
		h.sendAuthFail(c, err.Error())

	case <-time.After(authTimeout):
		log.Printf("auth timeout")
		h.sendAuthFail(c, "authentication timeout")
	}
}

func (h *WSHandler) waitForAuth(c *websocket.Conn) (*ws.DeviceConnection, error) {
	_, msgBytes, err := c.ReadMessage()
	if err != nil {
		return nil, err
	}

	var msg ws.Message
	if err := json.Unmarshal(msgBytes, &msg); err != nil {
		return nil, err
	}

	if msg.Type != ws.MsgTypeAuth {
		return nil, fiber.NewError(fiber.StatusBadRequest, "expected AUTH message")
	}

	var authData ws.AuthData
	if err := msg.UnmarshalData(&authData); err != nil {
		return nil, err
	}

	user, err := h.userService.GetByAPIKey(authData.APIKey)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "invalid api key")
	}

	device, err := h.deviceService.GetByDeviceUID(authData.DeviceUID)
	if err != nil {
		device, err = h.deviceService.Register(user.ID, "Unknown Device", authData.DeviceUID)
		if err != nil {
			return nil, err
		}
	}

	if device.UserID != user.ID {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "device belongs to another user")
	}

	conn, err := ws.NewDeviceConnection(device.ID.String(), user.ID.String(), c, h.hub)
	if err != nil {
		return nil, err
	}

	if err := h.deviceService.SetOnline(device.ID, 0); err != nil {
		log.Printf("failed to set device online: %v", err)
	}

	authOKMsg, err := ws.NewMessage(ws.MsgTypeAuthOK, &ws.AuthOKData{
		DeviceID: device.ID.String(),
	})
	if err != nil {
		return nil, err
	}

	msgBytes, err = json.Marshal(authOKMsg)
	if err != nil {
		return nil, err
	}

	if err := c.WriteMessage(websocket.TextMessage, msgBytes); err != nil {
		return nil, err
	}

	log.Printf("device %s authenticated for user %s", device.ID, user.ID)
	return conn, nil
}

func (h *WSHandler) handleDisconnect(conn *ws.DeviceConnection) {
	if err := h.deviceService.SetOffline(conn.DeviceID); err != nil {
		log.Printf("failed to set device offline: %v", err)
	}
	log.Printf("device %s disconnected", conn.DeviceID)
}

func (h *WSHandler) sendAuthFail(c *websocket.Conn, errorMsg string) {
	msg, err := ws.NewMessage(ws.MsgTypeAuthFail, &ws.AuthFailData{
		Error: errorMsg,
	})
	if err != nil {
		return
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return
	}

	c.WriteMessage(websocket.TextMessage, msgBytes)
}
