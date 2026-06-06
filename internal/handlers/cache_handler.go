package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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
	Key   string `json:"key" binding:"required"`
	Value string `json:"value" binding:"required"`
}

// CacheResponse represents the payload returned by GET /cache.
type CacheResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// MessageResponse represents simple textual messages.
type MessageResponse struct {
	Message string `json:"message"`
}

// ErrorResponse represents an error payload.
type ErrorResponse struct {
	Error string `json:"error"`
}

// GetCache handles GET /cache requests.
// @Summary      Get cached item
// @Description  Retrieves the value associated with a key from Redis.
// @Tags         cache
// @Produce      json
// @Param        key  query     string  true  "Cache Key"
// @Success      200  {object}  CacheResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /cache [get]
func (h *CacheHandler) GetCache(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "key query parameter is required"})
		return
	}

	val, err := h.service.Get(c.Request.Context(), key)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "key not found in cache"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, CacheResponse{
		Key:   key,
		Value: val,
	})
}

// SetCache handles POST /cache requests.
// @Summary      Save cache item
// @Description  Stores a key-value pair in Redis with a 24-hour default expiration.
// @Tags         cache
// @Accept       json
// @Produce      json
// @Param        body  body      CacheRequest  true  "Cache Entry Details"
// @Success      200   {object}  MessageResponse
// @Failure      400   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Router       /cache [post]
func (h *CacheHandler) SetCache(c *gin.Context) {
	var req CacheRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "key and value are required fields"})
		return
	}

	err := h.service.Set(c.Request.Context(), req.Key, req.Value, 24*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, MessageResponse{Message: "cached successfully"})
}
