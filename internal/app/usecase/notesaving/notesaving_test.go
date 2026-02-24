package notesaving_test

import (
	"errors"
	"testing"

	"protomorphine/tg-notes/internal/app/usecase/notesaving"
	"protomorphine/tg-notes/internal/app/usecase/notesaving/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var adderErr error = errors.New("failed to add")

func TestSave(t *testing.T) {
	testCases := []struct {
		name        string
		text        string
		setupAdder  func(m *mocks.NoteAdder)
		expectedErr error
	}{
		{
			name: "success",
			text: "test note content",
			setupAdder: func(m *mocks.NoteAdder) {
				m.EXPECT().Add(mock.Anything, mock.AnythingOfType("string"), "test note content").Return(nil).Once()
			},
			expectedErr: nil,
		},
		{
			name: "adder returns error",
			text: "test note content",
			setupAdder: func(m *mocks.NoteAdder) {
				m.EXPECT().Add(mock.Anything, mock.AnythingOfType("string"), "test note content").Return(adderErr).Once()
			},
			expectedErr: adderErr,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockAdder := mocks.NewNoteAdder(t)
			tc.setupAdder(mockAdder)

			uc := notesaving.New(mockAdder)
			err := uc.Save(t.Context(), tc.text)

			if tc.expectedErr != nil {
				require.ErrorIs(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}

			mockAdder.AssertExpectations(t)
		})
	}
}
