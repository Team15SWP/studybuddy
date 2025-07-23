package repository

import (
	"context"

	"study_buddy/internal/model"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ StatsRepository = (*StatsRepo)(nil)

type StatsRepo struct {
	db *pgxpool.Pool
}

func NewStatsRepo(db *pgxpool.Pool) *StatsRepo {
	return &StatsRepo{
		db: db,
	}
}

type StatsRepository interface {
	GetStatisticsData(ctx context.Context, userId int64) (*model.Statistics, error)
}

func (s *StatsRepo) GetStatisticsData(ctx context.Context, userId int64) (*model.Statistics, error) {
	pool, err := s.db.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer pool.Release()

	query, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Select("easy", "medium", "hard", "total").
		From("statistics").
		Where(sq.Eq{"user_id": userId}).
		ToSql()
	if err != nil {
		return nil, err
	}

	row := pool.QueryRow(ctx, query, args...)
	var stats model.Statistics

	err = row.Scan(&stats.Easy, &stats.Medium, &stats.Hard, &stats.Total)
	if err != nil {
		return nil, err
	}
	return &stats, nil
}
