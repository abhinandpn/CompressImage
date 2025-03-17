package service

import (
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/abhinandpn/CompressImage/server"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// ProcessAndCompressImage handles image processing via Imaginary API (concurrent)
func ProcessAndCompressImage(filename string, imageData []byte, size int64, originalWidth, originalHeight int) (map[string]string, error) {
	baseName := strings.TrimSuffix(filename, filepath.Ext(filename))
	baseName = strings.ReplaceAll(baseName, " ", "_")

	if cachedPaths, exists := GetCachedResult(baseName); exists {
		return cachedPaths, nil
	}

	var wg sync.WaitGroup
	resultChan := make(chan struct {
		key  string
		path string
		err  error
	}, 4)

	// Define compression quality levels and target sizes
	sizes := map[string]int{
		"original":  determineOriginalSizeReduction(size), // Original size with potential reduction
		"250-300KB": 80,                                   // Adjust quality for 250-300KB
		"150-200KB": 60,                                   // Adjust quality for 150-200KB
		"10-50KB":   20,                                   // Adjust quality for 10-50KB
	}

	aspectRatio := float64(originalWidth) / float64(originalHeight)

	for key, quality := range sizes {
		wg.Add(1)
		go func(k string, q int) {
			defer wg.Done()

			// Calculate dimensions based on aspect ratio
			var newWidth, newHeight int
			if k == "original" {
				newWidth, newHeight = originalWidth, originalHeight // Keep original dimensions
			} else {
				// Scale dimensions to maintain aspect ratio (start with a reasonable base width)
				baseWidth := 1200 // Starting point, adjust if needed
				newWidth = baseWidth
				newHeight = int(float64(newWidth) / aspectRatio)
			}

			// Process the image with consistent dimensions
			path, err := server.ProcessImageWithImaginary(imageData, q, baseName+"_"+k, newWidth, newHeight)
			resultChan <- struct {
				key  string
				path string
				err  error
			}{key: k, path: path, err: err}
		}(key, quality)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	imagePaths := make(map[string]string)
	for res := range resultChan {
		if res.err == nil {
			imagePaths[res.key] = res.path
		}
	}

	CacheResult(baseName, imagePaths)
	return imagePaths, nil
}

// S3ProcessAndCompressImage handles image processing and uploads to S3
func S3ProcessAndCompressImage(filename string, imageData []byte, size int64, originalWidth, originalHeight int) (map[string]string, error) {
	baseName := strings.TrimSuffix(filename, filepath.Ext(filename))
	baseName = strings.ReplaceAll(baseName, " ", "_")

	// Check if the image is cached
	if cachedPaths, exists := GetCachedResult(baseName); exists {
		return cachedPaths, nil
	}

	var wg sync.WaitGroup
	resultChan := make(chan struct {
		key  string
		path string
		err  error
	}, 4)

	// Define compression quality levels and target sizes
	sizes := map[string]int{
		"original":  determineOriginalSizeReduction(size), // Original size with potential reduction
		"250-300KB": 80,                                   // Adjust quality for 250-300KB
		"150-200KB": 60,                                   // Adjust quality for 150-200KB
		"10-50KB":   20,                                   // Adjust quality for 10-50KB
	}

	// Calculate aspect ratio from original dimensions
	aspectRatio := float64(originalWidth) / float64(originalHeight)

	// Process the image in different sizes concurrently
	for key, quality := range sizes {
		wg.Add(1)
		go func(k string, q int) {
			defer wg.Done()

			// Calculate dimensions based on aspect ratio
			var newWidth, newHeight int
			if k == "original" {
				newWidth, newHeight = originalWidth, originalHeight // Keep original dimensions
			} else {
				// Scale dimensions to maintain aspect ratio (start with a reasonable base width)
				baseWidth := 1200 // Starting point, adjust if needed
				newWidth = baseWidth
				newHeight = int(float64(newWidth) / aspectRatio)
			}

			// Process the image with consistent dimensions
			path, err := server.ProcessImageWithImaginary(imageData, q, baseName+"_"+k, newWidth, newHeight)
			if err != nil {
				resultChan <- struct {
					key  string
					path string
					err  error
				}{key: k, path: "", err: fmt.Errorf("failed to process image: %v", err)}
				return
			}

			// Open the processed image file
			file, err := os.Open(path)
			if err != nil {
				resultChan <- struct {
					key  string
					path string
					err  error
				}{key: k, path: "", err: fmt.Errorf("failed to open file: %v", err)}
				return
			}
			defer file.Close()

			// Upload the image to S3
			s3URL, uploadErr := S3Imageupload(file, baseName+"_"+k+".jpg", "image/jpeg")
			if uploadErr == nil {
				resultChan <- struct {
					key  string
					path string
					err  error
				}{key: k, path: s3URL, err: nil}
			} else {
				resultChan <- struct {
					key  string
					path string
					err  error
				}{key: k, path: "", err: uploadErr}
			}
		}(key, quality)
	}

	// Wait for all goroutines to finish
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect the results from each goroutine
	imagePaths := make(map[string]string)
	for res := range resultChan {
		if res.err == nil {
			imagePaths[res.key] = res.path
		}
	}

	// Cache the results
	CacheResult(baseName, imagePaths)
	return imagePaths, nil
}

// Determines the original image compression based on file size
func determineOriginalSizeReduction(size int64) int {
	switch {
	case size > 5*1024*1024: // If size > 5MB
		return 50
	case size > 2*1024*1024: // If size is between 2MB and 5MB
		return 70
	default: // If size is < 2MB
		return 100 // No compression
	}
}

// S3Imageupload uploads a file to the AWS S3 bucket
func S3Imageupload(file multipart.File, fileName string, fileType string) (string, error) {
	// Initialize AWS session
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(os.Getenv("AWS_BUCKET_REGION")),
		Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"), ""),
	})
	if err != nil {
		log.Println("Failed to initialize AWS session:", err)
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
		log.Printf("Failed to upload file: %v", err)
		return "", fmt.Errorf("failed to upload file: %v", err)
	}

	// Return the public URL
	return "https://" + bucket + ".s3.amazonaws.com/" + key, nil
}
