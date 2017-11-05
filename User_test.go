package arn

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	user := NewUser()

	assert.NotNil(t, user)
	assert.NotEmpty(t, user.ID)
}

func TestDatabaseErrorMessages(t *testing.T) {
	_, err := GetUser("NON EXISTENT USER ID")

	// We need to make sure that non-existent records return "not found" in the error message.
	assert.NotNil(t, err)
	assert.NotEmpty(t, err.Error())
	assert.NotEqual(t, -1, strings.Index(err.Error(), "not found"))
}
