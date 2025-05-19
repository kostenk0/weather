package config

import "os"

type Config struct {
	AppBaseURL string
	Port       string
	Email      EmailConfig
	External   ExternalAPIConfig
}

type EmailConfig struct {
	SMTPHost  string
	SMTPPort  string
	SMTPUser  string
	SMTPPass  string
	EmailFrom string
}

type ExternalAPIConfig struct {
	WeatherAPIKey string
}

var App Config

func Init() {
	App = Config{
		AppBaseURL: os.Getenv("APP_BASE_URL"),
		Port:       getEnv("PORT", "8080"),
		Email: EmailConfig{
			SMTPHost:  os.Getenv("SMTP_HOST"),
			SMTPPort:  os.Getenv("SMTP_PORT"),
			SMTPUser:  os.Getenv("SMTP_USER"),
			SMTPPass:  os.Getenv("SMTP_PASS"),
			EmailFrom: os.Getenv("EMAIL_FROM"),
		},
		External: ExternalAPIConfig{
			WeatherAPIKey: os.Getenv("WEATHER_API_KEY"),
		},
	}
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}
