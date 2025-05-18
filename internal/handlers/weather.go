package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
	"weather/internal/db"
	apierr "weather/internal/errors"
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
		apierr.Respond(c, apierr.New(http.StatusBadRequest, "City query parameter is required"))
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	weather, err := h.WeatherRepo.GetCached(ctx, city)
	if err != nil {
		apierr.Respond(c, apierr.New(http.StatusInternalServerError, "Failed to get weather"))
		return
	}

	if weather == nil {
		apierr.Respond(c, apierr.New(http.StatusNotFound, "Weather data not found for this city"))
		return
	}

	c.JSON(http.StatusOK, weather)
}
