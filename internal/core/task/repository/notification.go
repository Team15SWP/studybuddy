package repository

import (
	"context"
	"fmt"

	"study_buddy/internal/model"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ NotificationRepository = (*NotificationRepo)(nil)

type NotificationRepo struct {
	db *pgxpool.Pool
}

func NewNotificationRepo(db *pgxpool.Pool) *NotificationRepo {
	return &NotificationRepo{
		db: db,
	}
}

type NotificationRepository interface {
	GetNotification(ctx context.Context, userId int64) (*model.Notification, error)
	CreateNotification(ctx context.Context, notif *model.Notification) error
	UpdateNotification(ctx context.Context, notif *model.Notification) error
}

func (s *NotificationRepo) GetNotification(ctx context.Context, userId int64) (*model.Notification, error) {
	pool, err := s.db.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer pool.Release()

	query, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Select("user_id", "enabled", "time_24", "days").
		From("notifications").
		Where(sq.Eq{"user_id": userId}).
		ToSql()
	if err != nil {
		return nil, err
	}

	row := pool.QueryRow(ctx, query, args...)
	var stats model.Notification

	err = row.Scan(&stats.UserID, &stats.Enabled, &stats.Time24, &stats.Days)
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

func (s *NotificationRepo) CreateNotification(ctx context.Context, notif *model.Notification) error {
	pool, err := s.db.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("db acquire: %w", err)
	}
	defer pool.Release()

	query, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Insert("notifications").
		Columns("user_id", "enabled", "time_24", "days").
		Values(notif.UserID, notif.Enabled, notif.Time24, notif.Days).
		ToSql()
	if err != nil {
		return fmt.Errorf("build insert query: %w", err)
	}

	_, err = pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("exec insert: %w", err)
	}
	return nil
}

func (s *NotificationRepo) UpdateNotification(ctx context.Context, notif *model.Notification) error {
	pool, err := s.db.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("db acquire: %w", err)
	}
	defer pool.Release()

	query, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Update("notifications").
		Set("enabled", notif.Enabled).
		Set("time_24", notif.Time24).
		Set("days", notif.Days).
		Where(sq.Eq{"user_id": notif.UserID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build insert query: %w", err)
	}

	_, err = pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("exec insert: %w", err)
	}
	return nil
}
