package notesaving_test

import (
	"errors"
	"testing"

	"protomorphine/tg-notes/internal/app/usecases/notesaving"
	"protomorphine/tg-notes/internal/app/usecases/notesaving/mocks"
	"protomorphine/tg-notes/internal/config"
	"protomorphine/tg-notes/internal/domain"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	category     domain.Category             = "300 unknown"
	errAdderMock error                       = errors.New("failed to add")
	predictions  map[domain.Category]float64 = make(map[domain.Category]float64)
)

func TestSave(t *testing.T) {
	testCases := []struct {
		name            string
		text            string
		setupAdder      func(m *mocks.NoteAdder)
		setupClassifier func(m *mocks.Classifier)
		expectedErr     error
	}{
		{
			name: "success",
			text: "test note content",
			setupAdder: func(m *mocks.NoteAdder) {
				m.EXPECT().Add(mock.Anything, mock.AnythingOfType("domain.Note")).Return(nil).Once()
			},
			setupClassifier: func(m *mocks.Classifier) {
				m.EXPECT().Classify(mock.AnythingOfType("string")).Return(predictions, category)
			},
			expectedErr: nil,
		},
		{
			name: "adder returns error",
			text: "test note content",
			setupAdder: func(m *mocks.NoteAdder) {
				m.EXPECT().Add(mock.Anything, mock.AnythingOfType("domain.Note")).Return(errAdderMock).Once()
			},
			setupClassifier: func(m *mocks.Classifier) {
				m.EXPECT().Classify(mock.AnythingOfType("string")).Return(predictions, category)
			},
			expectedErr: errAdderMock,
		},
	}

	predictions[category] = .5

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockAdder := mocks.NewNoteAdder(t)
			tc.setupAdder(mockAdder)

			mockClassifier := mocks.NewClassifier(t)
			tc.setupClassifier(mockClassifier)

			uc := notesaving.New(mockAdder, mockClassifier, &config.NoteSaveConfig{CategoryThreshold: .1, DefaultCategory: "default"})
			_, err := uc.Save(t.Context(), tc.text)

			if tc.expectedErr != nil {
				require.ErrorIs(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
