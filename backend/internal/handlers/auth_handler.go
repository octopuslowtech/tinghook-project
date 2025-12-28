package handlers

import (
	"errors"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/octopuslowtech/tinghook-project/backend/internal/handlers/dto"
	"github.com/octopuslowtech/tinghook-project/backend/internal/models"
	"github.com/octopuslowtech/tinghook-project/backend/internal/services"
)

type AuthHandler struct {
	userService services.UserService
	jwtSecret   string
	jwtExpiry   time.Duration
}

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func NewAuthHandler(userService services.UserService, jwtSecret string, jwtExpiry time.Duration) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		jwtSecret:   jwtSecret,
		jwtExpiry:   jwtExpiry,
	}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "invalid request body"})
	}

	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "email and password are required"})
	}

	if !isValidEmail(req.Email) {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "invalid email format"})
	}

	if len(req.Password) < 8 {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "password must be at least 8 characters"})
	}

	user, err := h.userService.Register(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, services.ErrEmailAlreadyExists) {
			return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "email already exists"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Error: "failed to register user"})
	}

	token, err := h.generateToken(user.ID.String())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Error: "failed to generate token"})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.AuthResponse{
		User:  toUserDTO(user),
		Token: token,
	})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "invalid request body"})
	}

	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "email and password are required"})
	}

	user, err := h.userService.ValidateCredentials(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{Error: "invalid email or password"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Error: "failed to login"})
	}

	token, err := h.generateToken(user.ID.String())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Error: "failed to generate token"})
	}

	return c.JSON(dto.AuthResponse{
		User:  toUserDTO(user),
		Token: token,
	})
}

func (h *AuthHandler) GetMe(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	id, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "invalid user id"})
	}

	user, err := h.userService.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{Error: "user not found"})
	}

	return c.JSON(fiber.Map{"user": toUserDTO(user)})
}

func (h *AuthHandler) RefreshAPIKey(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	id, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "invalid user id"})
	}

	newKey, err := h.userService.RegenerateAPIKey(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Error: "failed to regenerate api key"})
	}

	return c.JSON(dto.APIKeyResponse{APIKey: newKey})
}

func (h *AuthHandler) generateToken(userID string) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(h.jwtExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.jwtSecret))
}

func toUserDTO(user *models.User) dto.UserDTO {
	return dto.UserDTO{
		ID:               user.ID.String(),
		Email:            user.Email,
		APIKey:           user.APIKey,
		SubscriptionPlan: user.SubscriptionPlan,
		Credits:          user.Credits,
	}
}

func isValidEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}
