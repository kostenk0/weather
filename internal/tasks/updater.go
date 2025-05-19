package tasks

import (
	"context"
	"fmt"
	"log"
	"time"
	email "weather/api"
	weatherapi "weather/api"
	"weather/internal/config"
	"weather/internal/db"
)

type WeatherUpdater struct {
	SubscriptionRepo *db.SubscriptionRepository
	WeatherRepo      *db.WeatherRepository
}

func NewWeatherUpdater(subRepo *db.SubscriptionRepository, weatherRepo *db.WeatherRepository) *WeatherUpdater {
	return &WeatherUpdater{
		SubscriptionRepo: subRepo,
		WeatherRepo:      weatherRepo,
	}
}

func (u *WeatherUpdater) SendWeatherFromCacheByFrequency(freq string) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	subs, err := u.SubscriptionRepo.GetConfirmedSubscriptionsByFrequency(ctx, freq)
	if err != nil {
		log.Printf("[CRON] failed to get subscriptions for %s: %v", freq, err)
		return
	}

	var maxAge time.Duration
	switch freq {
	case "hourly":
		maxAge = time.Hour
	case "daily":
		maxAge = 24 * time.Hour
	default:
		maxAge = time.Hour
	}

	baseURL := config.App.AppBaseURL

	for _, sub := range subs {
		weather, err := u.WeatherRepo.GetCached(ctx, sub.City)
		if err != nil {
			log.Printf("[CRON] DB error for %s: %v", sub.City, err)
			continue
		}

		if weather == nil || time.Since(weather.UpdatedAt) > maxAge {
			log.Printf("[CRON] Weather missing or stale for %s, fetching from API...", sub.City)

			weather, err = weatherapi.FetchWeather(sub.City)
			if err != nil {
				log.Printf("[CRON] Failed to fetch weather from API for %s: %v", sub.City, err)
				continue
			}

			if err := u.WeatherRepo.Save(ctx, weather); err != nil {
				log.Printf("[CRON] Failed to save weather for %s: %v", sub.City, err)
			}
		}

		unsubscribeLink := fmt.Sprintf("%s/api/unsubscribe/%s", baseURL, sub.Token)

		subject := fmt.Sprintf("Weather update for %s", weather.City)
		body := fmt.Sprintf(
			"Hello!\n\nHere is your weather update for %s:\n"+
				"Temperature: %.1f°C\nHumidity: %.0f%%\nCondition: %s\n\n"+
				"If you no longer want to receive updates, click here to unsubscribe:\n%s",
			weather.City, weather.Temperature, weather.Humidity, weather.Description, unsubscribeLink,
		)

		if err := email.Send(sub.Email, subject, body); err != nil {
			log.Printf("[EMAIL] Failed to send weather to %s: %v", sub.Email, err)
		} else {
			log.Printf("[EMAIL] Sent weather to %s", sub.Email)
		}
	}
}
