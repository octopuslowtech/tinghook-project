package models

import (
	"time"

	"github.com/google/uuid"
)

type ForwardingRule struct {
	ID            uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	DeviceID      *uuid.UUID `gorm:"type:uuid;index" json:"device_id,omitempty"`
	TriggerType   string     `gorm:"not null" json:"trigger_type"`
	SenderFilter  string     `json:"sender_filter"`
	ContentFilter string     `gorm:"type:text" json:"content_filter"`
	WebhookURL    string     `gorm:"type:text;not null" json:"webhook_url"`
	SecretHeader  string     `gorm:"type:text" json:"-"`
	Method        string     `gorm:"default:POST" json:"method"`
	IsActive      bool       `gorm:"default:true" json:"is_active"`
	CreatedAt     time.Time  `json:"created_at"`

	// Relations
	User   User    `gorm:"foreignKey:UserID" json:"-"`
	Device *Device `gorm:"foreignKey:DeviceID" json:"-"`
}

func (ForwardingRule) TableName() string {
	return "forwarding_rules"
}
