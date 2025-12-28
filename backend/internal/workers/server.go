package workers

import (
	"log"

	"github.com/hibiken/asynq"
	"github.com/octopuslowtech/tinghook-project/backend/internal/services"
)

type WorkerServer struct {
	srv *asynq.Server
}

func StartWorkerServer(redisAddr string, logService services.LogService) *WorkerServer {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"webhooks": 6,
				"default":  4,
			},
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx asynq.Context, task *asynq.Task, err error) {
				log.Printf("[worker] task %s failed: %v", task.Type(), err)
			}),
		},
	)

	handler := NewWebhookHandler(logService)
	mux := asynq.NewServeMux()
	mux.HandleFunc(TypeWebhookDispatch, handler.HandleWebhookTask)

	go func() {
		log.Printf("[worker] starting asynq server on redis=%s", redisAddr)
		if err := srv.Start(mux); err != nil {
			log.Fatalf("[worker] could not start worker server: %v", err)
		}
	}()

	return &WorkerServer{srv: srv}
}

func (w *WorkerServer) Shutdown() {
	if w.srv != nil {
		w.srv.Shutdown()
	}
}
