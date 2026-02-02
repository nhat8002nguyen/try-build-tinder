package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		name  string
		email string
		want  bool
	}{
		{"valid email", "test@example.com", true},
		{"invalid - no @", "testexample.com", false},
		{"invalid - no domain", "test@", false},
		{"empty email", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidEmail(tt.email)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIsValidPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		want     bool
	}{
		{"valid password", "password123", true},
		{"too short", "pass", false},
		{"exactly 8 chars", "pass1234", true},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidPassword(tt.password)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCalculateAge(t *testing.T) {
	// Test with a known birthdate
	birthdate := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)
	age := CalculateAge(birthdate)
	
	// Age should be current year - 1990
	currentYear := time.Now().Year()
	expectedAge := currentYear - 1990
	
	assert.GreaterOrEqual(t, age, expectedAge-1)
	assert.LessOrEqual(t, age, expectedAge)
}

func TestSanitizeString(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"trim spaces", "  hello  ", "hello"},
		{"no spaces", "hello", "hello"},
		{"empty string", "", ""},
		{"only spaces", "   ", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SanitizeString(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}
