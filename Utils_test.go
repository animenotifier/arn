package arn_test

import (
	"testing"

	"github.com/animenotifier/arn"
	"github.com/stretchr/testify/assert"
)

func TestContainsUnicodeLetters(t *testing.T) {
	assert.False(t, arn.ContainsUnicodeLetters("hello"))
	assert.True(t, arn.ContainsUnicodeLetters("こんにちは"))
	assert.True(t, arn.ContainsUnicodeLetters("hello こんにちは"))
}
