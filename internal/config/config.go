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

// GetAWSRegion returns the AWS region
func GetAWSRegion() string {
	region := os.Getenv("AWS_BUCKET_REGION")
	if region == "" {
		region = "ap-south-1" // Default to "ap-south-1" if not set
	}
	return region
}

// GetAWSAccessKey returns the AWS access key
func GetAWSAccessKey() string {
	return os.Getenv("AWS_ACCESS_KEY")
}

// GetAWSSecretKey returns the AWS secret key
func GetAWSSecretKey() string {
	return os.Getenv("AWS_SECRET_KEY")
}

// GetAWSBucketName returns the AWS S3 bucket name
func GetAWSBucketName() string {
	return os.Getenv("AWS_BUCKET_NAME")
}
