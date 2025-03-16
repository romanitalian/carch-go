package http

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/romanitalian/carch-go/internal/domain"
	"github.com/romanitalian/carch-go/internal/pkg/logger"
	"github.com/romanitalian/carch-go/internal/service"
)

type Handler struct {
	services *service.Services
	log      *logger.Logger
	mux      *http.ServeMux
}

func NewHandler(services *service.Services, log *logger.Logger) *Handler {
	h := &Handler{
		services: services,
		log:      log,
		mux:      http.NewServeMux(),
	}

	h.setupRoutes()
	return h
}

func (h *Handler) setupRoutes() {
	// REST API endpoints
	h.mux.HandleFunc("POST /api/v1/users", h.logRequest(h.createUser))
	h.mux.HandleFunc("GET /api/v1/users/{id}", h.logRequest(h.getUserByID))
	h.mux.HandleFunc("PUT /api/v1/users/{id}", h.logRequest(h.updateUser))
	h.mux.HandleFunc("DELETE /api/v1/users/{id}", h.logRequest(h.deleteUser))
	h.mux.HandleFunc("GET /api/v1/users", h.logRequest(h.listUsers))
}

// ServeHTTP implements the http.Handler interface
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

// Middleware for logging requests
func (h *Handler) logRequest(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a response wrapper to capture status code
		rw := newResponseWriter(w)

		// Process request
		next(rw, r)

		// Log after request is processed
		h.log.Info("HTTP Request", map[string]interface{}{
			"method":      r.Method,
			"path":        r.URL.Path,
			"status":      rw.statusCode,
			"duration_ms": time.Since(start).Milliseconds(),
			"remote_addr": r.RemoteAddr,
			"user_agent":  r.UserAgent(),
		})
	}
}

// responseWriter is a wrapper for http.ResponseWriter that captures the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Helper functions for handling requests and responses
func (h *Handler) decodeJSONBody(r *http.Request, dst interface{}) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.Unmarshal(body, dst)
}

func (h *Handler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			h.log.Error("Failed to encode response", err, nil)
		}
	}
}

func (h *Handler) respondError(w http.ResponseWriter, status int, err error) {
	h.respondJSON(w, status, errorRS{Error: err.Error()})
}

// For testing purposes
var pathValueFunc = func(r *http.Request, key string) string {
	return r.PathValue(key)
}

// Handler functions
func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	var req createUserRQ
	if err := h.decodeJSONBody(r, &req); err != nil {
		h.log.Error("Failed to decode request body", err, map[string]interface{}{"path": r.URL.Path})
		h.respondError(w, http.StatusBadRequest, err)
		return
	}

	// Validate required fields
	if req.Email == "" || req.Password == "" {
		h.log.Warn("Missing required fields", map[string]interface{}{"path": r.URL.Path})
		h.respondError(w, http.StatusBadRequest, domain.ErrInvalidInput)
		return
	}

	// Validate email format
	if !isValidEmail(req.Email) {
		h.log.Warn("Invalid email format", map[string]interface{}{"email": req.Email})
		h.respondError(w, http.StatusBadRequest, domain.ErrInvalidInput)
		return
	}

	user := &domain.User{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	}

	if err := h.services.User.Create(r.Context(), user); err != nil {
		h.log.Error("Failed to create user", err, map[string]interface{}{"email": req.Email})
		h.respondError(w, http.StatusInternalServerError, err)
		return
	}

	h.respondJSON(w, http.StatusCreated, user)
}

// isValidEmail checks if the email has a valid format
func isValidEmail(email string) bool {
	// Simple validation: check if it contains @ and a period after @
	atIndex := strings.Index(email, "@")
	if atIndex < 1 {
		return false
	}

	dotIndex := strings.LastIndex(email, ".")
	return dotIndex > atIndex && dotIndex < len(email)-1
}

func (h *Handler) getUserByID(w http.ResponseWriter, r *http.Request) {
	id := pathValueFunc(r, "id")
	if id == "" {
		h.log.Warn("Missing user ID", map[string]interface{}{"path": r.URL.Path})
		h.respondError(w, http.StatusBadRequest, domain.ErrInvalidInput)
		return
	}

	user, err := h.services.User.GetByID(r.Context(), id)
	if err != nil {
		if err == domain.ErrUserNotFound {
			h.log.Warn("User not found", map[string]interface{}{"user_id": id})
			h.respondError(w, http.StatusNotFound, err)
			return
		}
		h.log.Error("Failed to get user", err, map[string]interface{}{"user_id": id})
		h.respondError(w, http.StatusInternalServerError, err)
		return
	}

	h.respondJSON(w, http.StatusOK, user)
}

func (h *Handler) updateUser(w http.ResponseWriter, r *http.Request) {
	id := pathValueFunc(r, "id")
	if id == "" {
		h.log.Warn("Missing user ID", map[string]interface{}{"path": r.URL.Path})
		h.respondError(w, http.StatusBadRequest, domain.ErrInvalidInput)
		return
	}

	var req updateUserRQ
	if err := h.decodeJSONBody(r, &req); err != nil {
		h.log.Error("Failed to decode request body", err, map[string]interface{}{"path": r.URL.Path})
		h.respondError(w, http.StatusBadRequest, err)
		return
	}

	user := &domain.User{
		ID:    id,
		Email: req.Email,
		Name:  req.Name,
	}

	if err := h.services.User.Update(r.Context(), user); err != nil {
		if err == domain.ErrUserNotFound {
			h.log.Warn("User not found for update", map[string]interface{}{"user_id": id})
			h.respondError(w, http.StatusNotFound, err)
			return
		}
		h.log.Error("Failed to update user", err, map[string]interface{}{"user_id": id})
		h.respondError(w, http.StatusInternalServerError, err)
		return
	}

	h.respondJSON(w, http.StatusOK, user)
}

func (h *Handler) deleteUser(w http.ResponseWriter, r *http.Request) {
	id := pathValueFunc(r, "id")
	if id == "" {
		h.log.Warn("Missing user ID", map[string]interface{}{"path": r.URL.Path})
		h.respondError(w, http.StatusBadRequest, domain.ErrInvalidInput)
		return
	}

	if err := h.services.User.Delete(r.Context(), id); err != nil {
		if err == domain.ErrUserNotFound {
			h.log.Warn("User not found for deletion", map[string]interface{}{"user_id": id})
			h.respondError(w, http.StatusNotFound, err)
			return
		}
		h.log.Error("Failed to delete user", err, map[string]interface{}{"user_id": id})
		h.respondError(w, http.StatusInternalServerError, err)
		return
	}

	h.respondJSON(w, http.StatusNoContent, nil)
}

func (h *Handler) listUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.services.User.List(r.Context())
	if err != nil {
		h.log.Error("Failed to list users", err, nil)
		h.respondError(w, http.StatusInternalServerError, err)
		return
	}

	h.respondJSON(w, http.StatusOK, users)
}
