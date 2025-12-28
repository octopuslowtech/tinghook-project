package services

import (
	"errors"
	"math"

	"github.com/google/uuid"
	"github.com/octopuslowtech/tinghook-project/backend/internal/handlers/dto"
	"github.com/octopuslowtech/tinghook-project/backend/internal/models"
	"github.com/octopuslowtech/tinghook-project/backend/internal/repository"
)

var (
	ErrLogNotFound = errors.New("message log not found")
)

type LogService interface {
	Create(userID uuid.UUID, deviceID *uuid.UUID, direction models.MessageDirection, sender, receiver, content string, simSlot int) (*models.MessageLog, error)
	GetByID(id uint, userID uuid.UUID) (*models.MessageLog, error)
	List(userID uuid.UUID, params *dto.LogQueryParams) (*dto.PaginatedLogs, error)
	GetStats(userID uuid.UUID, params *dto.StatsQueryParams) (*dto.LogStats, error)
	UpdateStatus(id uint, status models.MessageStatus, errorMsg string) error
	IncrementRetry(id uint) error
}

type logService struct {
	repo repository.LogRepository
}

func NewLogService(repo repository.LogRepository) LogService {
	return &logService{repo: repo}
}

func (s *logService) Create(userID uuid.UUID, deviceID *uuid.UUID, direction models.MessageDirection, sender, receiver, content string, simSlot int) (*models.MessageLog, error) {
	log := &models.MessageLog{
		UserID:    userID,
		DeviceID:  deviceID,
		Direction: direction,
		Sender:    sender,
		Receiver:  receiver,
		Content:   content,
		SimSlot:   simSlot,
		Status:    models.StatusPending,
	}

	if err := s.repo.Create(log); err != nil {
		return nil, err
	}

	return log, nil
}

func (s *logService) GetByID(id uint, userID uuid.UUID) (*models.MessageLog, error) {
	log, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrLogNotFound) {
			return nil, ErrLogNotFound
		}
		return nil, err
	}

	if log.UserID != userID {
		return nil, ErrLogNotFound
	}

	return log, nil
}

func (s *logService) List(userID uuid.UUID, params *dto.LogQueryParams) (*dto.PaginatedLogs, error) {
	params.Normalize()

	logs, total, err := s.repo.FindByUserID(userID, params)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(params.Limit)))

	return &dto.PaginatedLogs{
		Data:       dto.ToLogDTOList(logs),
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	}, nil
}

func (s *logService) GetStats(userID uuid.UUID, params *dto.StatsQueryParams) (*dto.LogStats, error) {
	return s.repo.GetStats(userID, params.From, params.To)
}

func (s *logService) UpdateStatus(id uint, status models.MessageStatus, errorMsg string) error {
	err := s.repo.UpdateStatus(id, status, errorMsg)
	if errors.Is(err, repository.ErrLogNotFound) {
		return ErrLogNotFound
	}
	return err
}

func (s *logService) IncrementRetry(id uint) error {
	err := s.repo.IncrementRetry(id)
	if errors.Is(err, repository.ErrLogNotFound) {
		return ErrLogNotFound
	}
	return err
}
