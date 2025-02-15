package config

import (
	"os"
	"strings"
)

var (
	cfg = LoadConfig()
)

type Config struct {
	JWTSecret   string
	CorsOrigins []string
}

func LoadConfig() Config {

	// Get the CORS_ORIGIN env var
	corsOrigins := os.Getenv("CORS_ORIGIN")
	if corsOrigins == "" {
		corsOrigins = "*"
	}

	// Split the origins by comma
	origins := strings.Split(corsOrigins, ",")
	// Remove any empty strings and strip whitespace
	for i := 0; i < len(origins); i++ {
		origins[i] = strings.TrimSpace(origins[i])
		if origins[i] == "" {
			origins = append(origins[:i], origins[i+1:]...)
			i--
		}
	}

	// Print the origins
	for _, origin := range origins {
		println(origin)
	}

	return Config{
		JWTSecret:   os.Getenv("JWT_SECRET"),
		CorsOrigins: origins,
	}
}
