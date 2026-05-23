// Package notesaving provides usecase for saving notes
package notesaving

import (
	"context"
	"fmt"
	"time"

	"protomorphine/tg-notes/internal/app/models"
	"protomorphine/tg-notes/internal/config"
	"protomorphine/tg-notes/internal/domain"
)

// NoteCreator is an interface for creating notes..
//
//mockery:generate: true
type NoteCreator interface {
	Create(ctx context.Context, text string) (models.SaveResult, error)
}

// NoteSaver is an interface for saving the note.
//
//mockery:generate: true
type NoteSaver interface {
	Save(ctx context.Context, note domain.Note) error
}

// Classifier is an interface for text classification.
//
//mockery:generate: true
type Classifier interface {
	Classify(content string) (map[domain.Category]float64, domain.Category)
}

// Usecase represents the usecase for saving notes.
type Usecase struct {
	classifier Classifier
	adder      NoteSaver
	cfg        *config.NoteSaveConfig
}

// New creates a new Usecase.
func New(adder NoteSaver, classifier Classifier, cfg *config.NoteSaveConfig) *Usecase {
	return &Usecase{
		cfg:        cfg,
		adder:      adder,
		classifier: classifier,
	}
}

// Create creates a new note.
func (u *Usecase) Create(ctx context.Context, text string) (models.SaveResult, error) {
	const op = "app.usecase.notesaving.Save"

	probs, category := u.classifier.Classify(text)

	if probs[category] < u.cfg.CategoryThreshold {
		category = domain.Category(u.cfg.DefaultCategory)
	}

	title := fmt.Sprintf("note (%v)", time.Now().Format(time.DateTime))

	note := domain.Note{
		Title:    title,
		Content:  text,
		Category: category,
	}

	if err := u.adder.Save(ctx, note); err != nil {
		return models.SaveResult{}, fmt.Errorf("%s: error while saving note: %w", op, err)
	}

	return models.SaveResult{Title: note.Title, Category: note.Category}, nil
}
