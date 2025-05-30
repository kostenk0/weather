package handlers

import (
	"context"
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	apierr "weather/internal/errors"
)

func (h *SubscriptionHandler) ConfirmSubscription(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		apierr.Respond(c, apierr.New(http.StatusBadRequest, "Missing token"))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := h.Repo.ConfirmByToken(ctx, token)
	if errors.Is(err, sql.ErrNoRows) {
		apierr.Respond(c, apierr.New(http.StatusNotFound, "Invalid or already confirmed token"))
		return
	} else if err != nil {
		apierr.Respond(c, apierr.New(http.StatusInternalServerError, "Failed to confirm subscription"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subscription confirmed successfully"})
}
