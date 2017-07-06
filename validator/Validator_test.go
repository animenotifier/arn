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
