package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
	"weather/internal/db"
	"weather/internal/handlers"
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
