package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/google/uuid"
	"github.com/octopuslowtech/tinghook-project/backend/internal/models"
	"github.com/octopuslowtech/tinghook-project/backend/internal/repository"
)

var (
	ErrRuleNotFound       = errors.New("forwarding rule not found")
	ErrRuleAccessDenied   = errors.New("access denied to this rule")
	ErrInvalidRegex       = errors.New("invalid regex pattern")
	ErrInvalidTriggerType = errors.New("invalid trigger type")
	ErrInvalidMethod      = errors.New("invalid HTTP method")
)

type CreateRuleRequest struct {
	DeviceID      *string `json:"device_id"`
	TriggerType   string  `json:"trigger_type"`
	SenderFilter  string  `json:"sender_filter"`
	ContentFilter string  `json:"content_filter"`
	WebhookURL    string  `json:"webhook_url"`
	SecretHeader  string  `json:"secret_header"`
	Method        string  `json:"method"`
}

type UpdateRuleRequest struct {
	DeviceID      *string `json:"device_id"`
	TriggerType   *string `json:"trigger_type"`
	SenderFilter  *string `json:"sender_filter"`
	ContentFilter *string `json:"content_filter"`
	WebhookURL    *string `json:"webhook_url"`
	SecretHeader  *string `json:"secret_header"`
	Method        *string `json:"method"`
	IsActive      *bool   `json:"is_active"`
}

type WebhookTestResult struct {
	Success      bool   `json:"success"`
	StatusCode   int    `json:"status_code"`
	ResponseTime int64  `json:"response_time_ms"`
	Error        string `json:"error,omitempty"`
}

type RuleService interface {
	Create(userID uuid.UUID, req *CreateRuleRequest) (*models.ForwardingRule, error)
	GetByID(id uint, userID uuid.UUID) (*models.ForwardingRule, error)
	ListByUser(userID uuid.UUID) ([]models.ForwardingRule, error)
	Update(id uint, userID uuid.UUID, req *UpdateRuleRequest) (*models.ForwardingRule, error)
	Delete(id uint, userID uuid.UUID) error
	TestWebhook(id uint, userID uuid.UUID) (*WebhookTestResult, error)
	MatchRules(deviceID uuid.UUID, triggerType string, sender, content string) ([]models.ForwardingRule, error)
}

type ruleService struct {
	repo       repository.RuleRepository
	httpClient *http.Client
}

func NewRuleService(repo repository.RuleRepository) RuleService {
	return &ruleService{
		repo: repo,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *ruleService) Create(userID uuid.UUID, req *CreateRuleRequest) (*models.ForwardingRule, error) {
	if err := validateTriggerType(req.TriggerType); err != nil {
		return nil, err
	}

	method := req.Method
	if method == "" {
		method = "POST"
	}
	if err := validateMethod(method); err != nil {
		return nil, err
	}

	if req.SenderFilter != "" {
		if _, err := regexp.Compile(req.SenderFilter); err != nil {
			return nil, fmt.Errorf("%w: sender_filter: %v", ErrInvalidRegex, err)
		}
	}

	if req.ContentFilter != "" {
		if _, err := regexp.Compile(req.ContentFilter); err != nil {
			return nil, fmt.Errorf("%w: content_filter: %v", ErrInvalidRegex, err)
		}
	}

	var deviceID *uuid.UUID
	if req.DeviceID != nil && *req.DeviceID != "" {
		parsed, err := uuid.Parse(*req.DeviceID)
		if err != nil {
			return nil, fmt.Errorf("invalid device_id: %w", err)
		}
		deviceID = &parsed
	}

	rule := &models.ForwardingRule{
		UserID:        userID,
		DeviceID:      deviceID,
		TriggerType:   req.TriggerType,
		SenderFilter:  req.SenderFilter,
		ContentFilter: req.ContentFilter,
		WebhookURL:    req.WebhookURL,
		SecretHeader:  req.SecretHeader,
		Method:        method,
		IsActive:      true,
		CreatedAt:     time.Now(),
	}

	if err := s.repo.Create(rule); err != nil {
		return nil, err
	}

	return rule, nil
}

func (s *ruleService) GetByID(id uint, userID uuid.UUID) (*models.ForwardingRule, error) {
	rule, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrRuleNotFound) {
			return nil, ErrRuleNotFound
		}
		return nil, err
	}

	if rule.UserID != userID {
		return nil, ErrRuleAccessDenied
	}

	return rule, nil
}

func (s *ruleService) ListByUser(userID uuid.UUID) ([]models.ForwardingRule, error) {
	return s.repo.FindByUserID(userID)
}

