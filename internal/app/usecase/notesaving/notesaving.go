// Package notesaving provides usecase for saving notes
package notesaving

import (
	"context"
	"fmt"
	"time"
)

// NoteAdder is an interface for adding a note.
//
//mockery:generate: true
type NoteAdder interface {
	Add(ctx context.Context, title, text string) error
}

// Usecase represents the usecase for saving notes.
type Usecase struct {
	adder NoteAdder
}

// New creates a new Usecase.
func New(adder NoteAdder) *Usecase {
	return &Usecase{
		adder: adder,
	}
}

// Save saves a new note.
func (u *Usecase) Save(ctx context.Context, text string) error {
	const op = "app.usecase.notesaving.Save"

	title := fmt.Sprintf("note (%v)", time.Now().Format(time.DateTime))

	if err := u.adder.Add(ctx, title, text); err != nil {
		return fmt.Errorf("%s: error while saving note: %w", op, err)
	}

	return nil
}
