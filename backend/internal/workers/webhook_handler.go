package workers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/hibiken/asynq"
	"github.com/octopuslowtech/tinghook-project/backend/internal/models"
	"github.com/octopuslowtech/tinghook-project/backend/internal/services"
)

type WebhookHandler struct {
	httpClient *http.Client
	logService services.LogService
}

func NewWebhookHandler(logService services.LogService) *WebhookHandler {
	return &WebhookHandler{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		logService: logService,
	}
}

func (h *WebhookHandler) HandleWebhookTask(ctx context.Context, t *asynq.Task) error {
	var payload WebhookPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		log.Printf("[webhook] failed to unmarshal payload: %v", err)
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	log.Printf("[webhook] dispatching to %s for log_id=%d", payload.WebhookURL, payload.LogID)

	body, err := json.Marshal(payload.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook data: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, payload.Method, payload.WebhookURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "TingHook-Webhook/1.0")

	if payload.SecretHeader != "" {
		req.Header.Set("X-Webhook-Secret", payload.SecretHeader)
	}

	resp, err := h.httpClient.Do(req)
	if err != nil {
		h.updateLogStatus(payload.LogID, models.StatusFailed, fmt.Sprintf("request failed: %v", err))
		return fmt.Errorf("webhook request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))

	if resp.StatusCode >= 400 {
		errMsg := fmt.Sprintf("webhook returned status %d: %s", resp.StatusCode, string(respBody))
		h.updateLogStatus(payload.LogID, models.StatusFailed, errMsg)
		return fmt.Errorf("webhook returned status %d", resp.StatusCode)
	}

	log.Printf("[webhook] successfully delivered to %s (status=%d)", payload.WebhookURL, resp.StatusCode)
	h.updateLogStatus(payload.LogID, models.StatusDelivered, "")

	return nil
}

func (h *WebhookHandler) updateLogStatus(logID uint, status models.MessageStatus, errorMsg string) {
	if h.logService == nil || logID == 0 {
		return
	}
	if err := h.logService.UpdateStatus(logID, status, errorMsg); err != nil {
		log.Printf("[webhook] failed to update log status for id=%d: %v", logID, err)
	}
}
