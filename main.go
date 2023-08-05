package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	letters    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers    = "0123456789"
	codeLength = 12
)

var generatedCodes = make(map[string]bool)

func generateRandomCode() string {
	rand.Seed(time.Now().UnixNano())

	code := make([]byte, codeLength)
	for i := 0; i < codeLength; i++ {
		if i < 6 {
			code[i] = letters[rand.Intn(len(letters))]
		} else {
			code[i] = numbers[rand.Intn(len(numbers))]
		}
	}

	return string(code)
}

func codeExists(code string) bool {
	_, exists := generatedCodes[code]
	return exists
}

func main() {
	// Generate and check codes
	for i := 0; i < 10; i++ {
		randomCode := generateRandomCode()
		if codeExists(randomCode) {
			fmt.Println("Generated Code already exists:", randomCode)
		} else {
			generatedCodes[randomCode] = true
			fmt.Println("Generated Code:", randomCode)
		}
	}
}
