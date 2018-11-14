package validate_test

import (
	"testing"

	"github.com/animenotifier/arn/validate"
	"github.com/stretchr/testify/assert"
)

func TestIsValidNick(t *testing.T) {
	// Invalid nicknames
	assert.False(t, validate.Nick(""))
	assert.False(t, validate.Nick("A"))
	assert.False(t, validate.Nick("AB CD"))
	assert.False(t, validate.Nick("A123"))
	assert.False(t, validate.Nick("A!§$%&/()=?`"))
	assert.False(t, validate.Nick("__"))
	assert.False(t, validate.Nick("Tsun.Dere"))
	assert.False(t, validate.Nick("Tsun Dere"))

	// Valid nicknames
	assert.True(t, validate.Nick("Tsundere"))
	assert.True(t, validate.Nick("TsunDere"))
	assert.True(t, validate.Nick("Tsun_Dere"))
	assert.True(t, validate.Nick("Akyoto"))
}

func TestIsValidEmail(t *testing.T) {
	assert.False(t, validate.Email(""))
	assert.True(t, validate.Email("support@notify.moe"))
}

func TestIsValidDate(t *testing.T) {
	assert.False(t, validate.DateTime(""))
	assert.False(t, validate.DateTime("0001-01-01T01:01:00Z"))
	assert.False(t, validate.DateTime("292277026596-12-04T15:30:07Z"))
	assert.True(t, validate.DateTime("2017-03-09T10:25:00Z"))
}

func TestIsValidURI(t *testing.T) {
	assert.False(t, validate.URI(""))
	assert.False(t, validate.URI("a"))
	assert.False(t, validate.URI("google.com"))
	assert.True(t, validate.URI("https://google.com"))
	assert.True(t, validate.URI("https://google.com/"))
	assert.True(t, validate.URI("https://google.com/images"))
	assert.True(t, validate.URI("https://google.com/images/"))
}
