package handlers

import (
	"context"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	apierr "weather/internal/errors"
)

func (h *SubscriptionHandler) Unsubscribe(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		apierr.Respond(c, apierr.New(http.StatusBadRequest, "Missing token"))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := h.Repo.DeleteByToken(ctx, token)
	if err == sql.ErrNoRows {
		apierr.Respond(c, apierr.New(http.StatusNotFound, "Token not found"))
		return
	} else if err != nil {
		apierr.Respond(c, apierr.New(http.StatusInternalServerError, "Failed to unsubscribe"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Unsubscribed successfully"})
}
