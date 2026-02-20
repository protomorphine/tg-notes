package notesaving_test

import (
	"errors"
	"log/slog"
	"testing"

	notesavingUC "protomorphine/tg-notes/internal/app/usecase/notesaving"
	notesavingUCmocks "protomorphine/tg-notes/internal/app/usecase/notesaving/mocks"
	"protomorphine/tg-notes/internal/bot/handlers/notesaving"
	"protomorphine/tg-notes/internal/bot/handlers/notesaving/mocks"
	"protomorphine/tg-notes/internal/log"

	"github.com/go-telegram/bot/models"
	"github.com/stretchr/testify/mock"
)

func TestNilMessage(t *testing.T) {
	update := &models.Update{Message: nil}

	adder := notesavingUCmocks.NewNoteAdder(t)
	usecase := notesavingUC.New(adder)
	sender := mocks.NewMessageSender(t)

	logger := slog.New(log.NewDiscardHandler())
	h := notesaving.New(logger, usecase)

	h(t.Context(), sender, update)

	adder.AssertNotCalled(t, "Add", mock.Anything, mock.Anything, mock.Anything)
	sender.AssertNotCalled(t, "SendMessage", mock.Anything, mock.Anything)
}

func TestEmptyMessageText(t *testing.T) {
	update := &models.Update{
		Message: &models.Message{
			Text:    "",
			Caption: "valid caption",
		},
	}

	adder := notesavingUCmocks.NewNoteAdder(t)
	usecase := notesavingUC.New(adder)
	sender := mocks.NewMessageSender(t)

	adder.EXPECT().Add(mock.Anything, mock.AnythingOfType("string"), "valid caption").Return(nil)
	sender.EXPECT().SendMessage(mock.Anything, mock.Anything).Return(nil, nil)

	logger := slog.New(log.NewDiscardHandler())
	h := notesaving.New(logger, usecase)

	h(t.Context(), sender, update)

	adder.AssertExpectations(t)
	sender.AssertExpectations(t)
}

func TestTextAndCaptionEmpty(t *testing.T) {
	update := &models.Update{
		Message: &models.Message{
			Text:    "",
			Caption: "",
		},
	}

	adder := notesavingUCmocks.NewNoteAdder(t)
	usecase := notesavingUC.New(adder)
	sender := mocks.NewMessageSender(t)

	sender.EXPECT().SendMessage(mock.Anything, mock.Anything).Return(nil, nil)

	logger := slog.New(log.NewDiscardHandler())
	h := notesaving.New(logger, usecase)

	h(t.Context(), sender, update)

	sender.AssertExpectations(t)
	adder.AssertNotCalled(t, "Add", mock.Anything, mock.Anything, mock.Anything)
}

func TestAddNote(t *testing.T) {
	tests := []struct {
		name       string
		update     *models.Update
		adderSetup func(*notesavingUCmocks.NoteAdder)
	}{
		{
			name: "message text is not empty",
			update: &models.Update{
				Message: &models.Message{
					Text: "some text",
				},
			},
			adderSetup: func(adder *notesavingUCmocks.NoteAdder) {
				adder.EXPECT().Add(mock.Anything, mock.AnythingOfType("string"), "some text").Return(nil)
			},
		},
		{
			name: "message text is empty, caption is not empty",
			update: &models.Update{
				Message: &models.Message{
					Text:    "",
					Caption: "some caption",
				},
			},
			adderSetup: func(adder *notesavingUCmocks.NoteAdder) {
				adder.EXPECT().Add(mock.Anything, mock.AnythingOfType("string"), "some caption").Return(nil)
			},
		},
		{
			name: "Add returns err",
			update: &models.Update{
				Message: &models.Message{
					Text: "some text",
				},
			},
			adderSetup: func(adder *notesavingUCmocks.NoteAdder) {
				adder.EXPECT().Add(mock.Anything, mock.AnythingOfType("string"), "some text").Return(errors.New("internal adder error"))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			adder := notesavingUCmocks.NewNoteAdder(t)
			tc.adderSetup(adder)

			usecase := notesavingUC.New(adder)
			sender := mocks.NewMessageSender(t)

			sender.EXPECT().SendMessage(mock.Anything, mock.Anything).Return(nil, nil)

			logger := slog.New(log.NewDiscardHandler())
			h := notesaving.New(logger, usecase)

			h(t.Context(), sender, tc.update)

			adder.AssertExpectations(t)
			sender.AssertExpectations(t)
		})
	}
}
