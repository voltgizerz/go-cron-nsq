package config

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

// LoadENV - load env file.
func LoadENV() {
	if err := godotenv.Load(); err != nil {
		log.Warn("No .env file found")
	}
}
