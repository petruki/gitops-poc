package services

import (
	"testing"
)

func TestCheckoutAndCollectLastCommitDateWithGoGit(t *testing.T) {
	gs := &GoGitService{}

	repoURL := GetEnv("GIT_REPO_URL")
	token := GetEnv("GIT_TOKEN")

	lastCommitDate, err := gs.CheckoutAndCollectLastCommitDate(repoURL, token)
	AssertNil(t, err)
	AssertNotNil(t, lastCommitDate.String())
}

func TestCheckoutAndCollectLastCommitHashWithGoGit(t *testing.T) {
	gs := &GoGitService{}

	repoURL := GetEnv("GIT_REPO_URL")
	token := GetEnv("GIT_TOKEN")

	lastCommitHash, err := gs.CheckoutAndCollectLastCommitHash(repoURL, token)
	AssertNil(t, err)
	AssertNotNil(t, lastCommitHash)
}
