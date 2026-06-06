package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/ssjoy/website-health-checker/internal/services"
)

// CacheHandler handles caching-related requests.
type CacheHandler struct {
	service services.CacheService
}

// NewCacheHandler returns a new CacheHandler instance.
func NewCacheHandler(service services.CacheService) *CacheHandler {
	return &CacheHandler{service: service}
}

// CacheRequest represents the payload for POST /cache.
type CacheRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// HandleCache processes GET /cache and POST /cache requests.
func (h *CacheHandler) HandleCache(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		h.getCache(w, r)
	case http.MethodPost:
		h.setCache(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte(`{"error":"Method Not Allowed"}`))
	}
}

func (h *CacheHandler) getCache(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"key query parameter is required"}`))
		return
	}

	val, err := h.service.Get(r.Context(), key)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error":"key not found in cache"}`))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"key":   key,
		"value": val,
	})
}

func (h *CacheHandler) setCache(w http.ResponseWriter, r *http.Request) {
	var req CacheRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"invalid JSON body"}`))
		return
	}

	if req.Key == "" || req.Value == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"key and value are required fields"}`))
		return
	}

	// Store in cache with 24 hours default TTL (or pass 0 for no expiration)
	err := h.service.Set(r.Context(), req.Key, req.Value, 24*time.Hour)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"message":"cached successfully"}`))
}
