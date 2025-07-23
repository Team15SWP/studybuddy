package repository

import (
	"context"
	"fmt"

	"study_buddy/internal/model"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ Repository = (*SyllabusRepo)(nil)

type SyllabusRepo struct {
	db *pgxpool.Pool
}

func NewSyllabusRepo(db *pgxpool.Pool) *SyllabusRepo {
	return &SyllabusRepo{
		db: db,
	}
}

type Repository interface {
	GetSyllabus(ctx context.Context) ([]string, error)
	SaveSyllabus(ctx context.Context, syllabus []model.Schedule) ([]string, error)
	DeleteSyllabus(ctx context.Context) error
}

func (s *SyllabusRepo) GetSyllabus(ctx context.Context) ([]string, error) {
	pool, err := s.db.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("s.db.Acquire: %w", err)
	}
	defer pool.Release()

	query, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Select("topic").
		From("syllabus").
		OrderBy("id").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("query builder to sql: %w", err)
	}

	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("pool.Query: %w", err)
	}
	defer rows.Close()

	topics := make([]string, 0)
	for rows.Next() {
		var topic string
		if err := rows.Scan(&topic); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}
		topics = append(topics, topic)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err: %w", err)
	}

	return topics, nil
}

func (s *SyllabusRepo) DeleteSyllabus(ctx context.Context) error {
	pool, err := s.db.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("s.db.Acquire: %w", err)
	}
	defer pool.Release()

	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("transaction begin: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, "DELETE FROM syllabus")
	if err != nil {
		return fmt.Errorf("tx.Exec: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("tx.Commit: %w", err)
	}
	return nil
}

func (s *SyllabusRepo) SaveSyllabus(ctx context.Context, syllabus []model.Schedule) ([]string, error) {
	pool, err := s.db.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("s.db.Acquire: %w", err)
	}
	defer pool.Release()

	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("transaction begin: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, "DELETE FROM syllabus")
	if err != nil {
		return nil, err
	}

	insertBuilder := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Insert("syllabus").
		Columns("week", "topic")

	for _, item := range syllabus {
		insertBuilder = insertBuilder.Values(item.Week, item.Topic)
	}

	query, args, err := insertBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("query builder to sql: %w", err)
	}

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("tx.Exec: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("tx.Commit: %w", err)
	}

	topics := make([]string, 0)
	for _, item := range syllabus {
		topics = append(topics, item.Topic)
	}
	return topics, nil
}