func (s *ruleService) Update(id uint, userID uuid.UUID, req *UpdateRuleRequest) (*models.ForwardingRule, error) {
	rule, err := s.GetByID(id, userID)
	if err != nil {
		return nil, err
	}

	if req.TriggerType != nil {
		if err := validateTriggerType(*req.TriggerType); err != nil {
			return nil, err
		}
		rule.TriggerType = *req.TriggerType
	}

	if req.Method != nil {
		if err := validateMethod(*req.Method); err != nil {
			return nil, err
		}
		rule.Method = *req.Method
	}

	if req.SenderFilter != nil {
		if *req.SenderFilter != "" {
			if _, err := regexp.Compile(*req.SenderFilter); err != nil {
				return nil, fmt.Errorf("%w: sender_filter: %v", ErrInvalidRegex, err)
			}
		}
		rule.SenderFilter = *req.SenderFilter
	}

	if req.ContentFilter != nil {
		if *req.ContentFilter != "" {
			if _, err := regexp.Compile(*req.ContentFilter); err != nil {
				return nil, fmt.Errorf("%w: content_filter: %v", ErrInvalidRegex, err)
			}
		}
		rule.ContentFilter = *req.ContentFilter
	}

	if req.DeviceID != nil {
		if *req.DeviceID == "" {
			rule.DeviceID = nil
		} else {
			parsed, err := uuid.Parse(*req.DeviceID)
			if err != nil {
				return nil, fmt.Errorf("invalid device_id: %w", err)
			}
			rule.DeviceID = &parsed
		}
	}

	if req.WebhookURL != nil {
		rule.WebhookURL = *req.WebhookURL
	}

	if req.SecretHeader != nil {
		rule.SecretHeader = *req.SecretHeader
	}

	if req.IsActive != nil {
		rule.IsActive = *req.IsActive
	}

	if err := s.repo.Update(rule); err != nil {
		return nil, err
	}

	return rule, nil
}

func (s *ruleService) Delete(id uint, userID uuid.UUID) error {
	rule, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrRuleNotFound) {
			return ErrRuleNotFound
		}
		return err
	}

	if rule.UserID != userID {
		return ErrRuleAccessDenied
	}

	return s.repo.Delete(id)
}

func (s *ruleService) TestWebhook(id uint, userID uuid.UUID) (*WebhookTestResult, error) {
	rule, err := s.GetByID(id, userID)
	if err != nil {
		return nil, err
	}

	testPayload := map[string]interface{}{
		"test":        true,
		"rule_id":     rule.ID,
		"trigger_type": rule.TriggerType,
		"sender":      "test_sender",
		"content":     "This is a test message from TingHook",
		"timestamp":   time.Now().Unix(),
	}

	payloadBytes, err := json.Marshal(testPayload)
	if err != nil {
		return &WebhookTestResult{
			Success: false,
			Error:   "failed to marshal payload",
		}, nil
	}

	start := time.Now()

	req, err := http.NewRequest(rule.Method, rule.WebhookURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return &WebhookTestResult{
			Success: false,
			Error:   fmt.Sprintf("failed to create request: %v", err),
		}, nil
	}

	req.Header.Set("Content-Type", "application/json")
	if rule.SecretHeader != "" {
		req.Header.Set("X-Webhook-Secret", rule.SecretHeader)
	}

	resp, err := s.httpClient.Do(req)
	elapsed := time.Since(start).Milliseconds()

	if err != nil {
		return &WebhookTestResult{
			Success:      false,
			ResponseTime: elapsed,
			Error:        fmt.Sprintf("request failed: %v", err),
		}, nil
	}
	defer resp.Body.Close()

	return &WebhookTestResult{
		Success:      resp.StatusCode >= 200 && resp.StatusCode < 300,
		StatusCode:   resp.StatusCode,
		ResponseTime: elapsed,
	}, nil
}

func (s *ruleService) MatchRules(deviceID uuid.UUID, triggerType string, sender, content string) ([]models.ForwardingRule, error) {
	rules, err := s.repo.FindActiveByDeviceAndType(deviceID, triggerType)
	if err != nil {
		return nil, err
	}

	var matched []models.ForwardingRule
	for _, rule := range rules {
		if matchesFilters(&rule, sender, content) {
			matched = append(matched, rule)
		}
	}

	return matched, nil
}

func matchesFilters(rule *models.ForwardingRule, sender, content string) bool {
	if rule.SenderFilter != "" {
		re, err := regexp.Compile(rule.SenderFilter)
		if err != nil {
			return false
		}
		if !re.MatchString(sender) {
			return false
		}
	}

	if rule.ContentFilter != "" {
		re, err := regexp.Compile(rule.ContentFilter)
		if err != nil {
			return false
		}
		if !re.MatchString(content) {
			return false
		}
	}

	return true
}

func validateTriggerType(t string) error {
	if t != "sms" && t != "notification" {
		return ErrInvalidTriggerType
	}
	return nil
}

func validateMethod(m string) error {
	if m != "GET" && m != "POST" && m != "PUT" {
		return ErrInvalidMethod
	}
	return nil
}
