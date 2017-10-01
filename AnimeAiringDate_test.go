package arn

import "testing"
import "github.com/stretchr/testify/assert"

func TestValidDate(t *testing.T) {
	date := AnimeAiringDate{}
	assert.False(t, date.IsValid())

	date.Start = invalidDate
	assert.False(t, date.IsValid())

	date.Start = "0001-01-08T00:00:00Z"
	assert.False(t, date.IsValid())

	date.Start = "2017-09-30T15:00:00Z"
	assert.True(t, date.IsValid())
}
