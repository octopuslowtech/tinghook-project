package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"

	"github.com/google/uuid"
	"github.com/octopuslowtech/tinghook-project/backend/internal/models"
	"github.com/octopuslowtech/tinghook-project/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost   = 12
	apiKeyLength = 32 // 32 bytes = 64 hex characters
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailAlreadyExists = errors.New("email already exists")
)

type UserService interface {
	Register(email, password string) (*models.User, error)
	ValidateCredentials(email, password string) (*models.User, error)
	GetByID(id uuid.UUID) (*models.User, error)
	GetByAPIKey(apiKey string) (*models.User, error)
	RegenerateAPIKey(id uuid.UUID) (string, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) Register(email, password string) (*models.User, error) {
	passwordHash, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	apiKey, err := generateAPIKey()
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:               uuid.New(),
		Email:            email,
		PasswordHash:     passwordHash,
		APIKey:           apiKey,
		SubscriptionPlan: "free",
		Credits:          0,
	}

	if err := s.userRepo.Create(user); err != nil {
		if errors.Is(err, repository.ErrDuplicateEmail) {
			return nil, ErrEmailAlreadyExists
		}
		return nil, err
	}

	return user, nil
}

func (s *userService) ValidateCredentials(email, password string) (*models.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}

func (s *userService) GetByID(id uuid.UUID) (*models.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *userService) GetByAPIKey(apiKey string) (*models.User, error) {
	return s.userRepo.FindByAPIKey(apiKey)
}

func (s *userService) RegenerateAPIKey(id uuid.UUID) (string, error) {
	newKey, err := generateAPIKey()
	if err != nil {
		return "", err
	}

	if err := s.userRepo.UpdateAPIKey(id, newKey); err != nil {
		return "", err
	}

	return newKey, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func generateAPIKey() (string, error) {
	bytes := make([]byte, apiKeyLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
