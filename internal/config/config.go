package config

import "os"

type Config struct {
	Port  string
	DBURL string
}

func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	dbURL := getEnv("DATABASE_URL", "playground_user:playground_password@tcp(localhost:1444)/playground?charset=utf8mb4&parseTime=True&loc=Local")
	return &Config{
		Port:  port,
		DBURL: dbURL,
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
