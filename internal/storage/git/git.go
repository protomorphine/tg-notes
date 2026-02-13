// Package git provides storage implemented via git
package git

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"path"
	"sync"
	"time"

	"protomorphine/tg-notes/internal/config"
	"protomorphine/tg-notes/internal/log"

	"github.com/go-git/go-git/v6"
	gitCfg "github.com/go-git/go-git/v6/config"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/object"
	"github.com/go-git/go-git/v6/plumbing/transport/ssh"
)

const commitMsg string = "note from tg-notes"

var noErrNothingToDo error = errors.New("nothing to do")

type GitStorage struct {
	worktree *git.Worktree
	repo     *git.Repository
	pubKey   *ssh.PublicKeys

	config *config.GitRepository

	mu sync.Mutex

	buf       []string
	bufFullCh chan struct{}
}

// New creates instance of GitStorage
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
		config:    cfg,
		repo:      repo,
		worktree:  worktree,
		pubKey:    publicKeys,
		buf:       make([]string, 0, cfg.BufSize),
		bufFullCh: make(chan struct{}, 1),
	}, nil
}

// Add adds new note to storage
func (g *GitStorage) Add(ctx context.Context, title, text string) error {
	const op = "storage.git.Add"

	g.mu.Lock()
	defer g.mu.Unlock()

	path, err := g.createFile(title+".md", text)
	if err != nil {
		return fmt.Errorf("%s: error while saving file: %w", op, err)
	}

	g.buf = append(g.buf, path)

	if len(g.buf) == g.config.BufSize {
		select {
		case g.bufFullCh <- struct{}{}:
		default:
		}
	}

	return nil
}

// Processor starts update remote loop
func (g *GitStorage) Processor(ctx context.Context, logger *slog.Logger) {
	const op = "storage.git.Processor"
	logger = logger.With(log.Op(op))

	duration := g.config.UpdateDuratiion

	logger.Info("starting update storage", slog.String("duration", duration.String()))
	ticker := time.Tick(duration)

	for {
		select {
		case <-ctx.Done():
			close(g.bufFullCh)
			return

		case _, ok := <-g.bufFullCh:
			if !ok {
				continue
			}

			logger.Debug("start updating remote storage, reason: buffer is full")
			g.triggerUpdate(ctx, logger)

		case _, ok := <-ticker:
			if !ok {
				continue
			}

			logger.Debug("start updating remote storage, reason: timer")
			g.triggerUpdate(ctx, logger)
		}
	}
}

func (g *GitStorage) triggerUpdate(ctx context.Context, logger *slog.Logger) {
	saved, err := g.handlePendingNotes(ctx)
	if err != nil {

		if errors.Is(err, noErrNothingToDo) {
			logger.Debug("no new notes to save")
			return
		}

		logger.Error("error while handling pending notes", log.Err(err))
		return
	}

	logger.Info("notes saved successfully", slog.Int("count", saved))
}

func (g *GitStorage) handlePendingNotes(ctx context.Context) (int, error) {
	const op = "storage.git.handlePendingNotes"

	if err := ctx.Err(); err != nil {
		return 0, fmt.Errorf("%s: context err: %w", op, err)
	}

	g.mu.Lock()

	if len(g.buf) == 0 {
		g.mu.Unlock()
		return 0, noErrNothingToDo
	}

	buf := make([]string, len(g.buf))
	copy(buf, g.buf)

	g.buf = g.buf[:0]
	g.mu.Unlock()

	if err := g.prepareStorage(ctx); err != nil {
		return 0, fmt.Errorf("%s: error while preparing storage: %w", op, err)
	}

	for _, path := range buf {
		if _, err := g.worktree.Add(path); err != nil {
			return 0, fmt.Errorf("%s: error while add file %s to worktree: %w", op, path, err)
		}
	}

	if err := g.save(ctx); err != nil {
		return 0, fmt.Errorf("%s: error while saving notes: %w", op, err)
	}

	return len(buf), nil
}

func (g *GitStorage) createFile(filename, content string) (string, error) {
	path := path.Join(g.config.PathToSave, filename)

	file, err := g.worktree.Filesystem.Create(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = file.Write([]byte(content))
	return path, err
}

func (g *GitStorage) save(ctx context.Context) error {
	commitOpts := &git.CommitOptions{
		Committer: &object.Signature{
			Name: g.config.Committer.Name,
			When: time.Now(),
		},
	}

	if _, err := g.worktree.Commit(commitMsg, commitOpts); err != nil {
		return err
	}

	return g.repo.PushContext(ctx, &git.PushOptions{
		Auth:       g.pubKey,
		RemoteName: "origin",
	})
}

func (g *GitStorage) prepareStorage(ctx context.Context) error {
	pullOpts := &git.PullOptions{RemoteName: "origin", Auth: g.pubKey, Force: true}
	if err := g.worktree.PullContext(ctx, pullOpts); err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		return err
	}

	return nil
}
