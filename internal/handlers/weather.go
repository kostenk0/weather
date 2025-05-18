package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
	"weather/internal/db"
)

type WeatherHandler struct {
	WeatherRepo *db.WeatherRepository
}

func NewWeatherHandler(repo *db.WeatherRepository) *WeatherHandler {
	return &WeatherHandler{WeatherRepo: repo}
}

func (h *WeatherHandler) GetCachedWeather(c *gin.Context) {
	city := strings.TrimSpace(c.Query("city"))
	if city == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "city query parameter is required"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	weather, err := h.WeatherRepo.GetCached(ctx, city)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
		return
	}

	if weather == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "weather data not found for this city"})
		return
	}

	c.JSON(http.StatusOK, weather)
}
