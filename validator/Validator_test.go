package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidNick(t *testing.T) {
	// Invalid nicknames
	assert.False(t, IsValidNick(""))
	assert.False(t, IsValidNick("A"))
	assert.False(t, IsValidNick("AB CD"))
	assert.False(t, IsValidNick("A123"))
	assert.False(t, IsValidNick("A!§$%&/()=?`"))
	assert.False(t, IsValidNick("__"))
	assert.False(t, IsValidNick("Tsun.Dere"))
	assert.False(t, IsValidNick("Tsun Dere"))

	// Valid nicknames
	assert.True(t, IsValidNick("Tsundere"))
	assert.True(t, IsValidNick("TsunDere"))
	assert.True(t, IsValidNick("Tsun_Dere"))
	assert.True(t, IsValidNick("Akyoto"))
}

func TestIsValidEmail(t *testing.T) {
	assert.False(t, IsValidEmail(""))
	assert.True(t, IsValidEmail("support@notify.moe"))
}

func TestIsValidDate(t *testing.T) {
	assert.False(t, IsValidDate(""))
	assert.False(t, IsValidDate("0001-01-01T01:01:00Z"))
	assert.False(t, IsValidDate("292277026596-12-04T15:30:07Z"))
	assert.True(t, IsValidDate("2017-03-09T10:25:00Z"))
}
