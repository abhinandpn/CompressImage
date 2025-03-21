package repository

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// ReadFile reads file data into memory
func ReadFile(file io.Reader) ([]byte, error) {
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// SaveImageToStorage saves the image to the storage folder
func SaveImageToStorage(filename string, data []byte) (string, error) {
	path := fmt.Sprintf("storage/%s", filename)
	err := os.WriteFile(path, data, 0644)
	if err != nil {
		return "", err
	}
	return path, nil
}

func S3upload(file multipart.File, fileName string, fileType string) (string, error) {
	// Initialize AWS session
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(os.Getenv("AWS_BUCKET_REGION")),
		Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"), ""),
	})
	if err != nil {
		log.Println("failed to initialize AWS session: ", err)
		return "", fmt.Errorf("failed to initialize AWS session: %v", err)
	}

	svc := s3.New(sess)
	bucket := os.Getenv("AWS_BUCKET_NAME")
	key := "imaginary/" + fileName

	// Upload the file directly using PutObject
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        file,                 // Use the multipart.File directly
		ContentType: aws.String(fileType), // Ensure correct MIME type (e.g., image/jpeg)
	})
	if err != nil {
		log.Printf("failed to upload file: %v", err)
		return "", fmt.Errorf("failed to upload file: %v", err)
	}

	// Return the public URL
	return "https://" + bucket + ".s3.amazonaws.com/" + key, nil
}
