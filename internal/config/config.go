package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv loads environment variables
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using defaults")
	}
}

// GetImaginaryURL returns Imaginary API URL
func GetImaginaryURL() string {
	url := os.Getenv("IMAGINARY_URL")
	if url == "" {
		url = "http://localhost:9000" // Default URL
	}
	return url
}

// GetServerPort returns the server port
func GetServerPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}
