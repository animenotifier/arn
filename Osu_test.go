package arn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOsuUser(t *testing.T) {
	userName := "Aky"
	user, err := GetOsuUser(userName)

	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, user.UserName, userName)
	assert.NotEmpty(t, user.PPRaw)
	assert.NotEmpty(t, user.Level)
	assert.NotEmpty(t, user.Accuracy)
	assert.NotEmpty(t, user.PlayCount)
}
