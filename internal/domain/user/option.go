package domain_user

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type Option func() sq.Eq

func WithEmail(email string) Option {
	return func() sq.Eq {
		return sq.Eq{"email": email}
	}
}

func WithUsername(username string) Option {
	return func() sq.Eq {
		return sq.Eq{"username": username}
	}
}

func WithID(id uuid.UUID) Option {
	return func() sq.Eq {
		return sq.Eq{"id": id}
	}
}
