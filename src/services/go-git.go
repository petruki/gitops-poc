package services

import (
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
)

type GoGitService struct{}

func (gs *GoGitService) CheckoutAndCollectLastCommitDate(repoURL string, token string) (time.Time, error) {
	r, err := getRepository(repoURL, token)

	// Get the HEAD reference
	ref, err := r.Head()
	if err != nil {
		return time.Time{}, err
	}

	// Get the commit object from the reference
	c, err := r.CommitObject(ref.Hash())
	if err != nil {
		return time.Time{}, err
	}

	// Return the date of the commit
	return c.Author.When, nil
}

func getRepository(repoURL string, token string) (*git.Repository, error) {
	// Create a new repository object with the given URL
	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: repoURL,
		Auth: &http.BasicAuth{
			Username: "git-user",
			Password: token,
		},
	})
	if err != nil {
		return nil, err
	}

	return r, nil
}
