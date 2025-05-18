package handlers

import (
	"context"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func (h *SubscriptionHandler) Unsubscribe(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing token"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := h.Repo.DeleteByToken(ctx, token)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "invalid token"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Unsubscribed successfully"})
}
