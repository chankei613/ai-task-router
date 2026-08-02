package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/chankei613/ai-task-router/internal/db"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

var errRouteNotFound = &apiError{"routing log not found"}

type RouteInput struct {
	Source               string   `json:"source"`
	TaskType             string   `json:"task_type"`
	RequiredCapabilities []string `json:"required_capabilities"`
	MinQualityTier       string   `json:"min_quality_tier"`
}

type RouteResult struct {
	Log         db.RoutingLog `json:"log"`
	ChosenModel *db.ModelSpec `json:"chosen_model,omitempty"`
}

// Route はModelSpec全件を候補として決定論的ルーティングを行い、RoutingLogとして永続化する。
func (s *Server) Route(in RouteInput) (RouteResult, error) {
	var candidates []db.ModelSpec
	if err := s.DB.Find(&candidates).Error; err != nil {
		return RouteResult{}, err
	}

	chosen, reasoning := decideRoute(candidates, in.RequiredCapabilities, in.MinQualityTier)

	log := db.RoutingLog{
		ID:                   uuid.NewString(),
		Source:               in.Source,
		TaskType:             in.TaskType,
		RequiredCapabilities: in.RequiredCapabilities,
		MinQualityTier:       in.MinQualityTier,
		Reasoning:            reasoning,
		RequestedAt:          time.Now(),
	}
	if chosen != nil {
		log.ChosenModelID = chosen.ID
	}

	if err := s.DB.Create(&log).Error; err != nil {
		return RouteResult{}, err
	}

	return RouteResult{Log: log, ChosenModel: chosen}, nil
}

func (s *Server) ListRoutingLogs() ([]db.RoutingLog, error) {
	var logs []db.RoutingLog
	err := s.DB.Order("requested_at asc").Find(&logs).Error
	return logs, err
}

func (s *Server) GetRoutingLog(id string) (RouteResult, error) {
	var log db.RoutingLog
	if err := s.DB.First(&log, "id = ?", id).Error; err != nil {
		return RouteResult{}, errRouteNotFound
	}
	result := RouteResult{Log: log}
	if log.ChosenModelID != "" {
		var m db.ModelSpec
		if err := s.DB.First(&m, "id = ?", log.ChosenModelID).Error; err == nil {
			result.ChosenModel = &m
		}
	}
	return result, nil
}

// ─── HTTP handlers ──────────────────────────────────────────────────────

func (s *Server) httpRoute(w http.ResponseWriter, r *http.Request) {
	var in RouteInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	result, err := s.Route(in)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (s *Server) httpListRoutingLogs(w http.ResponseWriter, r *http.Request) {
	logs, err := s.ListRoutingLogs()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, logs)
}

func (s *Server) httpGetRoutingLog(w http.ResponseWriter, r *http.Request) {
	result, err := s.GetRoutingLog(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, result)
}
