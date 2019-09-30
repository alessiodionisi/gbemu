package display

import (
	"errors"
	"fmt"
	"github.com/adnsio/gbemu/pkg/gameboy/bits"
	"image"
	"image/color"
)

const (
	Start = 0x8000
	End   = 0x9fff

	TileDataStart = 0x8000
	TileDataEnd   = 0x97ff
	TileDataSize  = TileDataEnd - TileDataStart + 1

	BackgroundMapStart = 0x9800
	BackgroundMapEnd   = 0x9bff
	BackgroundMapSize  = BackgroundMapEnd - BackgroundMapStart + 1

	WindowMapStart = 0x9c00
	WindowMapEnd   = 0x9fff
	WindowMapSize  = WindowMapEnd - WindowMapStart + 1

	OamStart = 0xfe00
	OamEnd   = 0xfe9f
	OamSize  = OamEnd - OamStart + 1

	Width  = 160
	Height = 144

	ControlBackgroundEnabled             = 0
	ControlSpriteEnabled                 = 1
	ControlSpriteSizeSelect              = 2
	ControlBackgroundMapSelect           = 3
	ControlBackgroundAndWindowTileSelect = 4
	ControlWindowEnabled                 = 5
	ControlWindowMapSelect               = 6
	ControlDisplayEnabled                = 7

	StatusModeFlag0            = 0
	StatusModeFlag1            = 1
	StatusCoincidenceFlag      = 2
	StatusMode0HBlankInterrupt = 3
	StatusMode1VBlankInterrupt = 4
	StatusMode2OamInterrupt    = 5
	StatusCoincidenceInterrupt = 6
)

var (
	// 0 White
	// 1 Light gray
	// 2 Dark gray
	// 3 Black

	GrayShades = [4]color.RGBA{
		color.RGBA{R: 232, G: 232, B: 232, A: 255},
		color.RGBA{R: 160, G: 160, B: 160, A: 255},
		color.RGBA{R: 88, G: 88, B: 88, A: 255},
		color.RGBA{R: 16, G: 16, B: 16, A: 255},
	}

	GreenLCDShades = [4]color.RGBA{
		color.RGBA{R: 224, G: 248, B: 208, A: 255},
		color.RGBA{R: 136, G: 192, B: 112, A: 255},
		color.RGBA{R: 52, G: 104, B: 86, A: 255},
		color.RGBA{R: 8, G: 24, B: 32, A: 255},
	}
)

type Display struct {
	Image             *image.RGBA
	ShadesOfGray      [4]color.RGBA
	Control           uint8
	Status            uint8
	ScrollY           uint8
	ScrollX           uint8
	CurrentLine       uint8
	CompareLine       uint8
	WindowY           uint8
	WindowX           uint8
	BackgroundPalette uint8
	ObjectPalette0    uint8
	ObjectPalette1    uint8
	DmaTransfer       uint8
	TileDataBank0     [TileDataSize]uint8
	TileDataBank1     [TileDataSize]uint8 // CGB
	BackgroundMap     [BackgroundMapSize]uint8
	WindowMap         [WindowMapSize]uint8
	Oam               [OamSize]uint8
}

func NewDisplay() *Display {
	return &Display{
		Image:        image.NewRGBA(image.Rect(0, 0, Width, Height)),
		ShadesOfGray: GrayShades,
	}
}

func (d *Display) Emulate() {

}

/*func (dpl *Display) CombineTileValues(val1 uint8, val2 uint8, bit uint8) uint8 {
	val := bits.Get(val2, bit)
	val <<= 1
	val |= bits.Get(val1, bit)

	return val
}

func (dpl *Display) ReadPaletteShade(pal uint8, color uint8) uint8 {
	hi := uint8(0)
	lo := uint8(0)

	switch color {
	case 0:
		hi = 1
		lo = 0
	case 1:
		hi = 3
		lo = 2
	case 2:
		hi = 5
		lo = 4
	case 3:
		hi = 7
		lo = 6
	}

	shade := uint8(0)
	shade = bits.Get(pal, hi) << 1
	shade |= bits.Get(pal, lo)

	return shade
}*/

