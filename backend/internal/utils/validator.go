package utils

import (
	"regexp"
	"strings"
	"time"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func IsValidEmail(email string) bool {
	return emailRegex.MatchString(strings.ToLower(strings.TrimSpace(email)))
}

func IsValidPassword(password string) bool {
	return len(password) >= 8
}

func IsValidAge(birthdate *time.Time, minAge, maxAge int) bool {
	if birthdate == nil {
		return false
	}

	now := time.Now()
	age := now.Year() - birthdate.Year()
	if now.YearDay() < birthdate.YearDay() {
		age--
	}

	return age >= minAge && age <= maxAge
}

func CalculateAge(birthdate time.Time) int {
	now := time.Now()
	age := now.Year() - birthdate.Year()
	if now.YearDay() < birthdate.YearDay() {
		age--
	}
	return age
}

func SanitizeString(s string) string {
	return strings.TrimSpace(s)
}
