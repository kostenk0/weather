package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
	weatherapi "weather/api"
	"weather/internal/db"
	apierr "weather/internal/errors"
)

type WeatherHandler struct {
	WeatherRepo *db.WeatherRepository
}

func NewWeatherHandler(repo *db.WeatherRepository) *WeatherHandler {
	return &WeatherHandler{WeatherRepo: repo}
}

func (h *WeatherHandler) GetWeather(c *gin.Context) {
	city := strings.TrimSpace(c.Query("city"))
	if city == "" {
		apierr.Respond(c, apierr.New(http.StatusBadRequest, "City query parameter is required"))
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	const maxAge = time.Hour

	weather, err := h.WeatherRepo.GetCached(ctx, city)
	if err != nil {
		apierr.Respond(c, apierr.New(http.StatusInternalServerError, "Failed to get weather from DB"))
		return
	}

	if weather == nil || time.Since(weather.UpdatedAt) > maxAge {
		newWeather, fetchErr := weatherapi.FetchWeather(city)
		if fetchErr != nil {
			if weather == nil {
				apierr.Respond(c, apierr.New(http.StatusNotFound, "Weather data not found and fetch failed"))
				return
			}
			c.JSON(http.StatusOK, weather)
			return
		}

		if saveErr := h.WeatherRepo.Save(ctx, newWeather); saveErr != nil {
			c.JSON(http.StatusOK, newWeather)
			return
		}

		c.JSON(http.StatusOK, newWeather)
		return
	}

	c.JSON(http.StatusOK, weather)
}
