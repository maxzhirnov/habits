package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/maxzhirnov/habits/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockAppService struct {
	mock.Mock
}

func (m *MockAppService) AddNewHabit(habit models.Habit) error {
	args := m.Called(habit)
	return args.Error(0)
}

func TestAddNewHabitHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		mockAppService func() ApplicationService
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:        "Success",
			requestBody: `{"name": "Run", "user_id": "123"}`,
			mockAppService: func() ApplicationService {
				m := &MockAppService{}
				m.On("AddNewHabit", mock.Anything).Return(nil)
				return m
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"message":   "success",
				"habitName": "Run",
				"userID":    "123",
			},
		},
		// Add other test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(tt.requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			handler := &Handlers{ApplicationService: tt.mockAppService()}

			if err := handler.AddNewHabitHandler(c); err != nil && tt.expectedStatus != http.StatusBadRequest {
				t.Error(err)
			} else {
				assert.Equal(t, tt.expectedStatus, rec.Code)
				var responseBody map[string]interface{}
				json.Unmarshal(rec.Body.Bytes(), &responseBody)
				assert.Equal(t, tt.expectedBody, responseBody)
			}
		})
	}
}
