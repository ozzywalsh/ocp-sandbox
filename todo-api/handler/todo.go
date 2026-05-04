package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	todosGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "todo_items_total",
		Help: "Current number of TODO items",
	})

	requestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "todo_request_duration_seconds",
		Help:    "Duration of TODO API requests",
		Buckets: prometheus.DefBuckets,
	}, []string{"method"})
)

type Todo struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"createdAt"`
}

type Handler struct {
	mu    sync.RWMutex
	todos map[string]Todo
	seq   int
}

func New() *Handler {
	return &Handler{todos: make(map[string]Todo)}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	timer := prometheus.NewTimer(requestDuration.WithLabelValues("list"))
	defer timer.ObserveDuration()

	h.mu.RLock()
	items := make([]Todo, 0, len(h.todos))
	for _, t := range h.todos {
		items = append(items, t)
	}
	h.mu.RUnlock()

	slog.InfoContext(r.Context(), "listing todos", "count", len(items))
	writeJSON(w, http.StatusOK, items)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	timer := prometheus.NewTimer(requestDuration.WithLabelValues("create"))
	defer timer.ObserveDuration()

	var input struct {
		Title string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		slog.WarnContext(r.Context(), "invalid request body", "error", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if input.Title == "" {
		slog.WarnContext(r.Context(), "missing title in create request")
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}

	h.mu.Lock()
	h.seq++
	id := fmt.Sprintf("%d", h.seq)
	todo := Todo{
		ID:        id,
		Title:     input.Title,
		Completed: false,
		CreatedAt: time.Now(),
	}
	h.todos[id] = todo
	h.mu.Unlock()

	todosGauge.Inc()
	slog.InfoContext(r.Context(), "created todo", "id", id, "title", input.Title)
	writeJSON(w, http.StatusCreated, todo)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	timer := prometheus.NewTimer(requestDuration.WithLabelValues("get"))
	defer timer.ObserveDuration()

	id := r.PathValue("id")
	h.mu.RLock()
	todo, ok := h.todos[id]
	h.mu.RUnlock()

	if !ok {
		slog.WarnContext(r.Context(), "todo not found", "id", id)
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, todo)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	timer := prometheus.NewTimer(requestDuration.WithLabelValues("update"))
	defer timer.ObserveDuration()

	id := r.PathValue("id")
	var input struct {
		Title     *string `json:"title"`
		Completed *bool   `json:"completed"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		slog.WarnContext(r.Context(), "invalid request body", "error", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	h.mu.Lock()
	todo, ok := h.todos[id]
	if !ok {
		h.mu.Unlock()
		slog.WarnContext(r.Context(), "todo not found for update", "id", id)
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if input.Title != nil {
		todo.Title = *input.Title
	}
	if input.Completed != nil {
		todo.Completed = *input.Completed
	}
	h.todos[id] = todo
	h.mu.Unlock()

	slog.InfoContext(r.Context(), "updated todo", "id", id)
	writeJSON(w, http.StatusOK, todo)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	timer := prometheus.NewTimer(requestDuration.WithLabelValues("delete"))
	defer timer.ObserveDuration()

	id := r.PathValue("id")
	h.mu.Lock()
	_, ok := h.todos[id]
	if !ok {
		h.mu.Unlock()
		slog.WarnContext(r.Context(), "todo not found for delete", "id", id)
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	delete(h.todos, id)
	h.mu.Unlock()

	todosGauge.Dec()
	slog.InfoContext(r.Context(), "deleted todo", "id", id)
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
