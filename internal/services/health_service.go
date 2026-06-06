package services

import (
	"context"
	"time"

	"github.com/ssjoy/website-health-checker/internal/database"
	"github.com/ssjoy/website-health-checker/internal/redis"
)

// HealthStatus matches the JSON structure of /health.
type HealthStatus struct {
	Status    string `json:"status"`
	Service   string `json:"service"`
	MySQL     string `json:"mysql"`
	Redis     string `json:"redis"`
	Timestamp string `json:"timestamp"`
}

// DetailedHealthStatus matches the JSON structure of /health/details.
type DetailedHealthStatus struct {
	API   string      `json:"api"`
	MySQL DBStatus    `json:"mysql"`
	Redis CacheStatus `json:"redis"`
	Uptime string     `json:"uptime"`
}

type DBStatus struct {
	Status string `json:"status"`
}

type CacheStatus struct {
	Status string `json:"status"`
}

// HealthService defines the business logic operations for health checks.
type HealthService interface {
	GetHealth(ctx context.Context) *HealthStatus
	GetDetailedHealth(ctx context.Context) *DetailedHealthStatus
}

type healthService struct {
	db        *database.DB
	rdb       *redis.Client
	startTime time.Time
}

// NewHealthService creates a new HealthService instance.
func NewHealthService(db *database.DB, rdb *redis.Client, startTime time.Time) HealthService {
	return &healthService{
		db:        db,
		rdb:       rdb,
		startTime: startTime,
	}
}

// GetHealth runs ping checks on MySQL and Redis.
func (s *healthService) GetHealth(ctx context.Context) *HealthStatus {
	mysqlConnected := "connected"
	if err := s.db.PingContext(ctx); err != nil {
		mysqlConnected = "disconnected"
	}

	redisConnected := "connected"
	if err := s.rdb.CheckHealth(ctx); err != nil {
		redisConnected = "disconnected"
	}

	status := "ok"
	if mysqlConnected == "disconnected" || redisConnected == "disconnected" {
		status = "error"
	}

	return &HealthStatus{
		Status:    status,
		Service:   "health-checker",
		MySQL:     mysqlConnected,
		Redis:     redisConnected,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

// GetDetailedHealth checks detailed states of services.
func (s *healthService) GetDetailedHealth(ctx context.Context) *DetailedHealthStatus {
	dbHealth := "healthy"
	if err := s.db.PingContext(ctx); err != nil {
		dbHealth = "unhealthy"
	}

	redisHealth := "healthy"
	if err := s.rdb.CheckHealth(ctx); err != nil {
		redisHealth = "unhealthy"
	}

	// Calculate uptime
	uptime := time.Since(s.startTime).Truncate(time.Second)

	return &DetailedHealthStatus{
		API: "healthy",
		MySQL: DBStatus{
			Status: dbHealth,
		},
		Redis: CacheStatus{
			Status: redisHealth,
		},
		Uptime: uptime.String(),
	}
}
