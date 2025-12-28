package websockets

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = 30 * time.Second
	maxMessageSize = 8192
)

type MessageHandler func(conn *DeviceConnection, msg *Message)

func (c *DeviceConnection) ReadPump(handler MessageHandler) {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, msgBytes, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("websocket error: %v", err)
			}
			break
		}

		var msg Message
		if err := json.Unmarshal(msgBytes, &msg); err != nil {
			log.Printf("failed to unmarshal message: %v", err)
			continue
		}

		if handler != nil {
			handler(c, &msg)
		}
	}
}

func (c *DeviceConnection) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *DeviceConnection) Close() {
	c.Hub.unregister <- c
}

func (c *DeviceConnection) SendMessage(msg *Message) error {
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	select {
	case c.Send <- msgBytes:
		return nil
	default:
		return ErrDeviceNotConnected
	}
}

func NewDeviceConnection(deviceID, userID string, conn *websocket.Conn, hub *Hub) (*DeviceConnection, error) {
	did, err := parseUUID(deviceID)
	if err != nil {
		return nil, err
	}
	uid, err := parseUUID(userID)
	if err != nil {
		return nil, err
	}

	return &DeviceConnection{
		DeviceID: did,
		UserID:   uid,
		Conn:     conn,
		Send:     make(chan []byte, 256),
		Hub:      hub,
	}, nil
}

func parseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}
