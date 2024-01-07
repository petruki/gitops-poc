package services

import (
	"os"
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

func AssertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v, got %v", a, b)
	}
}
