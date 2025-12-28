package handlers

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/octopuslowtech/tinghook-project/backend/internal/handlers/dto"
	"github.com/octopuslowtech/tinghook-project/backend/internal/services"
)

type RuleHandler struct {
	ruleService services.RuleService
}

func NewRuleHandler(ruleService services.RuleService) *RuleHandler {
	return &RuleHandler{
		ruleService: ruleService,
	}
}

func (h *RuleHandler) RegisterRoutes(router fiber.Router, authMiddleware fiber.Handler) {
	rules := router.Group("/rules", authMiddleware)
	rules.Get("/", h.List)
	rules.Post("/", h.Create)
	rules.Get("/:id", h.Get)
	rules.Put("/:id", h.Update)
	rules.Delete("/:id", h.Delete)
	rules.Post("/:id/test", h.TestWebhook)
}

func (h *RuleHandler) List(c *fiber.Ctx) error {
	userID, err := getRuleUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	rules, err := h.ruleService.ListByUser(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch rules",
		})
	}

	return c.JSON(fiber.Map{
		"rules": dto.ToRuleDTOList(rules),
	})
}

func (h *RuleHandler) Create(c *fiber.Ctx) error {
	userID, err := getRuleUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	var req dto.CreateRuleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if req.TriggerType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "trigger_type is required",
		})
	}

	if req.WebhookURL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "webhook_url is required",
		})
	}

	serviceReq := &services.CreateRuleRequest{
		DeviceID:      req.DeviceID,
		TriggerType:   req.TriggerType,
		SenderFilter:  req.SenderFilter,
		ContentFilter: req.ContentFilter,
		WebhookURL:    req.WebhookURL,
		SecretHeader:  req.SecretHeader,
		Method:        req.Method,
	}

	rule, err := h.ruleService.Create(userID, serviceReq)
	if err != nil {
		if errors.Is(err, services.ErrInvalidTriggerType) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "trigger_type must be 'sms' or 'notification'",
			})
		}
		if errors.Is(err, services.ErrInvalidMethod) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "method must be GET, POST, or PUT",
			})
		}
		if errors.Is(err, services.ErrInvalidRegex) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create rule",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"rule": dto.ToRuleDTO(rule),
	})
}

func (h *RuleHandler) Get(c *fiber.Ctx) error {
	userID, err := getRuleUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	id, err := parseRuleID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid rule id",
		})
	}

	rule, err := h.ruleService.GetByID(id, userID)
	if err != nil {
		if errors.Is(err, services.ErrRuleNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "rule not found",
			})
		}
		if errors.Is(err, services.ErrRuleAccessDenied) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "access denied",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch rule",
		})
	}

	return c.JSON(fiber.Map{
		"rule": dto.ToRuleDTO(rule),
	})
}

func (h *RuleHandler) Update(c *fiber.Ctx) error {
	userID, err := getRuleUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	id, err := parseRuleID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid rule id",
		})
	}

	var req dto.UpdateRuleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	serviceReq := &services.UpdateRuleRequest{
		DeviceID:      req.DeviceID,
		TriggerType:   req.TriggerType,
		SenderFilter:  req.SenderFilter,
		ContentFilter: req.ContentFilter,
		WebhookURL:    req.WebhookURL,
		SecretHeader:  req.SecretHeader,
		Method:        req.Method,
		IsActive:      req.IsActive,
	}

	rule, err := h.ruleService.Update(id, userID, serviceReq)
	if err != nil {
		if errors.Is(err, services.ErrRuleNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "rule not found",
			})
		}
		if errors.Is(err, services.ErrRuleAccessDenied) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "access denied",
			})
		}
		if errors.Is(err, services.ErrInvalidTriggerType) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "trigger_type must be 'sms' or 'notification'",
			})
		}
		if errors.Is(err, services.ErrInvalidMethod) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "method must be GET, POST, or PUT",
			})
		}
		if errors.Is(err, services.ErrInvalidRegex) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to update rule",
		})
	}

	return c.JSON(fiber.Map{
		"rule": dto.ToRuleDTO(rule),
	})
}

func (h *RuleHandler) Delete(c *fiber.Ctx) error {
	userID, err := getRuleUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	id, err := parseRuleID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid rule id",
		})
	}

	err = h.ruleService.Delete(id, userID)
	if err != nil {
		if errors.Is(err, services.ErrRuleNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "rule not found",
			})
		}
		if errors.Is(err, services.ErrRuleAccessDenied) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "access denied",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to delete rule",
		})
	}

	return c.JSON(fiber.Map{
		"message": "rule deleted",
	})
}

func (h *RuleHandler) TestWebhook(c *fiber.Ctx) error {
	userID, err := getRuleUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	id, err := parseRuleID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid rule id",
		})
	}

	result, err := h.ruleService.TestWebhook(id, userID)
	if err != nil {
		if errors.Is(err, services.ErrRuleNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "rule not found",
			})
		}
		if errors.Is(err, services.ErrRuleAccessDenied) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "access denied",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to test webhook",
		})
	}

	return c.JSON(fiber.Map{
		"result": dto.WebhookTestResult{
			Success:      result.Success,
			StatusCode:   result.StatusCode,
			ResponseTime: result.ResponseTime,
			Error:        result.Error,
		},
	})
}

func getRuleUserID(c *fiber.Ctx) (uuid.UUID, error) {
	userIDVal := c.Locals("user_id")
	if userIDVal == nil {
		return uuid.Nil, errors.New("user_id not found in context")
	}

	switch v := userIDVal.(type) {
	case uuid.UUID:
		return v, nil
	case string:
		return uuid.Parse(v)
	default:
		return uuid.Nil, errors.New("invalid user_id type")
	}
}

func parseRuleID(c *fiber.Ctx) (uint, error) {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}
