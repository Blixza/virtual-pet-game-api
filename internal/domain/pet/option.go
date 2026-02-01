package domain_pet

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type Option func(sq.SelectBuilder) sq.SelectBuilder

func WithName(name string) Option {
	return func(sb sq.SelectBuilder) sq.SelectBuilder {
		return sb.Where(sq.Eq{"name": name})
	}
}

func WithID(id uuid.UUID) Option {
	return func(sb sq.SelectBuilder) sq.SelectBuilder {
		return sb.Where(sq.Eq{"id": id})
	}
}
