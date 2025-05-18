package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strings"
	"time"
	"weather/internal/db"
	error "weather/internal/errors"
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
		error.Respond(c, error.New(http.StatusBadRequest, "Invalid input"))
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
		if strings.Contains(err.Error(), "Duplicate key") {
			error.Respond(c, error.New(http.StatusConflict, "Email already subscribed"))
		} else {
			error.Respond(c, error.New(http.StatusInternalServerError, "Failed to create subscription"))
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Subscription successful. Please confirm via email.",
		"token":   sub.Token,
	})
}
