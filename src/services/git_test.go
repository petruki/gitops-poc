package services

import (
	"os"
	"testing"
)

func TestCheckoutAndCollectLastCommitDateWithCmd(t *testing.T) {
	gs := &GitService{GetDir() + "/../../TEMP"}

	repoURL := GetEnv("GIT_REPO_URL")
	token := GetEnv("GIT_TOKEN")

	lastCommitDate, err := gs.CheckoutAndCollectLastCommitDate(repoURL, token)
	AssertNil(t, err)
	AssertNotNil(t, lastCommitDate)
}

func TestMain(m *testing.M) {
	gs := &GitService{GetDir() + "/../../TEMP"}

	code := m.Run()
	gs.RemoveTempDirectory()
	os.Exit(code)
}
