package dto

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	User  UserDTO `json:"user"`
	Token string  `json:"token"`
}

type UserDTO struct {
	ID               string `json:"id"`
	Email            string `json:"email"`
	APIKey           string `json:"api_key"`
	SubscriptionPlan string `json:"subscription_plan"`
	Credits          int    `json:"credits"`
}

type APIKeyResponse struct {
	APIKey string `json:"api_key"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
