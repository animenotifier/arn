package arn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveSpecialCharacters(t *testing.T) {
	assert.Equal(t, RemoveSpecialCharacters("Hello World"), "Hello World")
	assert.Equal(t, RemoveSpecialCharacters("Aldnoah.Zero 2"), "Aldnoah Zero 2")
	assert.Equal(t, RemoveSpecialCharacters("Working!"), "Working ")
	assert.Equal(t, RemoveSpecialCharacters("Working!!"), "Working  ")
	assert.Equal(t, RemoveSpecialCharacters("Working!!!"), "Working   ")
	assert.Equal(t, RemoveSpecialCharacters("Lucky☆Star"), "Lucky Star")
	assert.Equal(t, RemoveSpecialCharacters("ChäoS;Child"), "ChäoS Child")
	assert.Equal(t, RemoveSpecialCharacters("僕だけがいない街"), "僕だけがいない街")
}
