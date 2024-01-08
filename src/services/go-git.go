package services

import (
	"time"

	"github.com/go-git/go-billy/v5/memfs"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
)

type GoGitService struct {
	repoURL    string
	token      string
	branchName string
}

func (gs *GoGitService) GetLastCommitDate() (time.Time, error) {
	c, _ := gs.getCommitObject()

	// Return the date of the commit
	return c.Author.When, nil
}

func (gs *GoGitService) GetLastCommitHash() (string, error) {
	c, _ := gs.getCommitObject()

	// Return the date of the commit
	return c.Hash.String(), nil
}

func (gs *GoGitService) GetFileContent(filePath string) (string, error) {
	r, _ := gs.getRepository()

	// Get the HEAD reference
	ref, _ := r.Head()

	// Get the commit object from the reference
	c, _ := r.CommitObject(ref.Hash())

	// Get the tree from the commit object
	tree, _ := c.Tree()

	// Get the file
	f, _ := tree.File(filePath)

	// Get the content
	content, _ := f.Contents()

	return content, nil
}

func (gs *GoGitService) getCommitObject() (*object.Commit, error) {
	r, _ := gs.getRepository()

	// Get the HEAD reference
	ref, _ := r.Head()

	// Get the commit object from the reference
	return r.CommitObject(ref.Hash())
}

func (gs *GoGitService) getRepository() (*git.Repository, error) {
	// Clone repository using in-memory storage
	r, _ := git.Clone(memory.NewStorage(), memfs.New(), &git.CloneOptions{
		URL: gs.repoURL,
		Auth: &http.BasicAuth{
			Username: "git-user",
			Password: gs.token,
		},
	})

	// Checkout branch
	return gs.checkoutBranch(*r)
}

func (gs *GoGitService) checkoutBranch(r git.Repository) (*git.Repository, error) {
	// Fetch worktree
	w, err := r.Worktree()

	if err != nil {
		return nil, err
	}

	// Checkout remote branch
	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewRemoteReferenceName("origin", gs.branchName),
	})

	return &r, err
}
