package notifications

import (
	"context"
	"path/filepath"
	"schooli-api/internal/models"

	firebase "firebase.google.com/go"
	"golang.org/x/exp/slog"

	"firebase.google.com/go/messaging"
	"github.com/pkg/errors"
	"google.golang.org/api/option"
)

type NotificationEngine struct {
	app       *firebase.App
	messenger *messaging.Client
}

func NewNotificationEngine(file string) (*NotificationEngine, error) {
	ctx := context.Background()
	serviceAccountKeyFilePath, err := filepath.Abs(file)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to load serviceAccountKeys.json file")
	}

	opt := option.WithCredentialsFile(serviceAccountKeyFilePath)

	// Firebase admin SDK initialization.
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, errors.Wrap(err, "Firebase loading error")
	}

	// Messaging client.
	client, err := app.Messaging(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Client fb loading error")
	}

	return &NotificationEngine{
		app:       app,
		messenger: client,
	}, nil
}

func (n *NotificationEngine) SendMessage(ctx context.Context, notification models.Notifications) error {
	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: notification.Title,
			Body:  notification.Body,
		},
		Token: notification.Device.Device,
	}

	response, err := n.messenger.Send(ctx, message)
	if err != nil {
		return errors.Wrap(err, "error sending message")
	}
	slog.Info("Successfully sent message:", response)
	return nil
}

func (n *NotificationEngine) SendBulkMessage(ctx context.Context, devices []string, m models.Notifications) error {
	message := &messaging.Message{

		Notification: &messaging.Notification{
			Title: m.Title,
			Body:  m.Body,
		},
	}

	response, err := n.messenger.SendMulticast(ctx, &messaging.MulticastMessage{
		Tokens:       devices,
		Notification: message.Notification,
	})
	if err != nil {
		return errors.Wrap(err, "error sending message")
	}
	slog.Info("Successfully sent message:", response)
	return nil
}
