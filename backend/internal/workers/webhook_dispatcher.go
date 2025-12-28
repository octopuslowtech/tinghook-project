package workers

import (
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
)

const (
	TypeWebhookDispatch = "webhook:dispatch"
)

type WebhookPayload struct {
	RuleID       uint        `json:"rule_id"`
	WebhookURL   string      `json:"webhook_url"`
	Method       string      `json:"method"`
	SecretHeader string      `json:"secret_header"`
	Data         WebhookData `json:"data"`
	LogID        uint        `json:"log_id"`
}

type WebhookData struct {
	Type      string `json:"type"`
	DeviceID  string `json:"device_id"`
	Sender    string `json:"sender,omitempty"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`

	AppPackage string `json:"app_package,omitempty"`
	AppName    string `json:"app_name,omitempty"`
	Title      string `json:"title,omitempty"`
}

type WebhookDispatcher struct {
	client *asynq.Client
}

func NewWebhookDispatcher(redisAddr string) *WebhookDispatcher {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	return &WebhookDispatcher{client: client}
}

func (d *WebhookDispatcher) Close() error {
	return d.client.Close()
}

func (d *WebhookDispatcher) Dispatch(payload *WebhookPayload) error {
	task, err := NewWebhookTask(payload)
	if err != nil {
		return err
	}
	_, err = d.client.Enqueue(task,
		asynq.MaxRetry(3),
		asynq.Timeout(30*time.Second),
		asynq.Queue("webhooks"),
	)
	return err
}

func NewWebhookTask(payload *WebhookPayload) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeWebhookDispatch, data), nil
}
