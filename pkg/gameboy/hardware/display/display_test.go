package display

import (
	"github.com/adnsio/gbemu/pkg/gameboy/bits"
	"testing"
)

/*func TestPPU_CombineTileValues(t *testing.T) {
	dpl := Display{}

	want := [8]uint8{
		2, // 10
		0, // 00
		3, // 11
		1, // 01
		2, // 10
		3, // 11
		2, // 10
		1, // 01
	}

	val1 := uint8(0x35) // 00110101
	val2 := uint8(0xae) // 10101110

	bitX := uint8(7)
	for x := 0; x < 8; x++ {
		val := dpl.CombineTileValues(val1, val2, bitX)

		if val != want[x] {
			t.Errorf("error: want %#02x, got %#02x", want[x], val)
		}

		bitX--
	}
}*/

func TestDisplay_DrawLine(t *testing.T) {
	d := NewDisplay()

	d.Control = bits.Set(d.Control, ControlBackgroundEnabled)
	d.Control = bits.Set(d.Control, ControlBackgroundAndWindowTileSelect)

	d.BackgroundPalette = 0xfc

	d.BackgroundMap[0x0000] = 1

	d.TileDataBank0[0x0010] = 0x35 // 00110101
	d.TileDataBank0[0x0011] = 0xae // 10101110

	//d.CurrentLine = 1

	d.DrawLine()
}
