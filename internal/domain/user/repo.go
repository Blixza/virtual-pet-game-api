package domain_user

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, user *Model) error
	Get(ctx context.Context, opts ...Option) (*Model, error)
}

type repo struct {
	db *pgxpool.Pool
	sb sq.StatementBuilderType
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repo{
		db: db,
		sb: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r repo) Create(ctx context.Context, user *Model) error {
	query, args, err := r.sb.Insert("users").
		Columns("email", "username", "password_hash", "created_at").
		Values(user.Email, user.Username, user.PasswordHash, user.CreatedAt).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return err
	}

	err = r.db.QueryRow(ctx, query, args...).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("failed to insert and scan id: %v", err)
	}

	return nil
}

func (r repo) Get(ctx context.Context, opts ...Option) (*Model, error) {
	query := r.sb.Select(
		"id", "email", "username", "password_hash", "created_at", "updated_at",
	).From("users")

	for _, opt := range opts {
		query = query.Where(opt())
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building query: %w", err)
	}

	var u Model
	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&u.ID,           // 1
		&u.Email,        // 2
		&u.Username,     // 3
		&u.PasswordHash, // 4
		&u.CreatedAt,    // 5
		&u.UpdatedAt,    // 6
	)

	if err != nil {
		return nil, err
	}

	return &u, nil
}
