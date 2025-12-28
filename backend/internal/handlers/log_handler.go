package handlers

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/octopuslowtech/tinghook-project/backend/internal/handlers/dto"
	"github.com/octopuslowtech/tinghook-project/backend/internal/services"
)

type LogHandler struct {
	logService services.LogService
}

func NewLogHandler(logService services.LogService) *LogHandler {
	return &LogHandler{logService: logService}
}

func (h *LogHandler) ListLogs(c *fiber.Ctx) error {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	params := new(dto.LogQueryParams)
	if err := c.QueryParser(params); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid query parameters",
		})
	}

	result, err := h.logService.List(userID, params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch logs",
		})
	}

	return c.JSON(result)
}

func (h *LogHandler) GetLog(c *fiber.Ctx) error {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid log id",
		})
	}

	log, err := h.logService.GetByID(uint(id), userID)
	if err != nil {
		if errors.Is(err, services.ErrLogNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "log not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch log",
		})
	}

	return c.JSON(dto.ToLogDTO(log))
}

func (h *LogHandler) GetStats(c *fiber.Ctx) error {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	params := new(dto.StatsQueryParams)
	if err := c.QueryParser(params); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid query parameters",
		})
	}

	stats, err := h.logService.GetStats(userID, params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch stats",
		})
	}

	return c.JSON(stats)
}

func getUserIDFromContext(c *fiber.Ctx) (uuid.UUID, error) {
	userIDStr := c.Locals("user_id")
	if userIDStr == nil {
		return uuid.Nil, errors.New("user_id not found in context")
	}

	switch v := userIDStr.(type) {
	case uuid.UUID:
		return v, nil
	case string:
		return uuid.Parse(v)
	default:
		return uuid.Nil, errors.New("invalid user_id type")
	}
}

func (h *LogHandler) RegisterRoutes(router fiber.Router) {
	logs := router.Group("/logs")
	logs.Get("/", h.ListLogs)
	logs.Get("/stats", h.GetStats)
	logs.Get("/:id", h.GetLog)
}
