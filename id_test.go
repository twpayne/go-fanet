package fanet_test

import (
	"testing"

	"github.com/alecthomas/assert/v2"

	"github.com/twpayne/go-fanet"
)

func TestID(t *testing.T) {
	id := fanet.ID{
		Manufacturer: 0x12,
		Device:       0x3456,
	}

	actualText, err := id.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, []byte("12:3456"), actualText)
	assert.Equal(t, 0x123456, id.Int())
	assert.False(t, id.IsZero())
	assert.Equal(t, "12:3456", id.String())
	assert.Equal(t, "id", id.Type())

	assert.EqualError(t, id.Set("AB:VXYZ"), "invalid ID")
	assert.Equal(t, 0x123456, id.Int())

	assert.NoError(t, id.Set("AB:CDEF"))
	assert.Equal(t, 0xABCDEF, id.Int())
}
