package repository

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/octopuslowtech/tinghook-project/backend/internal/models"
	"gorm.io/gorm"
)

var (
	ErrRuleNotFound = errors.New("forwarding rule not found")
)

type RuleRepository interface {
	Create(rule *models.ForwardingRule) error
	FindByID(id uint) (*models.ForwardingRule, error)
	FindByUserID(userID uuid.UUID) ([]models.ForwardingRule, error)
	FindByDeviceID(deviceID uuid.UUID) ([]models.ForwardingRule, error)
	FindActiveByDeviceAndType(deviceID uuid.UUID, triggerType string) ([]models.ForwardingRule, error)
	Update(rule *models.ForwardingRule) error
	Delete(id uint) error
}

type ruleRepository struct {
	db *gorm.DB
}

func NewRuleRepository(db *gorm.DB) RuleRepository {
	return &ruleRepository{db: db}
}

func (r *ruleRepository) Create(rule *models.ForwardingRule) error {
	if rule.CreatedAt.IsZero() {
		rule.CreatedAt = time.Now()
	}
	return r.db.Create(rule).Error
}

func (r *ruleRepository) FindByID(id uint) (*models.ForwardingRule, error) {
	var rule models.ForwardingRule
	err := r.db.Where("id = ?", id).First(&rule).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRuleNotFound
		}
		return nil, err
	}
	return &rule, nil
}

func (r *ruleRepository) FindByUserID(userID uuid.UUID) ([]models.ForwardingRule, error) {
	var rules []models.ForwardingRule
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&rules).Error
	return rules, err
}

func (r *ruleRepository) FindByDeviceID(deviceID uuid.UUID) ([]models.ForwardingRule, error) {
	var rules []models.ForwardingRule
	err := r.db.Where("device_id = ?", deviceID).Find(&rules).Error
	return rules, err
}

func (r *ruleRepository) FindActiveByDeviceAndType(deviceID uuid.UUID, triggerType string) ([]models.ForwardingRule, error) {
	var rules []models.ForwardingRule
	err := r.db.Where("device_id = ? AND trigger_type = ? AND is_active = ?", deviceID, triggerType, true).Find(&rules).Error
	return rules, err
}

func (r *ruleRepository) Update(rule *models.ForwardingRule) error {
	result := r.db.Save(rule)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrRuleNotFound
	}
	return nil
}

func (r *ruleRepository) Delete(id uint) error {
	result := r.db.Delete(&models.ForwardingRule{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrRuleNotFound
	}
	return nil
}
