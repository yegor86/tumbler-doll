package shared

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
)

type GitRepo struct {
	Url            string
	Branch         string
	CloneDir       string
	Changelog      bool
	Auth    	   transport.AuthMethod
	Poll           bool
	ProgressWriter io.Writer
}

func (r *GitRepo) CloneOrPull() error {
	if _, err := os.Stat(r.CloneDir); os.IsNotExist(err) {
		// Clone the repository
		fmt.Println("Cloning repository...")
		return r.cloneRepo()
	}

	// Pull changes in the existing repository
	fmt.Println("Pulling changes...")
	return r.pullRepo()
}

func (r *GitRepo) cloneRepo() error {
	options := &git.CloneOptions{
		URL:           r.Url,
		Progress:      r.ProgressWriter,
		ReferenceName: plumbing.NewBranchReferenceName(r.Branch),
		Auth: r.Auth,
	}

	_, err := git.PlainClone(r.CloneDir, false, options)
	return err
}

func (r *GitRepo) pullRepo() error {

	// Open the existing repository
	repo, err := git.PlainOpen(r.CloneDir)
	if err != nil {
		return fmt.Errorf("failed to open repository: %w", err)
	}

	// Get the working tree
	worktree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	// Pull changes
	options := &git.PullOptions{
		RemoteName:    "origin",
		ReferenceName: plumbing.NewBranchReferenceName(r.Branch),
		Progress:      r.ProgressWriter,
	}

	err = worktree.Pull(options)
	if err == git.NoErrAlreadyUpToDate {
		fmt.Println("Repository is already up-to-date.")
		return nil
	}
	return err
}

// deriveCloneDir derives a clone directory name from a Git repository URL.
func DeriveCloneDir(repoURL string) (string, error) {
	// Parse the URL
	parsedURL, err := url.Parse(repoURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}

	repoPath := parsedURL.Path

	// Trim trailing ".git" if present
	repoName := strings.TrimSuffix(path.Base(repoPath), ".git")

	// Validate the repository name
	if repoName == "" {
		return "", fmt.Errorf("could not derive repository name from URL: %s", repoURL)
	}

	return repoName, nil
}
