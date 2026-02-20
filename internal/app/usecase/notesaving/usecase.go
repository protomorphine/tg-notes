package notesaving

import (
	"context"
	"fmt"
	"time"
)

//mockery:generate: true
type NoteAdder interface {
	Add(ctx context.Context, title, text string) error
}

type Usecase struct {
	adder NoteAdder
}

func New(adder NoteAdder) *Usecase {
	return &Usecase{
		adder: adder,
	}
}

func (u *Usecase) Save(ctx context.Context, text string) error {
	const op = "app.usecase.notesaving.Save"

	title := fmt.Sprintf("note (%v)", time.Now().Format(time.DateTime))

	if err := u.adder.Add(ctx, title, text); err != nil {
		return fmt.Errorf("%s: error while saving note: %w", op, err)
	}

	return nil
}
