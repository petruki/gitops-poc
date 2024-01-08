package services

import (
	"os"
	"strings"
	"testing"
)

// Helpers

func GetDir() string {
	directory, err := os.Getwd()
	if err != nil {
		return ""
	}

	return directory
}

func AssertNotNil(t *testing.T, object interface{}) {
	if object == nil {
		t.Errorf("Object is nil")
	}
}

func AssertNil(t *testing.T, object interface{}) {
	if object != nil {
		t.Errorf("Object is not nil")
	}
}

func AssertEqual(t *testing.T, actual interface{}, expected interface{}) {
	if actual != expected {
		t.Errorf("Expected %v, got %v", actual, expected)
	}
}

func AssertContains(t *testing.T, actual string, expected string) {
	if !strings.Contains(actual, expected) {
		t.Errorf("Expected %v to contain %v", actual, expected)
	}
}
