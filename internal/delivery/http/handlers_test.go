package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/masatrio/bookstore-api/internal/domain/usecase"
	"github.com/masatrio/bookstore-api/internal/domain/usecase/mocks"
	"github.com/masatrio/bookstore-api/utils"
	"github.com/stretchr/testify/assert"
)

func TestRegisterHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mocks.NewMockUserUseCase(ctrl)
	handler := &Handler{userUseCase: mockUserUseCase}

	tests := []struct {
		name           string
		input          usecase.RegisterInput
		expectedStatus int
		mockError      error
	}{
		{
			name: "Success",
			input: usecase.RegisterInput{
				Name:     "satrio",
				Email:    "satrio@example.com",
				Password: "securepassword",
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Error from use case",
			input: usecase.RegisterInput{
				Name:     "satrio",
				Email:    "satrio@example.com",
				Password: "securepassword",
			},
			expectedStatus: http.StatusInternalServerError,
			mockError:      utils.NewCustomSystemError("System Error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputJSON, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(inputJSON))
			w := httptest.NewRecorder()

			if tt.mockError != nil {
				mockUserUseCase.EXPECT().
					Register(gomock.Any(), gomock.Eq(tt.input)).
					Return(nil, tt.mockError)
			} else {
				mockUserUseCase.EXPECT().
					Register(gomock.Any(), gomock.Eq(tt.input)).
					Return(&usecase.RegisterOutput{}, nil)
			}

			handler.RegisterHandler(w, req)

			res := w.Result()
			if res.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d; got %d", tt.expectedStatus, res.StatusCode)
			}
		})
	}
}

func TestLoginHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mocks.NewMockUserUseCase(ctrl)
	handler := NewHandler(mockUserUseCase, nil, nil)

	tests := []struct {
		name           string
		input          usecase.LoginInput
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Success",
			input: usecase.LoginInput{
				Email:    "satrio@example.com",
				Password: "securepassword",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Missing Email",
			input: usecase.LoginInput{
				Email:    "satrio@example.com",
				Password: "securepassword",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Email and password are required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(body))
			w := httptest.NewRecorder()

			if tt.expectedStatus == http.StatusOK {
				mockUserUseCase.EXPECT().Login(gomock.Any(), tt.input).Return(&usecase.LoginOutput{}, nil)
			} else {
				mockUserUseCase.EXPECT().Login(gomock.Any(), tt.input).Return(nil, utils.NewCustomUserError(tt.expectedError))
			}

			handler.LoginHandler(w, req)

			res := w.Result()
			assert.Equal(t, tt.expectedStatus, res.StatusCode)

			if tt.expectedError != "" {
				var errResponse map[string]string
				json.NewDecoder(w.Body).Decode(&errResponse)
				assert.Equal(t, tt.expectedError, errResponse["error"])
			}
		})
	}
}

func TestListBooksHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookUseCase := mocks.NewMockBookUseCase(ctrl)
	handler := &Handler{bookUseCase: mockBookUseCase}

	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		mockResponse   *usecase.ListBooksOutput
		mockError      error
	}{
		{
			name:           "Success",
			queryParams:    "?title=Go&author=Author1",
			expectedStatus: http.StatusOK,
			mockResponse:   &usecase.ListBooksOutput{},
		},
		{
			name:           "Error from use case",
			queryParams:    "?title=Go&author=Author1",
			expectedStatus: http.StatusInternalServerError,
			mockError:      utils.NewCustomSystemError("System Error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/books"+tt.queryParams, nil)
			w := httptest.NewRecorder()

			if tt.mockError != nil {
				mockBookUseCase.EXPECT().ListBooks(gomock.Any(), gomock.Any()).Return(nil, tt.mockError)
			} else {
				mockBookUseCase.EXPECT().ListBooks(gomock.Any(), gomock.Any()).Return(tt.mockResponse, nil)
			}

			handler.ListBooksHandler(w, req)

			res := w.Result()
			if res.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d; got %d", tt.expectedStatus, res.StatusCode)
			}
		})
	}
}

func TestHealthCheckHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	handler := &Handler{}
	handler.HealthCheckHandler(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status %d; got %d", http.StatusOK, res.StatusCode)
	}

	var response map[string]string
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Errorf("failed to decode response: %v", err)
	}

	if response["status"] != "healthy" {
		t.Errorf("expected 'healthy'; got %s", response["status"])
	}
}
