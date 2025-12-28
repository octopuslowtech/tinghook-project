package repository

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/octopuslowtech/tinghook-project/backend/internal/models"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrDuplicateEmail  = errors.New("email already exists")
	ErrDuplicateAPIKey = errors.New("api key already exists")
)

type UserRepository interface {
	Create(user *models.User) error
	FindByID(id uuid.UUID) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	FindByAPIKey(apiKey string) (*models.User, error)
	Update(user *models.User) error
	Delete(id uuid.UUID) error
	UpdateAPIKey(id uuid.UUID, newKey string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) error {
	result := r.db.Create(user)
	if result.Error != nil {
		errMsg := result.Error.Error()
		if isDuplicateKeyError(result.Error) && strings.Contains(errMsg, "email") {
			return ErrDuplicateEmail
		}
		if isDuplicateKeyError(result.Error) && strings.Contains(errMsg, "api_key") {
			return ErrDuplicateAPIKey
		}
		return result.Error
	}
	return nil
}

func (r *userRepository) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	result := r.db.First(&user, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, result.Error
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	result := r.db.First(&user, "email = ?", email)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, result.Error
	}
	return &user, nil
}

func (r *userRepository) FindByAPIKey(apiKey string) (*models.User, error) {
	var user models.User
	result := r.db.First(&user, "api_key = ?", apiKey)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, result.Error
	}
	return &user, nil
}

func (r *userRepository) Update(user *models.User) error {
	result := r.db.Save(user)
	if result.Error != nil {
		errMsg := result.Error.Error()
		if isDuplicateKeyError(result.Error) && strings.Contains(errMsg, "email") {
			return ErrDuplicateEmail
		}
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (r *userRepository) Delete(id uuid.UUID) error {
	result := r.db.Delete(&models.User{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (r *userRepository) UpdateAPIKey(id uuid.UUID, newKey string) error {
	result := r.db.Model(&models.User{}).Where("id = ?", id).Update("api_key", newKey)
	if result.Error != nil {
		if isDuplicateKeyError(result.Error) {
			return ErrDuplicateAPIKey
		}
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}
