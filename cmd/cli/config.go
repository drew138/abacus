package main

import "os"

type config struct {
	databasePath string
}

func loadConfig() config {
	return config{databasePath: envOr("DATABASE_PATH", "/data/abacus.db")}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
