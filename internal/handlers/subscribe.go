package handlers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"strings"
	"time"
	email "weather/api"
	"weather/internal/config"
	"weather/internal/db"
	apierr "weather/internal/errors"
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
		apierr.Respond(c, apierr.New(http.StatusBadRequest, "Invalid input"))
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
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "Duplicate key") {
			apierr.Respond(c, apierr.New(http.StatusConflict, "Email already subscribed"))
		} else {
			apierr.Respond(c, apierr.New(http.StatusInternalServerError, "Failed to create subscription"))
		}
		return
	}

	baseURL := config.App.AppBaseURL

	confirmLink := fmt.Sprintf("%s/api/confirm/%s", baseURL, sub.Token)

	subject := "Confirm your subscription"
	body := fmt.Sprintf(
		"Hello!\n\nPlease confirm your weather subscription by clicking the link below:\n\n%s\n\nIf you didn't request this, just ignore this email.",
		confirmLink,
	)

	if err := email.Send(sub.Email, subject, body); err != nil {
		log.Printf("[EMAIL] Failed to send confirmation to %s: %v", sub.Email, err)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Subscription successful. Please confirm via email.",
		"token":   sub.Token,
	})
}
