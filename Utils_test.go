package arn

import "testing"
import "github.com/stretchr/testify/assert"

func TestContainsUnicodeLetters(t *testing.T) {
	assert.False(t, ContainsUnicodeLetters("hello"))
	assert.True(t, ContainsUnicodeLetters("こんにちは"))
	assert.True(t, ContainsUnicodeLetters("hello こんにちは"))
}
