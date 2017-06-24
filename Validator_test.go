package arn

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

func TestFixUserNick(t *testing.T) {
	// Nickname autocorrect
	assert.True(t, FixUserNick("Akyoto") == "Akyoto")
	assert.True(t, FixUserNick("Tsundere") == "Tsundere")
	assert.True(t, FixUserNick("akyoto") == "Akyoto")
	assert.True(t, FixUserNick("aky123oto") == "Akyoto")
	assert.True(t, FixUserNick("__aky123oto%$§") == "Akyoto")
	assert.True(t, FixUserNick("__aky123oto%$§__") == "Akyoto")
	assert.True(t, FixUserNick("123%&/(__%") == "")
}
