package utils

import (
	"crypto/rand"
	"math/big"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash)
}

// ComparePassword compares a hashed password with a plain text password
func ComparePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func GenerateSecurePassword() string {
	const (
		lowercase = "abcdefghijklmnopqrstuvwxyz"
		uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		numbers   = "0123456789"
		// special   = "!@#$%^&*()_+-=[]{}|;:,.<>?"
		pwdLength = 16
	)

	var password strings.Builder
	password.Grow(pwdLength)

	// Ensure at least one from each category
	categories := []string{lowercase, uppercase, numbers}
	for _, category := range categories {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(category))))
		password.WriteByte(category[n.Int64()])
	}

	// Fill remaining length with random chars from all categories
	allChars := lowercase + uppercase + numbers
	remainingLength := pwdLength - password.Len()
	for i := 0; i < remainingLength; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(allChars))))
		password.WriteByte(allChars[n.Int64()])
	}

	// Convert to string and shuffle
	result := []rune(password.String())
	for i := len(result) - 1; i > 0; i-- {
		j, _ := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		result[i], result[j.Int64()] = result[j.Int64()], result[i]
	}

	return string(result)
}
