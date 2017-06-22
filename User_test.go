package arn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	user := NewUser()

	if !assert.NotNil(t, user) {
		return
	}

	assert.NotEmpty(t, user.ID)
}
