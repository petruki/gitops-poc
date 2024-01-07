package services

import (
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
)

type GoGitService struct{}

func (gs *GoGitService) GetLastCommitDate(repoURL string, token string) (time.Time, error) {
	c, _ := getCommitObject(repoURL, token)

	// Return the date of the commit
	return c.Author.When, nil
}

func (gs *GoGitService) GetLastCommitHash(repoURL string, token string) (string, error) {
	c, _ := getCommitObject(repoURL, token)

	// Return the date of the commit
	return c.Hash.String(), nil
}

func (gs *GoGitService) GetBranches(repoURL string, token string) ([]string, error) {
	r, _ := getRepository(repoURL, token)

	// Get the branchs
	branches, _ := r.Branches()

	var branchNames []string
	branches.ForEach(func(b *plumbing.Reference) error {
		branchNames = append(branchNames, b.Name().Short())
		return nil
	})

	return branchNames, nil
}

// Get file content from a repository
func (gs *GoGitService) GetFileContent(repoURL string, token string, filePath string) (string, error) {
	r, _ := getRepository(repoURL, token)

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

func getCommitObject(repoURL string, token string) (*object.Commit, error) {
	r, _ := getRepository(repoURL, token)

	// Get the HEAD reference
	ref, _ := r.Head()

	// Get the commit object from the reference
	return r.CommitObject(ref.Hash())
}

func getRepository(repoURL string, token string) (*git.Repository, error) {
	// Create a new repository object with the given URL
	return git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: repoURL,
		Auth: &http.BasicAuth{
			Username: "git-user",
			Password: token,
		},
	})
}
