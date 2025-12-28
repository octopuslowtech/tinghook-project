package models

import (
	"time"

	"github.com/google/uuid"
)

type MessageDirection string
type MessageStatus string

const (
	DirectionInbound  MessageDirection = "inbound"
	DirectionOutbound MessageDirection = "outbound"

	StatusPending   MessageStatus = "pending"
	StatusSent      MessageStatus = "sent"
	StatusDelivered MessageStatus = "delivered"
	StatusFailed    MessageStatus = "failed"
)

type MessageLog struct {
	ID           uint             `gorm:"primaryKey" json:"id"`
	UserID       uuid.UUID        `gorm:"type:uuid;not null;index" json:"user_id"`
	DeviceID     *uuid.UUID       `gorm:"type:uuid;index" json:"device_id,omitempty"`
	Direction    MessageDirection `gorm:"not null" json:"direction"`
	SimSlot      int              `gorm:"default:0" json:"sim_slot"`
	Sender       string           `gorm:"size:50" json:"sender"`
	Receiver     string           `gorm:"size:50" json:"receiver"`
	Content      string           `gorm:"type:text" json:"content"`
	Status       MessageStatus    `gorm:"default:pending" json:"status"`
	ErrorMessage string           `gorm:"type:text" json:"error_message,omitempty"`
	RetryCount   int              `gorm:"default:0" json:"retry_count"`
	CreatedAt    time.Time        `gorm:"index" json:"created_at"`
	ProcessedAt  *time.Time       `json:"processed_at,omitempty"`

	// Relations
	User   User    `gorm:"foreignKey:UserID" json:"-"`
	Device *Device `gorm:"foreignKey:DeviceID" json:"-"`
}

func (MessageLog) TableName() string {
	return "message_logs"
}
