package services

import (
	"testing"
)

func TestCheckoutAndCollectLastCommitDate(t *testing.T) {
	gs := &GoGitService{}

	repoURL := GetEnv("GIT_REPO_URL")
	token := GetEnv("GIT_TOKEN")

	lastCommitDate, err := gs.GetLastCommitDate(repoURL, token)
	AssertNil(t, err)
	AssertNotNil(t, lastCommitDate.String())
}

func TestCheckoutAndCollectLastCommitHash(t *testing.T) {
	gs := &GoGitService{}

	repoURL := GetEnv("GIT_REPO_URL")
	token := GetEnv("GIT_TOKEN")

	lastCommitHash, err := gs.GetLastCommitHash(repoURL, token)
	AssertNil(t, err)
	AssertNotNil(t, lastCommitHash)
}

func TestGetBranches(t *testing.T) {
	gs := &GoGitService{}

	repoURL := GetEnv("GIT_REPO_URL")
	token := GetEnv("GIT_TOKEN")

	branches, err := gs.GetBranches(repoURL, token)
	AssertNil(t, err)
	AssertNotNil(t, branches)
	AssertEqual(t, len(branches), 1)
	AssertEqual(t, branches[0], "master")
}

func TestGetFileContent(t *testing.T) {
	gs := &GoGitService{}

	repoURL := GetEnv("GIT_REPO_URL")
	token := GetEnv("GIT_TOKEN")

	filePath := "resources/default.json"

	content, err := gs.GetFileContent(repoURL, token, filePath)
	AssertNil(t, err)
	AssertNotNil(t, content)
	AssertContains(t, content, "Playground")
}
