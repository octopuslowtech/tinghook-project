package dto

import (
	"github.com/octopuslowtech/tinghook-project/backend/internal/models"
)

type CreateRuleRequest struct {
	DeviceID      *string `json:"device_id"`
	TriggerType   string  `json:"trigger_type" validate:"required,oneof=sms notification"`
	SenderFilter  string  `json:"sender_filter"`
	ContentFilter string  `json:"content_filter"`
	WebhookURL    string  `json:"webhook_url" validate:"required,url"`
	SecretHeader  string  `json:"secret_header"`
	Method        string  `json:"method" validate:"omitempty,oneof=GET POST PUT"`
}

type UpdateRuleRequest struct {
	DeviceID      *string `json:"device_id"`
	TriggerType   *string `json:"trigger_type" validate:"omitempty,oneof=sms notification"`
	SenderFilter  *string `json:"sender_filter"`
	ContentFilter *string `json:"content_filter"`
	WebhookURL    *string `json:"webhook_url" validate:"omitempty,url"`
	SecretHeader  *string `json:"secret_header"`
	Method        *string `json:"method" validate:"omitempty,oneof=GET POST PUT"`
	IsActive      *bool   `json:"is_active"`
}

type RuleDTO struct {
	ID            uint    `json:"id"`
	DeviceID      *string `json:"device_id,omitempty"`
	TriggerType   string  `json:"trigger_type"`
	SenderFilter  string  `json:"sender_filter"`
	ContentFilter string  `json:"content_filter"`
	WebhookURL    string  `json:"webhook_url"`
	Method        string  `json:"method"`
	IsActive      bool    `json:"is_active"`
	CreatedAt     string  `json:"created_at"`
}

type WebhookTestResult struct {
	Success      bool   `json:"success"`
	StatusCode   int    `json:"status_code"`
	ResponseTime int64  `json:"response_time_ms"`
	Error        string `json:"error,omitempty"`
}

func ToRuleDTO(rule *models.ForwardingRule) *RuleDTO {
	dto := &RuleDTO{
		ID:            rule.ID,
		TriggerType:   rule.TriggerType,
		SenderFilter:  rule.SenderFilter,
		ContentFilter: rule.ContentFilter,
		WebhookURL:    rule.WebhookURL,
		Method:        rule.Method,
		IsActive:      rule.IsActive,
		CreatedAt:     rule.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if rule.DeviceID != nil {
		deviceIDStr := rule.DeviceID.String()
		dto.DeviceID = &deviceIDStr
	}

	return dto
}

func ToRuleDTOList(rules []models.ForwardingRule) []RuleDTO {
	dtos := make([]RuleDTO, len(rules))
	for i, rule := range rules {
		dto := ToRuleDTO(&rule)
		dtos[i] = *dto
	}
	return dtos
}
