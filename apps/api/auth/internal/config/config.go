package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Port         string
	MongoURI     string
	MongoDB      string
	JWTSecret    []byte
	AccessTTL    time.Duration
	RefreshTTL   time.Duration
	CookieSecure bool
	CORSOrigins  []string
}

func Load() (*Config, error) {
	accessTTL, err := duration("ACCESS_TTL", 15*time.Minute)
	if err != nil {
		return nil, err
	}
	refreshTTL, err := duration("REFRESH_TTL", 30*24*time.Hour)
	if err != nil {
		return nil, err
	}
	cookieSecure, err := strconv.ParseBool(env("COOKIE_SECURE", "true"))
	if err != nil {
		return nil, errors.New("COOKIE_SECURE must be a bool")
	}

	secret := os.Getenv("JWT_SECRET")
	if len(secret) < 32 {
		return nil, errors.New("JWT_SECRET must be at least 32 characters")
	}

	var origins []string
	if v := os.Getenv("CORS_ORIGINS"); v != "" {
		origins = strings.Split(v, ",")
	}

	return &Config{
		Port:         env("PORT", "8080"),
		MongoURI:     env("MONGO_URI", "mongodb://localhost:27018"),
		MongoDB:      env("MONGO_DB", "to_on"),
		JWTSecret:    []byte(secret),
		AccessTTL:    accessTTL,
		RefreshTTL:   refreshTTL,
		CookieSecure: cookieSecure,
		CORSOrigins:  origins,
	}, nil
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func duration(key string, fallback time.Duration) (time.Duration, error) {
	v := os.Getenv(key)
	if v == "" {
		return fallback, nil
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return 0, errors.New(key + " must be a duration (e.g. 15m)")
	}
	return d, nil
}
