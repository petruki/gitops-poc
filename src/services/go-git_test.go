package services

import (
	"testing"
)

func TestCheckoutAndCollectLastCommitDate(t *testing.T) {
	gs := &GoGitService{
		repoURL: GetEnv("GIT_REPO_URL"),
		token:   GetEnv("GIT_TOKEN"),
	}

	lastCommitDate, err := gs.GetLastCommitDate()
	AssertNil(t, err)
	AssertNotNil(t, lastCommitDate.String())
}

func TestCheckoutAndCollectLastCommitHash(t *testing.T) {
	gs := &GoGitService{
		repoURL: GetEnv("GIT_REPO_URL"),
		token:   GetEnv("GIT_TOKEN"),
	}

	lastCommitHash, err := gs.GetLastCommitHash()
	println(lastCommitHash)
	AssertNil(t, err)
	AssertNotNil(t, lastCommitHash)
}

func TestGetFileContent(t *testing.T) {
	gs := &GoGitService{
		repoURL: GetEnv("GIT_REPO_URL"),
		token:   GetEnv("GIT_TOKEN"),
	}

	filePath := "resources/default.json"

	content, err := gs.GetFileContent(filePath)
	AssertNil(t, err)
	AssertNotNil(t, content)
	AssertContains(t, content, "Release 1")
}
