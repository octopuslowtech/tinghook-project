package repository

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/octopuslowtech/tinghook-project/backend/internal/models"
	"gorm.io/gorm"
)

var (
	ErrDeviceNotFound      = errors.New("device not found")
	ErrDuplicateDeviceUID  = errors.New("device UID already exists")
)

type DeviceRepository interface {
	Create(device *models.Device) error
	FindByID(id uuid.UUID) (*models.Device, error)
	FindByDeviceUID(uid string) (*models.Device, error)
	FindByUserID(userID uuid.UUID) ([]models.Device, error)
	Update(device *models.Device) error
	Delete(id uuid.UUID) error
	UpdateStatus(id uuid.UUID, status string) error
	UpdateLastSeen(id uuid.UUID, battery int) error
	UpdateFCMToken(id uuid.UUID, token string) error
}

type deviceRepository struct {
	db *gorm.DB
}

func NewDeviceRepository(db *gorm.DB) DeviceRepository {
	return &deviceRepository{db: db}
}

func (r *deviceRepository) Create(device *models.Device) error {
	if device.ID == uuid.Nil {
		device.ID = uuid.New()
	}
	if device.CreatedAt.IsZero() {
		device.CreatedAt = time.Now()
	}

	err := r.db.Create(device).Error
	if err != nil && isDuplicateKeyError(err) {
		return ErrDuplicateDeviceUID
	}
	return err
}

func (r *deviceRepository) FindByID(id uuid.UUID) (*models.Device, error) {
	var device models.Device
	err := r.db.Where("id = ?", id).First(&device).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDeviceNotFound
		}
		return nil, err
	}
	return &device, nil
}

func (r *deviceRepository) FindByDeviceUID(uid string) (*models.Device, error) {
	var device models.Device
	err := r.db.Where("device_uid = ?", uid).First(&device).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDeviceNotFound
		}
		return nil, err
	}
	return &device, nil
}

func (r *deviceRepository) FindByUserID(userID uuid.UUID) ([]models.Device, error) {
	var devices []models.Device
	err := r.db.Where("user_id = ?", userID).Find(&devices).Error
	return devices, err
}

func (r *deviceRepository) Update(device *models.Device) error {
	result := r.db.Save(device)
	if result.Error != nil {
		if isDuplicateKeyError(result.Error) {
			return ErrDuplicateDeviceUID
		}
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrDeviceNotFound
	}
	return nil
}

func (r *deviceRepository) Delete(id uuid.UUID) error {
	result := r.db.Delete(&models.Device{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrDeviceNotFound
	}
	return nil
}

func (r *deviceRepository) UpdateStatus(id uuid.UUID, status string) error {
	result := r.db.Model(&models.Device{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrDeviceNotFound
	}
	return nil
}

func (r *deviceRepository) UpdateLastSeen(id uuid.UUID, battery int) error {
	now := time.Now()
	result := r.db.Model(&models.Device{}).Where("id = ?", id).Updates(map[string]interface{}{
		"last_seen_at":  now,
		"battery_level": battery,
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrDeviceNotFound
	}
	return nil
}

func (r *deviceRepository) UpdateFCMToken(id uuid.UUID, token string) error {
	result := r.db.Model(&models.Device{}).Where("id = ?", id).Update("fcm_token", token)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrDeviceNotFound
	}
	return nil
}

func isDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return contains(errStr, "duplicate key") || contains(errStr, "UNIQUE constraint") || contains(errStr, "1062")
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsRune(s, substr))
}

func containsRune(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
