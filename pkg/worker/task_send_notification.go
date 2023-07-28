package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"schooli-api/internal/models"

	"github.com/hibiken/asynq"
	"golang.org/x/exp/slog"
)

func (processor *RedisTaskProcessor) ProcessTaskSendNotification(ctx context.Context, task *asynq.Task) error {
	var payload models.Notifications
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	if err := processor.notificationEngine.SendMessage(ctx, payload); err != nil {
		return fmt.Errorf("unable to index business %w", err)
	}
	slog.Info("processed task", "type", task.Type(), "payload", string(task.Payload()))
	return nil
}

func (distributor *RedisTaskDistributor) DistributeTaskSendNotification(ctx context.Context, notification models.Notifications, opts ...asynq.Option) error {
	jsonPayload, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}

	task := asynq.NewTask(TaskSendNotification, jsonPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	slog.Info("enqueued task", "type", task.Type(), "payload", string(task.Payload()), "queue", info.Queue, "max_retry", info.MaxRetry)
	return nil
}
