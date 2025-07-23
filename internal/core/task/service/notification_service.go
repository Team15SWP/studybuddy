package service

import (
	"context"
	"fmt"

	"study_buddy/internal/model"
)

func (t *TaskService) GetNotificationSettings(ctx context.Context, userId int64) (*model.Notification, error) {
	response, err := t.notificationRepo.GetNotification(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("t.notificationRepo.GetNotification: %w", err)
	}
	return response, nil
}

func (t *TaskService) SetNotificationSettings(ctx context.Context, notif *model.Notification) error {
	err := t.notificationRepo.UpdateNotification(ctx, notif)
	if err != nil {
		return fmt.Errorf("t.notificationRepo.UpdateNotification: %w", err)
	}
	return nil
}
