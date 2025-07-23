package repository

import (
	"context"
	"fmt"
	"time"

	"study_buddy/internal/model"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ Repository = (*NotifyRepo)(nil)

type NotifyRepo struct {
	db *pgxpool.Pool
}

func NewNotifyRepo(db *pgxpool.Pool) *NotifyRepo {
	return &NotifyRepo{
		db: db,
	}
}

type Repository interface {
	GetAllUsersEmail(ctx context.Context, userIDs []int64) ([]*model.User, error)
	GetUserIDs(ctx context.Context, now *time.Time) ([]int64, error)
}

func (n *NotifyRepo) GetAllUsersEmail(ctx context.Context, userIDs []int64) ([]*model.User, error) {
	pool, err := n.db.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer pool.Release()

	query, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Select("name", "email").
		From("users").
		Where(sq.Eq{"id": userIDs, "is_confirmed": true}).
		OrderBy("id").
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*model.User, 0)
	for rows.Next() {
		user := &model.User{}
		if err = rows.Scan(&user.Name, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

type NotificationData struct {
	UserID int64
	Time   time.Time
	Days   []int
}

func (n *NotifyRepo) GetUserIDs(ctx context.Context, now *time.Time) ([]int64, error) {
	pool, err := n.db.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer pool.Release()

	day := int(now.Weekday())
	if day == 0 {
		day = 7
	}

	query, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Select("user_id", "time_24", "days").
		From("notifications").
		Where(sq.Eq{"enabled": true}).
		OrderBy("id").
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]int64, 0)
	for rows.Next() {
		user := &NotificationData{}
		if err = rows.Scan(&user.UserID, &user.Time, &user.Days); err != nil {
			return nil, err
		}
		diff := HourDifferenceOnly(*now, user.Time)
		fmt.Println(diff)
		if contains(day, user.Days) && diff >= time.Duration(0) && diff < time.Minute {
			users = append(users, user.UserID)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func contains(day int, days []int) bool {
	for _, dayInt := range days {
		if dayInt == day {
			return true
		}
	}
	return false
}

func HourDifferenceOnly(t1, t2 time.Time) time.Duration {
	t1Time := time.Date(2000, 1, 1, t1.Hour(), t1.Minute(), t1.Second(), 0, time.UTC)
	t2Time := time.Date(2000, 1, 1, t2.Hour(), t2.Minute(), t2.Second(), 0, time.UTC)
	return t1Time.Sub(t2Time)
}
