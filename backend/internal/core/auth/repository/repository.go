package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"study_buddy/internal/model"
	"study_buddy/pkg/constants"
	"study_buddy/pkg/errlist"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ Repository = (*AuthRepo)(nil)

type AuthRepo struct {
	db *pgxpool.Pool
}

func NewAuthRepo(db *pgxpool.Pool) *AuthRepo {
	return &AuthRepo{
		db: db,
	}
}

type Repository interface {
	GetUserByEmailOrUsername(ctx context.Context, username string) (*model.UserData, error)
	CreateUser(ctx context.Context, username, email, password string) (*model.UserData, error)
	UpdateUser(ctx context.Context, user *model.UserData) error
}

func (a *AuthRepo) GetUserByEmailOrUsername(ctx context.Context, username string) (*model.UserData, error) {
	pool, err := a.db.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer pool.Release()

	query, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Select("id", "role", "name", "email", "password",
			"created_at", "updated_at", "is_confirmed").
		From("users").
		Where(
			sq.Or{
				sq.Eq{"name": username},
				sq.Eq{"email": username},
			},
		).ToSql()
	if err != nil {
		return nil, err
	}

	row := pool.QueryRow(ctx, query, args...)
	var user model.UserData

	err = row.Scan(&user.ID, &user.Role, &user.Name, &user.Email, &user.Password,
		&user.CreatedAt, &user.UpdatedAt, &user.IsConfirmed)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errlist.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (a *AuthRepo) CreateUser(ctx context.Context, username, email, password string) (*model.UserData, error) {
	tx, err := a.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			_ = tx.Commit(ctx)
		}
	}()

	createdAt := time.Now()
	userColumns := []string{"role", "name", "email", "password", "created_at", "updated_at", "is_confirmed"}

	userQuery, userArgs, err := sq.Insert("users").
		Columns(userColumns...).
		Values(constants.User, username, email, password, createdAt, createdAt, false).
		Suffix("RETURNING id, " + strings.Join(userColumns, ", ")).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("format user insert SQL: %w", err)
	}

	user := new(model.UserData)
	err = tx.QueryRow(ctx, userQuery, userArgs...).Scan(
		&user.ID, &user.Role, &user.Name, &user.Email, &user.Password,
		&user.CreatedAt, &user.UpdatedAt, &user.IsConfirmed,
	)
	if err != nil {
		return nil, fmt.Errorf("insert user: %w", err)
	}

	// Insert statistics entry
	statsColumns := []string{"user_id", "easy", "medium", "hard", "total"}
	statsQuery, statsArgs, err := sq.Insert("statistics").
		Columns(statsColumns...).
		Values(user.ID, 0, 0, 0, 0).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("format stats insert SQL: %w", err)
	}

	_, err = tx.Exec(ctx, statsQuery, statsArgs...)
	if err != nil {
		return nil, fmt.Errorf("insert stats: %w", err)
	}

	return user, nil
}

func (a *AuthRepo) UpdateUser(ctx context.Context, user *model.UserData) error {
	pool, err := a.db.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("db acquire: %w", err)
	}
	defer pool.Release()

	query, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Update("users").
		Set("is_confirmed", user.IsConfirmed).
		Set("updated_at", user.UpdatedAt).
		Where(sq.Eq{"id": user.ID}).
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
