package services

import (
	"os"
	"os/exec"
	"strings"
)

type GitService struct {
	Directory string
}

func (gs *GitService) CheckoutAndCollectLastCommitDate(repoURL string, token string) (string, error) {
	// Insert the token into the repository URL
	repoURL = strings.Replace(repoURL, "https://", "https://"+token+"@", 1)

	// Clone the repository into the specified directory
	cmd := exec.Command("git", "clone", repoURL, gs.Directory)
	if err := cmd.Run(); err != nil {
		return "", err
	}

	// Get the date of the last commit
	cmd = exec.Command("git", "log", "-1", "--format=%cd")
	cmd.Dir = gs.Directory
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}

func (gs *GitService) RemoveTempDirectory() error {
	return os.RemoveAll(gs.Directory)
}
