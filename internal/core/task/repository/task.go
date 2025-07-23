package repository

import (
	"context"
	"fmt"
	"strings"

	"study_buddy/internal/model"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ TaskRepository = (*TaskRepo)(nil)

type TaskRepo struct {
	db *pgxpool.Pool
}

func NewTaskRepo(db *pgxpool.Pool) *TaskRepo {
	return &TaskRepo{
		db: db,
	}
}

type TaskRepository interface {
	CreateTask(ctx context.Context, task *model.GeneratedTask) error
	GetTask(ctx context.Context, userId int64, taskName string) (*model.GeneratedTask, error)
	UpdateTaskSolved(ctx context.Context, task *model.GeneratedTask, stats *model.Statistics) error
}

func (t *TaskRepo) CreateTask(ctx context.Context, task *model.GeneratedTask) error {
	pool, err := t.db.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("t.db.Acquire: %w", err)
	}
	defer pool.Release()

	columns := []string{"user_id", "task", "description", "solution", "hint1", "hint2", "hint3", "difficulty", "solved"}

	query, args, err := sq.Insert("tasks").
		Columns(columns...).
		Values(task.UserID, task.TaskName, task.TaskDescription, task.Solution, task.Hints.Hint1, task.Hints.Hint2,
			task.Hints.Hint3, task.Difficulty, task.Solved).
		Suffix("RETURNING id, " + strings.Join(columns, ", ")).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("format user insert SQL: %w", err)
	}

	_, err = pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("insert task: %w", err)
	}
	return nil
}

func (t *TaskRepo) GetTask(ctx context.Context, userId int64, taskName string) (*model.GeneratedTask, error) {
	pool, err := t.db.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer pool.Release()

	query, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Select("id", "user_id", "task", "description", "solution", "hint1", "hint2", "hint3", "difficulty", "solved").
		From("tasks").
		Where(sq.Eq{"user_id": userId, "task": taskName}).
		ToSql()
	if err != nil {
		return nil, err
	}

	row := pool.QueryRow(ctx, query, args...)
	var task model.GeneratedTask

	err = row.Scan(&task.ID, &task.UserID, &task.TaskName, &task.TaskDescription, &task.Solution,
		&task.Hints.Hint1, &task.Hints.Hint2, &task.Hints.Hint3, &task.Difficulty, &task.Solved)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (t *TaskRepo) UpdateTaskSolved(ctx context.Context, task *model.GeneratedTask, stats *model.Statistics) error {
	tx, err := t.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			_ = tx.Commit(ctx)
		}
	}()

	query, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Update("tasks").
		Set("solved", task.Solved).
		Where(sq.Eq{"user_id": task.UserID, "task": task.TaskName}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("tx.Exec: %w", err)
	}

	query, args, err = sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Update("statistics").
		Set("easy", stats.Easy).
		Set("medium", stats.Medium).
		Set("hard", stats.Hard).
		Set("total", stats.Total).
		Where(sq.Eq{"user_id": task.UserID}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("tx.Exec: %w", err)
	}
	return nil
}
