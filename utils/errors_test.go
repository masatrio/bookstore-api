package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCustomUserError(t *testing.T) {
	message := "This is a user error"
	err := NewCustomUserError(message)

	assert.NotNil(t, err)
	assert.Equal(t, UserError, err.Type)
	assert.Equal(t, message, err.Message)
	assert.True(t, err.IsUserError())
	assert.False(t, err.IsSystemError())
}

func TestNewCustomSystemError(t *testing.T) {
	message := "This is a system error"
	err := NewCustomSystemError(message)

	assert.NotNil(t, err)
	assert.Equal(t, SystemError, err.Type)
	assert.Equal(t, message, err.Message)
	assert.False(t, err.IsUserError())
	assert.True(t, err.IsSystemError())
}

func TestCustomError_ErrorMethod(t *testing.T) {
	message := "This is a custom error"
	err := NewCustomUserError(message)

	assert.Equal(t, message, err.Error())
}
