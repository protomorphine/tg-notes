package handlers_test

import (
	"errors"
	"log/slog"
	"testing"

	"protomorphine/tg-notes/internal/bot/handlers"
	"protomorphine/tg-notes/internal/bot/handlers/mocks"
	"protomorphine/tg-notes/internal/log"

	"github.com/go-telegram/bot/models"
	"github.com/stretchr/testify/mock"
)

func TestNilMessage(t *testing.T) {
	update := &models.Update{Message: nil}

	adder := mocks.NewNoteAdder(t)
	sender := mocks.NewMessageSender(t)

	logger := slog.New(log.NewDiscardHandler())
	h := handlers.NewNoteSaving(logger, adder)

	h(t.Context(), sender, update)

	if !adder.AssertNotCalled(t, "Add", mock.Anything, mock.Anything) {
		t.Error("adder.Add called when message is nil")
	}

	if !sender.AssertNotCalled(t, "SendMessage", mock.Anything, mock.Anything) {
		t.Error("sender.SendMessage called when message is nil")
	}
}

func TestEmptyMessageText(t *testing.T) {
	update := &models.Update{
		Message: &models.Message{
			Text:    "",
			Caption: "valid caption",
		},
	}

	adder := mocks.NewNoteAdder(t)
	sender := mocks.NewMessageSender(t)

	adder.EXPECT().Add(t.Context(), mock.Anything, "valid caption").Return(nil)
	sender.EXPECT().SendMessage(t.Context(), mock.Anything).Return(nil, nil)

	logger := slog.New(log.NewDiscardHandler())
	h := handlers.NewNoteSaving(logger, adder)

	h(t.Context(), sender, update)
}

func TestTextAndCaptionEmpty(t *testing.T) {
	update := &models.Update{
		Message: &models.Message{
			Text:    "",
			Caption: "",
		},
	}

	adder := mocks.NewNoteAdder(t)
	sender := mocks.NewMessageSender(t)

	sender.EXPECT().SendMessage(t.Context(), mock.Anything).Return(nil, nil)

	logger := slog.New(log.NewDiscardHandler())
	h := handlers.NewNoteSaving(logger, adder)

	h(t.Context(), sender, update)
}

func TestAddNote(t *testing.T) {
	tests := []struct {
		name       string
		update     *models.Update
		adderSetup func(*mocks.NoteAdder)
	}{
		{
			name: "message text is not empty",
			update: &models.Update{
				Message: &models.Message{
					Text: "some text",
				},
			},

			adderSetup: func(adder *mocks.NoteAdder) {
				adder.EXPECT().Add(mock.Anything, mock.Anything, "some text").Return(nil)
			},
		},
		{
			name: "message text is empty, caption is empty",
			update: &models.Update{
				Message: &models.Message{
					Text:    "",
					Caption: "some caption",
				},
			},

			adderSetup: func(adder *mocks.NoteAdder) {
				adder.EXPECT().Add(mock.Anything, mock.Anything, "some caption").Return(nil)
			},
		},
		{
			name: "text is not empty, Add returns err",
			update: &models.Update{
				Message: &models.Message{
					Text:    "",
					Caption: "some caption",
				},
			},

			adderSetup: func(adder *mocks.NoteAdder) {
				adder.EXPECT().Add(mock.Anything, mock.Anything, "some caption").Return(errors.New("internal adder error"))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			adder := mocks.NewNoteAdder(t)
			tc.adderSetup(adder)

			sender := mocks.NewMessageSender(t)

			sender.EXPECT().SendMessage(t.Context(), mock.Anything).Return(nil, nil)

			logger := slog.New(log.NewDiscardHandler())
			h := handlers.NewNoteSaving(logger, adder)

			h(t.Context(), sender, tc.update)
		})
	}
}
