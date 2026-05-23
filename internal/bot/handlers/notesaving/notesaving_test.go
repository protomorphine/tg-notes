package notesaving_test

import (
	"errors"
	"log/slog"
	"testing"

	appmodels "protomorphine/tg-notes/internal/app/models"
	ucmocks "protomorphine/tg-notes/internal/app/usecases/notesaving/mocks"
	"protomorphine/tg-notes/internal/bot/handlers/notesaving"
	"protomorphine/tg-notes/internal/bot/handlers/notesaving/mocks"
	"protomorphine/tg-notes/internal/log"

	"github.com/go-telegram/bot/models"
	"github.com/stretchr/testify/mock"
)

func TestNilMessage(t *testing.T) {
	update := &models.Update{Message: nil}

	creator := ucmocks.NewNoteCreator(t)
	sender := mocks.NewMessageSender(t)

	logger := slog.New(log.NewDiscardHandler())
	h := notesaving.New(logger, creator)

	h(t.Context(), sender, update)

	creator.AssertNotCalled(t, "Save")
	sender.AssertNotCalled(t, "SendMessage")
}

func TestEmptyMessageText(t *testing.T) {
	update := &models.Update{
		Message: &models.Message{
			Text:    "",
			Caption: "valid caption",
		},
	}

	creator := ucmocks.NewNoteCreator(t)
	sender := mocks.NewMessageSender(t)

	creator.EXPECT().Create(mock.Anything, mock.AnythingOfType("string")).Return(appmodels.SaveResult{}, nil)
	sender.EXPECT().SendMessage(mock.Anything, mock.Anything).Return(nil, nil).Maybe()

	logger := slog.New(log.NewDiscardHandler())
	h := notesaving.New(logger, creator)

	h(t.Context(), sender, update)
}

func TestTextAndCaptionEmpty(t *testing.T) {
	update := &models.Update{
		Message: &models.Message{
			Text:    "",
			Caption: "",
		},
	}

	creator := ucmocks.NewNoteCreator(t)
	sender := mocks.NewMessageSender(t)

	sender.EXPECT().SendMessage(mock.Anything, mock.Anything).Return(nil, nil).Maybe()

	logger := slog.New(log.NewDiscardHandler())
	h := notesaving.New(logger, creator)

	h(t.Context(), sender, update)

	sender.AssertExpectations(t)
	creator.AssertNotCalled(t, "Save")
}

func TestAddNote(t *testing.T) {
	tests := []struct {
		name         string
		update       *models.Update
		setupCreator func(*ucmocks.NoteCreator)
	}{
		{
			name: "message text is not empty",
			update: &models.Update{
				Message: &models.Message{
					Text: "some text",
				},
			},
			setupCreator: func(adder *ucmocks.NoteCreator) {
				adder.EXPECT().Create(mock.Anything, mock.AnythingOfType("string")).Return(appmodels.SaveResult{}, nil)
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
			setupCreator: func(adder *ucmocks.NoteCreator) {
				adder.EXPECT().Create(mock.Anything, mock.AnythingOfType("string")).Return(appmodels.SaveResult{}, nil)
			},
		},
		{
			name: "Add returns err",
			update: &models.Update{
				Message: &models.Message{
					Text: "some text",
				},
			},
			setupCreator: func(adder *ucmocks.NoteCreator) {
				adder.EXPECT().Create(mock.Anything, mock.AnythingOfType("string")).Return(appmodels.SaveResult{}, errors.New("internal adder error"))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			creator := ucmocks.NewNoteCreator(t)
			tc.setupCreator(creator)

			sender := mocks.NewMessageSender(t)

			sender.EXPECT().SendMessage(mock.Anything, mock.Anything).Return(nil, nil).Maybe()

			logger := slog.New(log.NewDiscardHandler())
			h := notesaving.New(logger, creator)

			h(t.Context(), sender, tc.update)
		})
	}
}
