package websockets

import (
	"encoding/json"
	"time"
)

const (
	MsgTypeAuth          = "AUTH"
	MsgTypeAuthOK        = "AUTH_OK"
	MsgTypeAuthFail      = "AUTH_FAIL"
	MsgTypePing          = "PING"
	MsgTypePong          = "PONG"
	MsgTypeSMSReceived   = "SMS_RECEIVED"
	MsgTypeNotifReceived = "NOTIFICATION_RECEIVED"
	MsgTypeSendSMS       = "SEND_SMS"
	MsgTypeSMSSent       = "SMS_SENT"
	MsgTypeSMSFailed     = "SMS_FAILED"
)

type Message struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type AuthData struct {
	APIKey    string `json:"api_key"`
	DeviceUID string `json:"device_uid"`
}

type AuthOKData struct {
	DeviceID string `json:"device_id"`
}

type AuthFailData struct {
	Error string `json:"error"`
}

type PingData struct {
	Battery int `json:"battery"`
	Signal  int `json:"signal"`
}

type PongData struct {
	Timestamp time.Time `json:"timestamp"`
}

type SMSReceivedData struct {
	Sender    string    `json:"sender"`
	Content   string    `json:"content"`
	SimSlot   int       `json:"sim_slot"`
	Timestamp time.Time `json:"timestamp"`
}

type NotificationReceivedData struct {
	PackageName string    `json:"package_name"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Timestamp   time.Time `json:"timestamp"`
}

type SendSMSData struct {
	RequestID string `json:"request_id"`
	Phone     string `json:"phone"`
	Content   string `json:"content"`
	SimSlot   int    `json:"sim_slot"`
}

type SMSSentData struct {
	RequestID string `json:"request_id"`
}

type SMSFailedData struct {
	RequestID string `json:"request_id"`
	Error     string `json:"error"`
}

func NewMessage(msgType string, data interface{}) (*Message, error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return &Message{
		Type: msgType,
		Data: dataBytes,
	}, nil
}

func (m *Message) UnmarshalData(v interface{}) error {
	return json.Unmarshal(m.Data, v)
}
