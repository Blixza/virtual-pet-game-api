package handler_pet

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"
	domain_pet "virtual_pet_game/internal/domain/pet"
	service_pet "virtual_pet_game/internal/service/pet"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type Handler struct {
	svc *service_pet.Service
	log *zap.Logger
}

func New(svc *service_pet.Service, log *zap.Logger) *Handler {
	return &Handler{svc: svc, log: log}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		ID             uuid.UUID    `json:"id"`
		Name           string       `json:"name"`
		Kind           string       `json:"kind"`
		Breed          string       `json:"breed"`
		AgeDays        int          `json:"age_days"`
		Level          int          `json:"level"`
		LastTrainingAt sql.NullTime `json:"last_training_at"`
		CreatedAt      time.Time    `json:"created_at"`
		UpdatedAt      sql.NullTime `json:"updated_at"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.log.Error("failed to decode request", zap.Error(err))
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	petExists, err := h.svc.Get(r.Context(), domain_pet.WithName(req.Name), domain_pet.WithOwnerID(userID))
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			h.log.Error("database error", zap.Error(err))
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
	}

	if petExists != nil {
		h.log.Error("pet with this name already exists", zap.Error(err))
		http.Error(w, "pet with this name already exists", http.StatusBadRequest)
		return
	}

	pet, err := h.svc.Create(r.Context(), userID, req.Name, req.Kind, req.Breed)
	if err != nil {
		h.log.Error("failed to create pet", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	h.log.Info("pet created", zap.Any("pet_id", pet.ID))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pet)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	var opts []domain_pet.Option

	idStr := r.PathValue("id")
	if idStr != "" {
		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, "invalid uuid", http.StatusBadRequest)
			return
		}
		opts = append(opts, domain_pet.WithID(id))
	}

	name := r.PathValue("name")
	if name != "" {
		opts = append(opts, domain_pet.WithName(name))
	}

	pet, err := h.svc.Get(r.Context(), opts...)
	if err != nil {
		h.log.Error("pet not found", zap.Error(err))
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pet)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	petID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid pet id", http.StatusBadRequest)
		return
	}

	var req struct {
		Name  string `json:"name"`
		Kind  string `json:"kind"`
		Breed string `json:"breed"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("failed to decode request", zap.Error(err))
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	petUpdate := domain_pet.Model{
		ID:    petID,
		Name:  req.Name,
		Kind:  req.Kind,
		Breed: req.Breed,
	}

	pet, err := h.svc.Update(r.Context(), &petUpdate)
	if err != nil {
		h.log.Error("failed to update pet", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	h.log.Info("pet updated", zap.Any("pet_id", pet.ID))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pet)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	petID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid pet id", http.StatusBadRequest)
		return
	}

	err = h.svc.Delete(r.Context(), domain_pet.WithID(petID))
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			h.log.Error("database error", zap.Error(err))
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
	}

	h.log.Info("pet deleted", zap.Any("pet_id", petID))
	w.Header().Set("Content-Type", "application/json")
}
