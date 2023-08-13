package repos

import (
	"errors"
	"github.com/maxzhirnov/habits/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockDB struct {
	mock.Mock
}

func (mdb *MockDB) Insert(collectionName string, document interface{}) (string, error) {
	args := mdb.Called(collectionName, document)
	return args.String(0), args.Error(1)
}

func (mdb *MockDB) Exists(collectionName string, filters map[string]interface{}) (bool, error) {
	args := mdb.Called(collectionName, filters)
	return args.Bool(0), args.Error(1)
}

func TestCreateNewHabit(t *testing.T) {
	habit := models.Habit{
		Name:   "testHabit",
		UserID: "testUserID",
	}

	tests := []struct {
		name        string
		exists      bool
		existsErr   error
		insertErr   error
		expectedErr error
	}{
		{
			name:        "HabitExists",
			exists:      true,
			expectedErr: ErrHabitExists,
		},
		{
			name:   "SuccessfulCreation",
			exists: false,
		},
		{
			name:        "ExistsCheckError",
			exists:      false,
			existsErr:   errors.New("some error"),
			expectedErr: errors.New("some error"),
		},
		{
			name:        "InsertError",
			exists:      false,
			insertErr:   errors.New("insert error"),
			expectedErr: errors.New("insert error"),
		},
		{
			name:        "SuccessfulInsertionWithID",
			exists:      false,
			insertErr:   nil, // нет ошибки
			expectedErr: nil, // нет ошибки
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			repo := New(mockDB)

			mockDB.On("Exists", "habits", map[string]interface{}{
				"name":    habit.Name,
				"user_id": habit.UserID,
			}).Return(tt.exists, tt.existsErr)

			if !tt.exists && tt.existsErr == nil {
				mockDB.On("Insert", "habits", habit).Return("someID", tt.insertErr)
			}

			err := repo.CreateNewHabit(habit)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
