package worker

import (
	"context"
	"schooli-api/internal/models"

	"github.com/hibiken/asynq"
)

type TaskDistributor interface {
	DistributeTaskSendNotification(ctx context.Context, notification models.Notifications, opts ...asynq.Option) error
}

type RedisTaskDistributor struct {
	client *asynq.Client
}

func NewRedisTaskDistributor(redisOpt asynq.RedisClientOpt) TaskDistributor {
	client := asynq.NewClient(redisOpt)
	return &RedisTaskDistributor{
		client: client,
	}
}
