package worker

import (
	"context"
	"log"
	"schooli-api/cmd/config"
	"schooli-api/internal/repository/postgresql/db"
	"schooli-api/pkg/notifications"
	"time"

	"github.com/hibiken/asynq"
	"golang.org/x/exp/slog"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)
const TaskSendNotification = "task:send_notification"

type TaskProcessor interface {
	ProcessTaskSendNotification(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server             *asynq.Server
	c                  *config.Config
	redisOpt           asynq.RedisClientOpt
	notificationEngine *notifications.NotificationEngine
	store              db.Store
	l                  *slog.Logger
	worker             TaskDistributor
}

func NewRedisTaskProcessor(
	redisOpt asynq.RedisClientOpt,
	c *config.Config,
	notificationEngine *notifications.NotificationEngine,
	store db.Store,
	l *slog.Logger,
	worker TaskDistributor,

) *RedisTaskProcessor {
	server := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				QueueCritical: 10,
				QueueDefault:  5,
			},
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				log.Print(err)
				log.Print(task.Type())
			}),
			Logger: NewLogger(l),
		},
	)

	return &RedisTaskProcessor{
		server:             server,
		c:                  c,
		notificationEngine: notificationEngine,
		store:              store,
		l:                  l,
		worker:             worker,
	}
}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskSendNotification, processor.ProcessTaskSendNotification)
	return processor.server.Start(mux)
}

func (processor *RedisTaskProcessor) StartScheduler() error {
	l, err := time.LoadLocation("Africa/Nairobi")
	if err != nil {
		return err
	}
	mux := asynq.NewScheduler(processor.redisOpt, &asynq.SchedulerOpts{
		Location: l,
		Logger:   NewLogger(processor.l),
	})

	slog.Info("cron tasks scheduled", slog.Group("tasks"))
	return mux.Run()
}
