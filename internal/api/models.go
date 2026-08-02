package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/chankei613/ai-task-router/internal/db"
	"github.com/go-chi/chi/v5"
)

var (
	errModelIDRequired = &apiError{"provider and model_id are required"}
	errModelNotFound   = &apiError{"model not found"}
)

type UpsertModelInput struct {
	Provider         string   `json:"provider"`
	ModelID          string   `json:"model_id"`
	QualityTier      string   `json:"quality_tier"`
	InputPricePer1M  float64  `json:"input_price_per_1m"`
	OutputPricePer1M float64  `json:"output_price_per_1m"`
	Capabilities     []string `json:"capabilities"`
	Enabled          *bool    `json:"enabled"`
}

func (s *Server) UpsertModel(in UpsertModelInput) (db.ModelSpec, error) {
	if in.Provider == "" || in.ModelID == "" {
		return db.ModelSpec{}, errModelIDRequired
	}
	enabled := true
	if in.Enabled != nil {
		enabled = *in.Enabled
	}
	id := in.Provider + ":" + in.ModelID
	now := time.Now()

	spec := db.ModelSpec{
		ID: id, Provider: in.Provider, ModelID: in.ModelID,
		QualityTier: in.QualityTier, InputPricePer1M: in.InputPricePer1M, OutputPricePer1M: in.OutputPricePer1M,
		Capabilities: in.Capabilities, Enabled: enabled, UpdatedAt: now,
	}

	var existing db.ModelSpec
	if err := s.DB.First(&existing, "id = ?", id).Error; err != nil {
		spec.CreatedAt = now
		if err := s.DB.Create(&spec).Error; err != nil {
			return db.ModelSpec{}, err
		}
		return spec, nil
	}

	spec.CreatedAt = existing.CreatedAt
	if err := s.DB.Save(&spec).Error; err != nil {
		return db.ModelSpec{}, err
	}
	return spec, nil
}

func (s *Server) ListModels() ([]db.ModelSpec, error) {
	var models []db.ModelSpec
	err := s.DB.Order("provider asc, model_id asc").Find(&models).Error
	return models, err
}

func (s *Server) SetModelEnabled(id string, enabled bool) (db.ModelSpec, error) {
	var m db.ModelSpec
	if err := s.DB.First(&m, "id = ?", id).Error; err != nil {
		return db.ModelSpec{}, errModelNotFound
	}
	m.Enabled = enabled
	m.UpdatedAt = time.Now()
	if err := s.DB.Save(&m).Error; err != nil {
		return db.ModelSpec{}, err
	}
	return m, nil
}

func (s *Server) DeleteModel(id string) error {
	res := s.DB.Delete(&db.ModelSpec{}, "id = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errModelNotFound
	}
	return nil
}

// ─── HTTP handlers ──────────────────────────────────────────────────────

func (s *Server) httpUpsertModel(w http.ResponseWriter, r *http.Request) {
	var in UpsertModelInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	m, err := s.UpsertModel(in)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusOK, m)
}

func (s *Server) httpListModels(w http.ResponseWriter, r *http.Request) {
	models, err := s.ListModels()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, models)
}

type setEnabledRequest struct {
	Enabled bool `json:"enabled"`
}

func (s *Server) httpSetModelEnabled(w http.ResponseWriter, r *http.Request) {
	var body setEnabledRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	m, err := s.SetModelEnabled(chi.URLParam(r, "id"), body.Enabled)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, m)
}

func (s *Server) httpDeleteModel(w http.ResponseWriter, r *http.Request) {
	if err := s.DeleteModel(chi.URLParam(r, "id")); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