func (d *Display) ReadBackgroundPalette(val uint8) uint8 {
	switch val {
	case 0:
		return d.BackgroundPalette >> 0 & 0x3
	case 1:
		return d.BackgroundPalette >> 2 & 0x3
	case 2:
		return d.BackgroundPalette >> 4 & 0x3
	case 3:
		return d.BackgroundPalette >> 6 & 0x3
	default:
		panic(errors.New(fmt.Sprintf("display: invalid palette %#02x", val)))
	}
}

func (d *Display) DrawLine() {
	if bits.Test(d.Control, ControlBackgroundEnabled) {
		bgMapSelect := bits.Get(d.Control, ControlBackgroundMapSelect)
		bgTileSelect := bits.Get(d.Control, ControlBackgroundAndWindowTileSelect)

		y := d.CurrentLine + d.ScrollY
		row := int(y / 8)

		for i := 0; i < 160; i++ {
			x := uint8(i) + d.ScrollX
			col := int(x / 8)

			var rawBgTileNum int
			// Bit 3 - BG Tile Map Display Select     (0=9800-9BFF, 1=9C00-9FFF)
			bgMapAddr := row*32 + col
			if bgMapSelect == 1 {
				rawBgTileNum = int(d.WindowMap[bgMapAddr])
			} else {
				rawBgTileNum = int(d.BackgroundMap[bgMapAddr])
			}

			var bgTileNum int
			// Bit 4 - BG & Window Tile Data Select   (0=8800-97FF, 1=8000-8FFF)
			if bgTileSelect == 1 {
				bgTileNum = rawBgTileNum
			} else {
				bgTileNum = 128 + (rawBgTileNum + 128)
			}

			bgTileAddr := uint16(bgTileNum * 16)

			line := (y % 8) * 2
			data1Addr := bgTileAddr + uint16(line)
			data2Addr := bgTileAddr + uint16(line+1)
			data1 := d.TileDataBank0[data1Addr]
			data2 := d.TileDataBank0[data2Addr]

			bit := (x%8 - 7) * 0xff
			colorVal := (bits.Get(data2, bit) << 1) | bits.Get(data1, bit)
			color := d.ReadBackgroundPalette(colorVal)
			d.Image.SetRGBA(int(i), int(d.CurrentLine), d.ShadesOfGray[color])
		}
	}

	if bits.Test(d.Control, ControlWindowEnabled) && d.WindowY <= d.CurrentLine {

	}
}

func (d *Display) Write(addr uint16, val uint8) {
	switch {
	case addr >= TileDataStart && addr <= TileDataEnd:
		// todo handle CGB
		d.TileDataBank0[addr-TileDataStart] = val
	case addr >= BackgroundMapStart && addr <= BackgroundMapEnd:
		d.BackgroundMap[addr-BackgroundMapStart] = val
	case addr >= WindowMapStart && addr <= WindowMapEnd:
		d.WindowMap[addr-WindowMapStart] = val
	case addr >= OamStart && addr <= OamEnd:
		// todo write only if OAM dma is enabled
		d.Oam[addr-OamStart] = val
	default:
		panic(errors.New(fmt.Sprintf("display: writing invalid address (%#04x)", addr)))
	}
}

func (d *Display) Read(addr uint16) uint8 {
	switch {
	case addr >= TileDataStart && addr <= TileDataEnd:
		// todo handle CGB
		return d.TileDataBank0[addr-TileDataStart]
	case addr >= BackgroundMapStart && addr <= BackgroundMapEnd:
		return d.BackgroundMap[addr-BackgroundMapStart]
	case addr >= WindowMapStart && addr <= WindowMapEnd:
		return d.WindowMap[addr-WindowMapStart]
	case addr >= OamStart && addr <= OamEnd:
		return d.Oam[addr-OamStart]
	default:
		panic(errors.New(fmt.Sprintf("display: reading invalid address (%#04x)", addr)))
	}
}
