package websockets

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
)

var (
	ErrDeviceNotConnected = errors.New("device not connected")
	ErrDeviceNotFound     = errors.New("device not found")
)

type Hub struct {
	devices    map[uuid.UUID]*DeviceConnection
	register   chan *DeviceConnection
	unregister chan *DeviceConnection
	broadcast  chan *Message
	mu         sync.RWMutex
}

type DeviceConnection struct {
	DeviceID uuid.UUID
	UserID   uuid.UUID
	Conn     *websocket.Conn
	Send     chan []byte
	Hub      *Hub
}

func NewHub() *Hub {
	return &Hub{
		devices:    make(map[uuid.UUID]*DeviceConnection),
		register:   make(chan *DeviceConnection),
		unregister: make(chan *DeviceConnection),
		broadcast:  make(chan *Message),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case conn := <-h.register:
			h.mu.Lock()
			h.devices[conn.DeviceID] = conn
			h.mu.Unlock()

		case conn := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.devices[conn.DeviceID]; ok {
				delete(h.devices, conn.DeviceID)
				close(conn.Send)
			}
			h.mu.Unlock()

		case msg := <-h.broadcast:
			h.mu.RLock()
			msgBytes, err := json.Marshal(msg)
			if err != nil {
				h.mu.RUnlock()
				continue
			}
			for _, conn := range h.devices {
				select {
				case conn.Send <- msgBytes:
				default:
					h.mu.RUnlock()
					h.mu.Lock()
					delete(h.devices, conn.DeviceID)
					close(conn.Send)
					h.mu.Unlock()
					h.mu.RLock()
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) RegisterDevice(conn *DeviceConnection) {
	h.register <- conn
}

func (h *Hub) UnregisterDevice(deviceID uuid.UUID) {
	h.mu.RLock()
	conn, ok := h.devices[deviceID]
	h.mu.RUnlock()
	if ok {
		h.unregister <- conn
	}
}

func (h *Hub) SendToDevice(deviceID uuid.UUID, msg *Message) error {
	h.mu.RLock()
	conn, ok := h.devices[deviceID]
	h.mu.RUnlock()

	if !ok {
		return ErrDeviceNotConnected
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	select {
	case conn.Send <- msgBytes:
		return nil
	default:
		return ErrDeviceNotConnected
	}
}

func (h *Hub) GetDeviceStatus(deviceID uuid.UUID) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, ok := h.devices[deviceID]
	return ok
}

func (h *Hub) GetOnlineDevices(userID uuid.UUID) []uuid.UUID {
	h.mu.RLock()
	defer h.mu.RUnlock()

	var devices []uuid.UUID
	for _, conn := range h.devices {
		if conn.UserID == userID {
			devices = append(devices, conn.DeviceID)
		}
	}
	return devices
}

func (h *Hub) GetConnection(deviceID uuid.UUID) (*DeviceConnection, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	conn, ok := h.devices[deviceID]
	return conn, ok
}

func (h *Hub) Broadcast(msg *Message) {
	h.broadcast <- msg
}
