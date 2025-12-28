package dto

import (
	"time"

	"github.com/octopuslowtech/tinghook-project/backend/internal/models"
)

type LogQueryParams struct {
	Page      int       `query:"page"`
	Limit     int       `query:"limit"`
	Direction string    `query:"direction"`
	Status    string    `query:"status"`
	DeviceID  string    `query:"device_id"`
	From      time.Time `query:"from"`
	To        time.Time `query:"to"`
}

func (p *LogQueryParams) Normalize() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.Limit < 1 {
		p.Limit = 20
	}
	if p.Limit > 100 {
		p.Limit = 100
	}
}

type PaginatedLogs struct {
	Data       []LogDTO `json:"data"`
	Total      int64    `json:"total"`
	Page       int      `json:"page"`
	Limit      int      `json:"limit"`
	TotalPages int      `json:"total_pages"`
}

type LogDTO struct {
	ID           uint   `json:"id"`
	DeviceID     string `json:"device_id,omitempty"`
	Direction    string `json:"direction"`
	SimSlot      int    `json:"sim_slot"`
	Sender       string `json:"sender"`
	Receiver     string `json:"receiver"`
	Content      string `json:"content"`
	Status       string `json:"status"`
	ErrorMessage string `json:"error_message,omitempty"`
	RetryCount   int    `json:"retry_count"`
	CreatedAt    string `json:"created_at"`
	ProcessedAt  string `json:"processed_at,omitempty"`
}

func ToLogDTO(log *models.MessageLog) LogDTO {
	dto := LogDTO{
		ID:           log.ID,
		Direction:    string(log.Direction),
		SimSlot:      log.SimSlot,
		Sender:       log.Sender,
		Receiver:     log.Receiver,
		Content:      log.Content,
		Status:       string(log.Status),
		ErrorMessage: log.ErrorMessage,
		RetryCount:   log.RetryCount,
		CreatedAt:    log.CreatedAt.Format(time.RFC3339),
	}

	if log.DeviceID != nil {
		dto.DeviceID = log.DeviceID.String()
	}

	if log.ProcessedAt != nil {
		dto.ProcessedAt = log.ProcessedAt.Format(time.RFC3339)
	}

	return dto
}

func ToLogDTOList(logs []models.MessageLog) []LogDTO {
	dtos := make([]LogDTO, len(logs))
	for i, log := range logs {
		dtos[i] = ToLogDTO(&log)
	}
	return dtos
}

type LogStats struct {
	TotalInbound  int64 `json:"total_inbound"`
	TotalOutbound int64 `json:"total_outbound"`
	TotalSent     int64 `json:"total_sent"`
	TotalFailed   int64 `json:"total_failed"`
	TotalPending  int64 `json:"total_pending"`
}

type StatsQueryParams struct {
	From time.Time `query:"from"`
	To   time.Time `query:"to"`
}
