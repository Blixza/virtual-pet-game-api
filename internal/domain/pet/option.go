package domain_pet

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type Option func() sq.Eq

func WithName(name string) Option {
	return func() sq.Eq {
		return sq.Eq{"name": name}
	}
}

func WithID(id uuid.UUID) Option {
	return func() sq.Eq {
		return sq.Eq{"id": id}
	}
}

func WithOwnerID(id uuid.UUID) Option {
	return func() sq.Eq {
		return sq.Eq{"owner_id": id}
	}
}