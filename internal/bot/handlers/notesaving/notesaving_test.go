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

	saver := ucmocks.NewNoteSaver(t)
	sender := mocks.NewMessageSender(t)

	logger := slog.New(log.NewDiscardHandler())
	h := notesaving.New(logger, saver)

	h(t.Context(), sender, update)

	saver.AssertNotCalled(t, "Save")
	sender.AssertNotCalled(t, "SendMessage")
}

func TestEmptyMessageText(t *testing.T) {
	update := &models.Update{
		Message: &models.Message{
			Text:    "",
			Caption: "valid caption",
		},
	}

	saver := ucmocks.NewNoteSaver(t)
	sender := mocks.NewMessageSender(t)

	saver.EXPECT().Save(mock.Anything, mock.AnythingOfType("string")).Return(appmodels.SaveResult{}, nil)
	sender.EXPECT().SendMessage(mock.Anything, mock.Anything).Return(nil, nil).Maybe()

	logger := slog.New(log.NewDiscardHandler())
	h := notesaving.New(logger, saver)

	h(t.Context(), sender, update)
}

func TestTextAndCaptionEmpty(t *testing.T) {
	update := &models.Update{
		Message: &models.Message{
			Text:    "",
			Caption: "",
		},
	}

	saver := ucmocks.NewNoteSaver(t)
	sender := mocks.NewMessageSender(t)

	sender.EXPECT().SendMessage(mock.Anything, mock.Anything).Return(nil, nil).Maybe()

	logger := slog.New(log.NewDiscardHandler())
	h := notesaving.New(logger, saver)

	h(t.Context(), sender, update)

	sender.AssertExpectations(t)
	saver.AssertNotCalled(t, "Save")
}

func TestAddNote(t *testing.T) {
	tests := []struct {
		name       string
		update     *models.Update
		setupSaver func(*ucmocks.NoteSaver)
	}{
		{
			name: "message text is not empty",
			update: &models.Update{
				Message: &models.Message{
					Text: "some text",
				},
			},
			setupSaver: func(adder *ucmocks.NoteSaver) {
				adder.EXPECT().Save(mock.Anything, mock.AnythingOfType("string")).Return(appmodels.SaveResult{}, nil)
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
			setupSaver: func(adder *ucmocks.NoteSaver) {
				adder.EXPECT().Save(mock.Anything, mock.AnythingOfType("string")).Return(appmodels.SaveResult{}, nil)
			},
		},
		{
			name: "Add returns err",
			update: &models.Update{
				Message: &models.Message{
					Text: "some text",
				},
			},
			setupSaver: func(adder *ucmocks.NoteSaver) {
				adder.EXPECT().Save(mock.Anything, mock.AnythingOfType("string")).Return(appmodels.SaveResult{}, errors.New("internal adder error"))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			saver := ucmocks.NewNoteSaver(t)
			tc.setupSaver(saver)

			sender := mocks.NewMessageSender(t)

			sender.EXPECT().SendMessage(mock.Anything, mock.Anything).Return(nil, nil).Maybe()

			logger := slog.New(log.NewDiscardHandler())
			h := notesaving.New(logger, saver)

			h(t.Context(), sender, tc.update)
		})
	}
}
