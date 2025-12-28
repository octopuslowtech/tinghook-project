package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID               uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Email            string    `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash     string    `gorm:"not null" json:"-"`
	APIKey           string    `gorm:"uniqueIndex;not null;size:64" json:"-"`
	SubscriptionPlan string    `gorm:"default:free" json:"subscription_plan"`
	Credits          int       `gorm:"default:0" json:"credits"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`

	// Relations
	Devices         []Device         `gorm:"foreignKey:UserID" json:"devices,omitempty"`
	ForwardingRules []ForwardingRule `gorm:"foreignKey:UserID" json:"forwarding_rules,omitempty"`
}

func (User) TableName() string {
	return "users"
}
