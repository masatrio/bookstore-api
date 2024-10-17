package middleware

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUserIDFromContext(t *testing.T) {
	tests := []struct {
		name       string
		ctx        context.Context
		expectedID int64
		expectedOk bool
	}{
		{
			name:       "UserID Present in Context",
			ctx:        context.WithValue(context.Background(), userIDKey, int64(12345)),
			expectedID: 12345,
			expectedOk: true,
		},
		{
			name:       "UserID Not Present in Context",
			ctx:        context.Background(),
			expectedID: 0,
			expectedOk: false,
		},
		{
			name:       "UserID Present as Different Type",
			ctx:        context.WithValue(context.Background(), userIDKey, "not an int64"),
			expectedID: 0,
			expectedOk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, ok := GetUserIDFromContext(tt.ctx)
			assert.Equal(t, tt.expectedID, userID)
			assert.Equal(t, tt.expectedOk, ok)
		})
	}
}
