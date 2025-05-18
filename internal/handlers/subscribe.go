package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
	"weather/internal/db"
	"weather/internal/models"
)

type SubscriptionHandler struct {
	Repo *db.SubscriptionRepository
}

func NewSubscriptionHandler(repo *db.SubscriptionRepository) *SubscriptionHandler {
	return &SubscriptionHandler{Repo: repo}
}

func (h *SubscriptionHandler) Subscribe(c *gin.Context) {
	var req struct {
		Email     string `form:"email" binding:"required,email"`
		City      string `form:"city" binding:"required"`
		Frequency string `form:"frequency" binding:"required,oneof=hourly daily"`
	}

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token := uuid.New().String()

	sub := &models.Subscription{
		Email:     req.Email,
		City:      req.City,
		Frequency: req.Frequency,
		Confirmed: false,
		Token:     token,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := h.Repo.Create(ctx, sub); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Subscription successful. Please confirm via email.",
		"token":   sub.Token,
	})
}
