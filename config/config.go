package config

import (
    "log"
    "os"
)

var (
    JwtSecret []byte
    DbURL     string
    Port      string
)

func Load() {
    DbURL = os.Getenv("DB_URL")
    if DbURL == "" {
        log.Fatal("DB_URL is required")
    }

    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        log.Fatal("JWT_SECRET is required")
    }
    JwtSecret = []byte(secret)

    Port = os.Getenv("PORT")
    if Port == "" {
        log.Fatal("PORT is required")
    }
}
