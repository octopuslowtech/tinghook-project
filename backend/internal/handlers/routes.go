package handlers

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/octopuslowtech/tinghook-project/backend/internal/middleware"
	"github.com/octopuslowtech/tinghook-project/backend/internal/services"
)

type Handlers struct {
	Auth *AuthHandler
	WS   *WSHandler
	SMS  *SMSHandler
}

func SetupRoutes(app *fiber.App, h *Handlers, jwtSecret string, userService services.UserService) {
	api := app.Group("/api")

	auth := api.Group("/auth")
	auth.Post("/register", h.Auth.Register)
	auth.Post("/login", h.Auth.Login)
	auth.Get("/me", middleware.JWTMiddleware(jwtSecret), h.Auth.GetMe)
	auth.Post("/refresh-key", middleware.JWTMiddleware(jwtSecret), h.Auth.RefreshAPIKey)

	v1 := api.Group("/v1")
	v1.Use(middleware.APIKeyMiddleware(userService))
	v1.Post("/sms/send", h.SMS.SendSMS)
	v1.Get("/devices/status", h.SMS.GetDevicesStatus)

	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/device", websocket.New(h.WS.HandleDeviceWS))
}
