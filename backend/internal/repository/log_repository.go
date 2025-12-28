package repository

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/octopuslowtech/tinghook-project/backend/internal/handlers/dto"
	"github.com/octopuslowtech/tinghook-project/backend/internal/models"
	"gorm.io/gorm"
)

var (
	ErrLogNotFound = errors.New("message log not found")
)

type LogRepository interface {
	Create(log *models.MessageLog) error
	FindByID(id uint) (*models.MessageLog, error)
	FindByUserID(userID uuid.UUID, params *dto.LogQueryParams) ([]models.MessageLog, int64, error)
	UpdateStatus(id uint, status models.MessageStatus, errorMsg string) error
	IncrementRetry(id uint) error
	GetStats(userID uuid.UUID, from, to time.Time) (*dto.LogStats, error)
}

type logRepository struct {
	db *gorm.DB
}

func NewLogRepository(db *gorm.DB) LogRepository {
	return &logRepository{db: db}
}

func (r *logRepository) Create(log *models.MessageLog) error {
	if log.CreatedAt.IsZero() {
		log.CreatedAt = time.Now()
	}
	return r.db.Create(log).Error
}

func (r *logRepository) FindByID(id uint) (*models.MessageLog, error) {
	var log models.MessageLog
	err := r.db.Where("id = ?", id).First(&log).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrLogNotFound
		}
		return nil, err
	}
	return &log, nil
}

func (r *logRepository) FindByUserID(userID uuid.UUID, params *dto.LogQueryParams) ([]models.MessageLog, int64, error) {
	var logs []models.MessageLog
	var total int64

	query := r.db.Model(&models.MessageLog{}).Where("user_id = ?", userID)

	if params.Direction != "" {
		query = query.Where("direction = ?", params.Direction)
	}

	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	if params.DeviceID != "" {
		deviceUUID, err := uuid.Parse(params.DeviceID)
		if err == nil {
			query = query.Where("device_id = ?", deviceUUID)
		}
	}

	if !params.From.IsZero() {
		query = query.Where("created_at >= ?", params.From)
	}

	if !params.To.IsZero() {
		query = query.Where("created_at <= ?", params.To)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.Limit
	err := query.Order("created_at DESC").Offset(offset).Limit(params.Limit).Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

func (r *logRepository) UpdateStatus(id uint, status models.MessageStatus, errorMsg string) error {
	now := time.Now()
	updates := map[string]interface{}{
		"status":       status,
		"processed_at": now,
	}

	if errorMsg != "" {
		updates["error_message"] = errorMsg
	}

	result := r.db.Model(&models.MessageLog{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrLogNotFound
	}
	return nil
}

func (r *logRepository) IncrementRetry(id uint) error {
	result := r.db.Model(&models.MessageLog{}).Where("id = ?", id).
		UpdateColumn("retry_count", gorm.Expr("retry_count + 1"))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrLogNotFound
	}
	return nil
}

func (r *logRepository) GetStats(userID uuid.UUID, from, to time.Time) (*dto.LogStats, error) {
	stats := &dto.LogStats{}

	query := r.db.Model(&models.MessageLog{}).Where("user_id = ?", userID)

	if !from.IsZero() {
		query = query.Where("created_at >= ?", from)
	}
	if !to.IsZero() {
		query = query.Where("created_at <= ?", to)
	}

	if err := query.Where("direction = ?", models.DirectionInbound).Count(&stats.TotalInbound).Error; err != nil {
		return nil, err
	}

	query = r.db.Model(&models.MessageLog{}).Where("user_id = ?", userID)
	if !from.IsZero() {
		query = query.Where("created_at >= ?", from)
	}
	if !to.IsZero() {
		query = query.Where("created_at <= ?", to)
	}

	if err := query.Where("direction = ?", models.DirectionOutbound).Count(&stats.TotalOutbound).Error; err != nil {
		return nil, err
	}

	query = r.db.Model(&models.MessageLog{}).Where("user_id = ?", userID)
	if !from.IsZero() {
		query = query.Where("created_at >= ?", from)
	}
	if !to.IsZero() {
		query = query.Where("created_at <= ?", to)
	}

	if err := query.Where("status = ?", models.StatusSent).Count(&stats.TotalSent).Error; err != nil {
		return nil, err
	}

	query = r.db.Model(&models.MessageLog{}).Where("user_id = ?", userID)
	if !from.IsZero() {
		query = query.Where("created_at >= ?", from)
	}
	if !to.IsZero() {
		query = query.Where("created_at <= ?", to)
	}

	if err := query.Where("status = ?", models.StatusFailed).Count(&stats.TotalFailed).Error; err != nil {
		return nil, err
	}

	query = r.db.Model(&models.MessageLog{}).Where("user_id = ?", userID)
	if !from.IsZero() {
		query = query.Where("created_at >= ?", from)
	}
	if !to.IsZero() {
		query = query.Where("created_at <= ?", to)
	}

	if err := query.Where("status = ?", models.StatusPending).Count(&stats.TotalPending).Error; err != nil {
		return nil, err
	}

	return stats, nil
}
