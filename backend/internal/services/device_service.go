package services

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/octopuslowtech/tinghook-project/backend/internal/models"
	"github.com/octopuslowtech/tinghook-project/backend/internal/repository"
)

var (
	ErrDeviceNotFound     = errors.New("device not found")
	ErrDuplicateDeviceUID = errors.New("device UID already exists")
)

type DeviceService interface {
	Register(userID uuid.UUID, name, deviceUID string) (*models.Device, error)
	GetByID(id uuid.UUID) (*models.Device, error)
	GetByDeviceUID(uid string) (*models.Device, error)
	ListByUser(userID uuid.UUID) ([]models.Device, error)
	UpdateName(id uuid.UUID, name string) error
	Delete(id uuid.UUID) error
	SetOnline(id uuid.UUID, battery int) error
	SetOffline(id uuid.UUID) error
	UpdateFCMToken(id uuid.UUID, token string) error
	GeneratePairingToken(userID uuid.UUID) (string, error)
	GenerateDeviceUID() (string, error)
}

type deviceService struct {
	repo repository.DeviceRepository
}

func NewDeviceService(repo repository.DeviceRepository) DeviceService {
	return &deviceService{repo: repo}
}

func (s *deviceService) Register(userID uuid.UUID, name, deviceUID string) (*models.Device, error) {
	device := &models.Device{
		ID:        uuid.New(),
		UserID:    userID,
		Name:      name,
		DeviceUID: deviceUID,
		Status:    models.DeviceStatusOffline,
		CreatedAt: time.Now(),
	}

	err := s.repo.Create(device)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateDeviceUID) {
			return nil, ErrDuplicateDeviceUID
		}
		return nil, err
	}

	return device, nil
}

func (s *deviceService) GetByID(id uuid.UUID) (*models.Device, error) {
	device, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrDeviceNotFound) {
			return nil, ErrDeviceNotFound
		}
		return nil, err
	}
	return device, nil
}

func (s *deviceService) GetByDeviceUID(uid string) (*models.Device, error) {
	device, err := s.repo.FindByDeviceUID(uid)
	if err != nil {
		if errors.Is(err, repository.ErrDeviceNotFound) {
			return nil, ErrDeviceNotFound
		}
		return nil, err
	}
	return device, nil
}

func (s *deviceService) ListByUser(userID uuid.UUID) ([]models.Device, error) {
	return s.repo.FindByUserID(userID)
}

func (s *deviceService) UpdateName(id uuid.UUID, name string) error {
	device, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrDeviceNotFound) {
			return ErrDeviceNotFound
		}
		return err
	}

	device.Name = name
	return s.repo.Update(device)
}

func (s *deviceService) Delete(id uuid.UUID) error {
	err := s.repo.Delete(id)
	if errors.Is(err, repository.ErrDeviceNotFound) {
		return ErrDeviceNotFound
	}
	return err
}

func (s *deviceService) SetOnline(id uuid.UUID, battery int) error {
	if err := s.repo.UpdateStatus(id, models.DeviceStatusOnline); err != nil {
		if errors.Is(err, repository.ErrDeviceNotFound) {
			return ErrDeviceNotFound
		}
		return err
	}
	return s.repo.UpdateLastSeen(id, battery)
}

func (s *deviceService) SetOffline(id uuid.UUID) error {
	err := s.repo.UpdateStatus(id, models.DeviceStatusOffline)
	if errors.Is(err, repository.ErrDeviceNotFound) {
		return ErrDeviceNotFound
	}
	return err
}

func (s *deviceService) UpdateFCMToken(id uuid.UUID, token string) error {
	err := s.repo.UpdateFCMToken(id, token)
	if errors.Is(err, repository.ErrDeviceNotFound) {
		return ErrDeviceNotFound
	}
	return err
}

func (s *deviceService) GeneratePairingToken(userID uuid.UUID) (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	token := base64.URLEncoding.EncodeToString(bytes)
	return token, nil
}

func (s *deviceService) GenerateDeviceUID() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
