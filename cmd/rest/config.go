package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type config struct {
	httpAddr     string
	databasePath string
	importToken  string
	secureCookie bool
	sessionTTL   time.Duration
}

func loadConfig() (config, error) {
	importToken, err := required("IMPORT_TOKEN")
	if err != nil {
		return config{}, err
	}
	if len(importToken) < 32 {
		return config{}, fmt.Errorf("IMPORT_TOKEN must be at least 32 characters")
	}
	secureCookie, err := boolean("SECURE_COOKIE", true)
	if err != nil {
		return config{}, err
	}
	sessionTTL, err := duration("SESSION_TTL", 30*24*time.Hour)
	if err != nil {
		return config{}, err
	}
	return config{
		httpAddr:     envOr("HTTP_ADDR", ":8080"),
		databasePath: envOr("DATABASE_PATH", "/data/abacus.db"),
		importToken:  importToken,
		secureCookie: secureCookie,
		sessionTTL:   sessionTTL,
	}, nil
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func required(key string) (string, error) {
	v := os.Getenv(key)
	if v == "" {
		return "", fmt.Errorf("%s is required", key)
	}
	return v, nil
}

func boolean(key string, fallback bool) (bool, error) {
	v := os.Getenv(key)
	if v == "" {
		return fallback, nil
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return false, fmt.Errorf("%s must be true or false", key)
	}
	return b, nil
}

func duration(key string, fallback time.Duration) (time.Duration, error) {
	v := os.Getenv(key)
	if v == "" {
		return fallback, nil
	}
	d, err := time.ParseDuration(v)
	if err != nil || d <= 0 {
		return 0, fmt.Errorf("%s must be a positive duration such as 720h", key)
	}
	return d, nil
}
