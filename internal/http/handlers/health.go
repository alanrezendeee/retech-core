package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type HealthHandler struct {
	Mongo *mongo.Client
}

func NewHealthHandler(m *mongo.Client) *HealthHandler {
	return &HealthHandler{Mongo: m}
}

func (h *HealthHandler) Get(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 2*time.Second)
	defer cancel()
	err := h.Mongo.Ping(ctx, nil)
	status := "ok"
	if err != nil {
		status = "degraded"
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"status":  status,
		"mongo":   status,
		"ts":      time.Now().UTC(),
	})
}

