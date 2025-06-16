package handlers

import (
	"Daemon/internal/container"
	_ "Daemon/internal/models"
	_ "context"
	"encoding/json"
	"net/http"
)

type ContainerHandler struct {
	Service *container.Service
}

type CreateRequest struct {
	Name    string `json:"name"`
	EggName string `json:"egg"`
}

// POST /containers
func (h *ContainerHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	container, err := h.Service.CreateContainer(r.Context(), req.Name, req.EggName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, container)
}

// GET /containers/:name
func (h *ContainerHandler) Get(w http.ResponseWriter, r *http.Request, name string) {
	container, err := h.Service.GetContainer(r.Context(), name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, container)
}

// DELETE /containers/:name
func (h *ContainerHandler) Delete(w http.ResponseWriter, r *http.Request, name string) {
	err := h.Service.RemoveContainer(r.Context(), name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// POST /containers/:name/start
func (h *ContainerHandler) Start(w http.ResponseWriter, r *http.Request, name string) {
	err := h.Service.StartContainer(r.Context(), name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// POST /containers/:name/stop
func (h *ContainerHandler) Stop(w http.ResponseWriter, r *http.Request, name string) {
	err := h.Service.StopContainer(r.Context(), name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}
