package bits_test

import (
	"testing"

	"github.com/adnsio/gbemu/pkg/bits"
	"gotest.tools/assert"
)

func TestGetAndSetLower(t *testing.T) {
	data := uint16(0x0000)
	data = bits.SetLower(data, 0x01)

	assert.Equal(t, uint8(0x01), bits.GetLower(data))
}

func TestGetAndSetHigher(t *testing.T) {
	data := uint16(0x0000)
	data = bits.SetHigher(data, 0x01)

	assert.Equal(t, uint8(0x01), bits.GetHigher(data))
}
