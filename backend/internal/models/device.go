package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	DeviceStatusOnline  = "online"
	DeviceStatusOffline = "offline"
)

type Device struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID       uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	Name         string     `gorm:"not null" json:"name"`
	DeviceUID    string     `gorm:"uniqueIndex;not null" json:"device_uid"`
	FCMToken     string     `gorm:"type:text" json:"-"`
	Status       string     `gorm:"default:offline" json:"status"`
	BatteryLevel int        `gorm:"default:0" json:"battery_level"`
	AppVersion   string     `json:"app_version"`
	LastSeenAt   *time.Time `json:"last_seen_at"`
	CreatedAt    time.Time  `json:"created_at"`

	// Relations
	User            User             `gorm:"foreignKey:UserID" json:"-"`
	ForwardingRules []ForwardingRule `gorm:"foreignKey:DeviceID" json:"forwarding_rules,omitempty"`
}

func (Device) TableName() string {
	return "devices"
}
