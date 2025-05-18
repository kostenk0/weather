package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"log"
	"os"
	"weather/internal/db"
	"weather/internal/handlers"
	"weather/internal/tasks"
)

func main() {
	_ = godotenv.Load()

	database, err := db.Connect()
	if err != nil {
		log.Fatalf("DB error: %v", err)
	}
	defer func(database *sql.DB) {
		err := database.Close()
		if err != nil {
			log.Fatalf("DB close: %v", err)
		}
	}(database)

	repo := db.NewSubscriptionRepository(database)
	handler := handlers.NewSubscriptionHandler(repo)

	weatherRepo := db.NewWeatherRepository(database)
	weatherHandler := handlers.NewWeatherHandler(weatherRepo)

	weatherUpdater := tasks.NewWeatherUpdater(repo, weatherRepo)

	c := cron.New()

	c.AddFunc("@every 1m", func() {
		weatherUpdater.SendWeatherFromCacheByFrequency("hourly")
	})

	c.AddFunc("@every 24h", func() {
		weatherUpdater.SendWeatherFromCacheByFrequency("daily")
	})

	c.Start()

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.POST("/api/subscribe", handler.Subscribe)

	r.GET("/api/confirm/:token", handler.ConfirmSubscription)

	r.GET("/api/unsubscribe/:token", handler.Unsubscribe)

	r.GET("/api/weather", weatherHandler.GetCachedWeather)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
