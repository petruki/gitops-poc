package services

import (
	"os"
	"testing"
)

func TestCheckoutAndCollectLastCommitDate(t *testing.T) {
	gs := &GitService{GetDir() + "/../../TEMP"}

	// Use a public repository for testing
	repoURL := GetEnv("GIT_REPO_URL")
	token := GetEnv("GIT_TOKEN")

	_, err := gs.CheckoutAndCollectLastCommitDate(repoURL, token)
	if err != nil {
		t.Errorf("CheckoutAndCollectLastCommitDate failed with error: %v", err)
	}
}

func TestMain(m *testing.M) {
	gs := &GitService{GetDir() + "/../../TEMP"}

	code := m.Run()
	gs.RemoveTempDirectory()
	os.Exit(code)
}

// Helpers

func GetDir() string {
	directory, err := os.Getwd()
	if err != nil {
		return ""
	}

	return directory
}
