package handlers

import (
	"encoding/json"
	"minicloudstack/internal/state"
	"net/http"
)

type HealthHandler struct {
	state state.State
}

type HealthResponse struct {
	Status string `json:"status"`
}

type ReadyStatus string

var (
	Ready    ReadyStatus = "ready"
	NotReady ReadyStatus = "not_ready"
)

func NewHealthHandler(state state.State) *HealthHandler {
	return &HealthHandler{state: state}
}

func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	h.Live(w, r)
}

func (h *HealthHandler) Live(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
	ready, err := h.state.Ready(r.Context())
	if err != nil || !ready {
		w.WriteHeader(http.StatusServiceUnavailable)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"status": "not_ready",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status": "ready",
	})
}
