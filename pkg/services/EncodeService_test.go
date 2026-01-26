package services

import (
	"io"
	"testing"

	"github.com/rs/zerolog"
)

func TestEncodeSha256(t *testing.T) {
	logger := zerolog.New(io.Discard)
	service := NewEncodeService(&logger)

	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "empty string",
			content:  "",
			expected: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
		{
			name:     "hello world",
			content:  "hello world",
			expected: "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.EncodeSha256(tt.content)
			if result != tt.expected {
				t.Errorf("EncodeSha256() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestComputeSha256(t *testing.T) {
	logger := zerolog.New(io.Discard)
	service := NewEncodeService(&logger)

	tests := []struct {
		name          string
		input         string
		comparisonSha string
		expected      bool
	}{
		{
			name:          "matching sha",
			input:         "hello world",
			comparisonSha: "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9",
			expected:      true,
		},
		{
			name:          "non-matching sha",
			input:         "hello world",
			comparisonSha: "wrong-sha",
			expected:      false,
		},
		{
			name:          "empty input matching sha",
			input:         "",
			comparisonSha: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			expected:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.ComputeSha256(tt.input, tt.comparisonSha)
			if result != tt.expected {
				t.Errorf("ComputeSha256() = %v, want %v", result, tt.expected)
			}
		})
	}
}
