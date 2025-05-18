package tasks

import (
	"context"
	"log"
	"time"
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	subs, err := u.SubscriptionRepo.GetConfirmedSubscriptionsByFrequency(ctx, freq)
	if err != nil {
		log.Printf("[CRON] failed to get subscriptions for %s: %v", freq, err)
		return
	}

	for _, sub := range subs {
		weather, err := u.WeatherRepo.GetCached(ctx, sub.City)
		if err != nil || weather == nil {
			log.Printf("[CRON] no cached weather for %s", sub.City)
			continue
		}

		log.Printf("[CRON] Would send email to %s: %s — %.1f°C, %s",
			sub.Email, weather.City, weather.Temperature, weather.Description)
	}
}
