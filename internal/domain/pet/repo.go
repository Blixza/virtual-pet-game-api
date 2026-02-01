package domain_pet

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Repository interface {
	Create(ctx context.Context, pet *Model) error
	Get(ctx context.Context, opts ...Option) (*Model, error)
	Update(ctx context.Context, pet *Model) (*Model, error)
}

type repo struct {
	db  *pgxpool.Pool
	log *zap.Logger
	sb  sq.StatementBuilderType
}

func NewRepository(db *pgxpool.Pool, log *zap.Logger) Repository {
	return &repo{
		db:  db,
		log: log,
		sb:  sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r repo) Create(ctx context.Context, pet *Model) error {
	query, args, err := r.sb.Insert("pets").
		SetMap(map[string]interface{}{
			"name":             pet.Name,
			"kind":             pet.Kind,
			"breed":            pet.Breed,
			"age_days":         pet.AgeDays,
			"level":            pet.Level,
			"last_training_at": pet.LastTrainingAt,
			"created_at":       pet.CreatedAt,
			"updated_at":       pet.UpdatedAt,
		}).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return err
	}

	err = r.db.QueryRow(ctx, query, args...).Scan(&pet.ID)
	if err != nil {
		return fmt.Errorf("failed to insert and scan id: %v", err)
	}

	return nil
}

func (r repo) Get(ctx context.Context, opts ...Option) (*Model, error) {
	query := r.sb.Select(
		"id", "name", "kind", "breed", "((EXTRACT(EPOCH FROM (NOW() - created_at)) / 86400)::INT) as age_days",
		"level", "last_training_at", "created_at", "updated_at",
	).From("pets")

	for _, opt := range opts {
		query = opt(query)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building query: %w", err)
	}

	var p Model
	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&p.ID,             // 1
		&p.Name,           // 2
		&p.Kind,           // 3
		&p.Breed,          // 4
		&p.AgeDays,        // 5
		&p.Level,          // 6
		&p.LastTrainingAt, // 7
		&p.CreatedAt,      // 8
		&p.UpdatedAt,      // 9
	)

	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r repo) Update(ctx context.Context, pet *Model) (*Model, error) {
	query, args, err := r.sb.Update("pets").
		SetMap(map[string]interface{}{
			"name":             pet.Name,
			"kind":             pet.Kind,
			"breed":            pet.Breed,
			"age_days":         pet.AgeDays,
			"level":            pet.Level,
			"last_training_at": pet.LastTrainingAt,
			"updated_at":       pet.UpdatedAt,
		}).
		Where(sq.Eq{"id": pet.ID}).
		Suffix("RETURNING updated_at").
		ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(ctx, query, args...).Scan(&pet.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to update and scan id: %v", err)
	}

	return pet, nil
}
