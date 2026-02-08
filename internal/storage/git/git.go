// Package git provides storage implemented via git
package git

import (
	"errors"
	"fmt"
	"path"
	"sync"
	"time"

	"protomorphine/tg-notes/internal/config"

	"github.com/go-git/go-git/v6"
	gitCfg "github.com/go-git/go-git/v6/config"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/object"
	"github.com/go-git/go-git/v6/plumbing/transport/ssh"
)

const commitMsg string = "note from tg-notes"

type GitStorage struct {
	*git.Repository
	*git.Worktree
	pubKey *ssh.PublicKeys
	config *config.GitRepository
	mu     sync.Mutex
}

func New(cfg *config.GitRepository) (*GitStorage, error) {
	const op = "storage.git.New"

	publicKeys, err := ssh.NewPublicKeys("git", []byte(cfg.Key), cfg.KeyPassword)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	repo, err := git.PlainOpenWithOptions(cfg.Path, &git.PlainOpenOptions{DetectDotGit: true})
	if errors.Is(err, git.ErrRepositoryNotExists) {
		repo, err = git.PlainClone(cfg.Path, &git.CloneOptions{
			Auth: publicKeys,
			URL:  cfg.URL,
		})
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// refs/heads/<cfg.Branch>
	localBranch := plumbing.NewBranchReferenceName(cfg.Branch)
	// refs/remotes/origin/<cfg.Branch>
	remoteBranch := plumbing.NewRemoteReferenceName("origin", cfg.Branch)

	// go-git do not see an existing branch by default
	// so we need to add it manually
	_ = repo.CreateBranch(&gitCfg.Branch{
		Name:   cfg.Branch,
		Remote: "origin",
		Merge:  localBranch,
	})

	// thx a lot for @ajschmidt8
	// https://github.com/go-git/go-git/issues/279#issuecomment-816359799
	newReference := plumbing.NewSymbolicReference(localBranch, remoteBranch)
	if err := repo.Storer.SetReference(newReference); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = repo.Branch(cfg.Branch)

	checkoutOpts := &git.CheckoutOptions{
		Branch: localBranch,
		Create: errors.Is(err, git.ErrBranchNotFound),
	}

	if err := worktree.Checkout(checkoutOpts); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &GitStorage{
		Repository: repo,
		Worktree:   worktree,
		pubKey:     publicKeys,
		config:     cfg,
	}, nil
}

func (g *GitStorage) Add(title, text string) error {
	const op = "storage.git.Add"

	g.mu.Lock()
	defer g.mu.Unlock()

	if err := g.prepareStorage(); err != nil {
		return fmt.Errorf("%s: error while preparing storage: %w", op, err)
	}

	path, err := g.createFile(title+".md", text)
	if err != nil {
		return fmt.Errorf("%s: error while saving file: %w", op, err)
	}

	if _, err := g.Worktree.Add(path); err != nil {
		return fmt.Errorf("%s: error while adding file: %w", op, err)
	}

	if err := g.save(); err != nil {
		return fmt.Errorf("%s: error while saving note: %w", op, err)
	}

	return nil
}

// createFile method    saves a new file with given filename and content. returns a path to saved file or error
func (g *GitStorage) createFile(filename, content string) (string, error) {
	path := path.Join(g.config.PathToSave, filename)

	file, err := g.Filesystem.Create(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = file.Write([]byte(content))
	return path, err
}

// save method    commit and push changes.
func (g *GitStorage) save() error {
	commitOpts := &git.CommitOptions{
		Committer: &object.Signature{
			Name: g.config.Committer.Name,
			When: time.Now(),
		},
	}

	if _, err := g.Commit(commitMsg, commitOpts); err != nil {
		return err
	}

	return g.Push(&git.PushOptions{
		Auth:       g.pubKey,
		RemoteName: "origin",
	})
}

func (g *GitStorage) prepareStorage() error {

	/*
	waiting for https://github.com/go-git/go-git/pull/1815

	pullOpts := &git.PullOptions{RemoteName: "origin", Auth: g.pubKey}
	if err := g.Pull(pullOpts); err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		return err
	}

	*/

	return nil
}
