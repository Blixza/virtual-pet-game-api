package service_pet

import (
	"context"
	"database/sql"
	"time"
	domain_pet "virtual_pet_game/internal/domain/pet"

	"github.com/google/uuid"
)

type Service struct {
	repo domain_pet.Repository
}

func New(repo domain_pet.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(
	ctx context.Context, id uuid.UUID,
	name, kind, breed string,
) (*domain_pet.Model, error) {
	pet := &domain_pet.Model{
		Name:           name,
		Kind:           kind,
		Breed:          breed,
		AgeDays:        0,
		Level:          0,
		LastTrainingAt: sql.NullTime{},
		CreatedAt:      time.Now(),
		UpdatedAt:      sql.NullTime{},
	}

	err := s.repo.Create(ctx, pet)
	if err != nil {
		return nil, err
	}

	return pet, nil
}

func (s *Service) Get(
	ctx context.Context, opts ...domain_pet.Option,
) (*domain_pet.Model, error) {
	return s.repo.Get(ctx, opts...)
}

func (s *Service) Update(
	ctx context.Context, pet *domain_pet.Model,
) (*domain_pet.Model, error) {
	pet.UpdatedAt = sql.NullTime{Time: time.Now(), Valid: true}

	updated, err := s.repo.Update(ctx, pet)
	if err != nil {
		return nil, err
	}

	return updated, nil
}
