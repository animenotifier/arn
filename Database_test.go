package arn

import "testing"
import "github.com/stretchr/testify/assert"

func TestConnect(t *testing.T) {
	assert.NotEmpty(t, DB.Node().Address().String())
}
