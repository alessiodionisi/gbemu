package hardware

import (
	"errors"
	"fmt"
	"github.com/adnsio/gbemu/pkg/gameboy/hardware/audio"
	"github.com/adnsio/gbemu/pkg/gameboy/hardware/bootrom"
	"github.com/adnsio/gbemu/pkg/gameboy/hardware/cartridge"
	"github.com/adnsio/gbemu/pkg/gameboy/hardware/display"
	"github.com/adnsio/gbemu/pkg/gameboy/hardware/timer"
)

const (
	DATA_START      = 0x0000
	DATA_END        = 0xffff
	DATA_SIZE       = DATA_END - DATA_START + 1
	BOOTROM_START   = 0x0000
	BOOTROM_END     = 0x00ff
	BOOTROM_SIZE    = BOOTROM_END - BOOTROM_START + 1
	CARTRIDGE_START = 0x0000
	CARTRIDGE_END   = 0x7fff
	//CARTRIDGE_SIZE  = CARTRIDGE_END - CARTRIDGE_START + 1

	LCDC_BACKGROUND_ENABLED = 0
	LCDC_SPRITES_ENABLED    = 0
	LCDC_DISPLAY_ENABLED    = 7

	STAT_COINCIDENCE_FLAG       = 2
	STAT_MODE0_HBLANK_INTERRUPT = 3
	STAT_MODE1_VBLANK_INTERRUPT = 4
	STAT_MODE2_OAM_INTERRUPT    = 5
	STAT_COINCIDENCE_INTERRUPT  = 6

	IoStart = 0xff00
	IoEnd   = 0xff7f

	NotUsableStart = 0xfea0
	NotUsableEnd   = 0xfeff

	EchoStart = 0xe000
	EchoEnd   = 0xfdff

	HighRamStart = 0xff80
	HighRamEnd   = 0xfffe
	HighRamSize  = HighRamEnd - HighRamStart + 1

	WorkRamBank0Start = 0xc000
	WorkRamBank0End   = 0xcfff
	WorkRamBank0Size  = WorkRamBank0End - WorkRamBank0Start + 1

	WorkRamBankNStart = 0xd000
	WorkRamBankNEnd   = 0xdfff
	WorkRamBankNSize  = WorkRamBankNEnd - WorkRamBankNStart + 1
)

const (
	IO_SB              = 0xff01
	IO_SC              = 0xff02
	IO_DIV             = 0xff04
	IO_TIMA            = 0xff05
	IO_TMA             = 0xff06
	IO_TAC             = 0xff07
	IO_IF              = 0xff0f
	IO_NR11            = 0xff11
	IO_NR12            = 0xff12
	IO_NR13            = 0xff13
	IO_NR14            = 0xff14
	IO_NR50            = 0xff24
	IO_NR51            = 0xff25
	IO_NR52            = 0xff26
	IO_LCDC            = 0xff40
	IO_STAT            = 0xff41
	IO_SCY             = 0xff42
	IO_SCX             = 0xff43
	IO_LY              = 0xff44
	IO_LYC             = 0xff45
	IO_BGP             = 0xff47
	IO_OBP0            = 0xff48
	IO_OBP1            = 0xff49
	IO_DISABLE_BOOTROM = 0xff50
	IO_IE              = 0xffff
)

type Hardware struct {
	Bootrom   *bootrom.Bootrom
	Cartrdige *cartridge.Cartridge
	Display   *display.Display
	//Irq           *irq.Irq
	Timer        *timer.Timer
	Audio        *audio.Audio
	HighRam      [HighRamSize]uint8
	WorkRamBank0 [WorkRamBank0Size]uint8
	WorkRamBankN [WorkRamBankNSize]uint8 // CGB
	//EmulationTime int
}

func NewHardware() *Hardware {
	return &Hardware{
		Bootrom:   bootrom.NewBootrom(),
		Cartrdige: cartridge.NewCartridge(),
		Display:   display.NewDisplay(),
		//Irq:       irq.NewIrq(),
		Timer: timer.NewTimer(),
		Audio: audio.NewAudio(),
	}
}

func (h *Hardware) Read(addr uint16) uint8 {
	switch {
	case addr >= bootrom.Start && addr <= bootrom.End:
		if h.Bootrom.Enabled {
			return h.Bootrom.Read(addr)
		} else {
			return h.Cartrdige.Read(addr)
		}
	case addr >= cartridge.Start && addr <= cartridge.End:
		return h.Cartrdige.Read(addr)
	case addr >= display.Start && addr <= display.End:
		return h.Display.Read(addr)
	case addr >= cartridge.RamStart && addr <= cartridge.RamEnd:
		return h.Cartrdige.Read(addr)
	case addr >= WorkRamBank0Start && addr <= WorkRamBank0End:
		return h.WorkRamBank0[addr-WorkRamBank0Start]
	case addr >= WorkRamBankNStart && addr <= WorkRamBankNEnd:
		return h.WorkRamBankN[addr-WorkRamBankNStart]
	case addr >= EchoStart && addr <= EchoEnd:
		echoWRamAddr := addr - EchoStart + WorkRamBank0Start
		return h.Read(echoWRamAddr)
	case addr >= display.OamStart && addr <= display.OamEnd:
		return h.Display.Read(addr)
	case addr >= NotUsableStart && addr <= NotUsableEnd:
		fmt.Printf("memory: reading not usable (%#04x)\n", addr)
		return 0
	case addr >= IoStart && addr <= IoEnd:
		ioAddr := addr & 0xff

		switch ioAddr {
		case 0x00:
			// todo joypad register
			fmt.Printf("memory: reading joypad register io (%#04x)\n", addr)
			return 0
		case 0x01:
			// todo serial data
			fmt.Printf("memory: reading serial data io (%#04x)\n", addr)
			return 0
		case 0x02:
			// todo serial control
			fmt.Printf("memory: reading serial control io (%#04x)\n", addr)
			return 0
		case 0x04:
			return h.Timer.DividerRegister
		case 0x05:
			return h.Timer.Counter
		case 0x06:
			return h.Timer.Modulo
		case 0x07:
			return h.Timer.Control
		case 0x40:
			return h.Display.Control
		case 0x41:
			return h.Display.Status
		case 0x42:
			return h.Display.ScrollY
		case 0x43:
			return h.Display.ScrollX
		case 0x44:
			return h.Display.CurrentLine
		case 0x45:
			return h.Display.CompareLine
		case 0x46:
			return h.Display.DmaTransfer
		case 0x47:
			return h.Display.BackgroundPalette
		case 0x48:
			return h.Display.ObjectPalette0
		case 0x49:
			return h.Display.ObjectPalette1
		case 0x4a:
			return h.Display.WindowY
		case 0x4b:
			return h.Display.WindowX
		default:
			//panic(errors.New(fmt.Sprintf("memory: reading invalid io (%#04x)", addr)))
			fmt.Printf("memory: reading invalid io (%#04x)\n", addr)
			return 0
		}
	case addr >= HighRamStart && addr <= HighRamEnd:
		return h.HighRam[addr-HighRamStart]
	default:
		panic(errors.New(fmt.Sprintf("memory: reading invalid address (%#04x)", addr)))
	}
}

func (h *Hardware) Read16(addr uint16) uint16 {
	return uint16(h.Read(addr)) | uint16(h.Read(addr+1))<<8
}

func (h *Hardware) Write(addr uint16, val uint8) {
	switch {
	case addr >= cartridge.Start && addr <= cartridge.End:
		fmt.Printf("memory: writing cartridge (%#04x) %#02x\n", addr, val)
	case addr >= display.Start && addr <= display.End:
		h.Display.Write(addr, val)
	case addr >= WorkRamBank0Start && addr <= WorkRamBank0End:
		h.WorkRamBank0[addr-WorkRamBank0Start] = val
	case addr >= WorkRamBankNStart && addr <= WorkRamBankNEnd:
		h.WorkRamBankN[addr-WorkRamBankNStart] = val
	case addr >= EchoStart && addr <= EchoEnd:
		echoWRamAddr := addr - EchoStart + WorkRamBank0Start
		h.Write(echoWRamAddr, val)
	case addr >= display.OamStart && addr <= display.OamEnd:
		h.Display.Write(addr, val)
	case addr >= NotUsableStart && addr <= NotUsableEnd:
		fmt.Printf("memory: writing not usable (%#04x) %#02x\n", addr, val)
	case addr >= IoStart && addr <= IoEnd:
		ioAddr := addr & 0xff

		switch ioAddr {
		case 0x00:
			// todo joypad register
			fmt.Printf("memory: writing joypad register io (%#04x)\n", addr)
		case 0x01:
			// todo serial data
			fmt.Printf("memory: writing serial data io (%#04x)\n", addr)
		case 0x02:
			// todo serial control
			fmt.Printf("memory: writing serial control io (%#04x)\n", addr)
		case 0x04:
			h.Timer.DividerRegister = 0
		case 0x05:
			h.Timer.Counter = val
		case 0x06:
			h.Timer.Modulo = val
		case 0x07:
			h.Timer.Control = val
		case 0x0f:
			// todo interrupt flag
			fmt.Printf("memory: writing interrupt flag io (%#04x)\n", addr)
		//case ioAddr >= 0x10 && ioAddr <= 0x3f:
		// todo sound io
		//	fmt.Printf("memory: writing sound io (%#04x)\n", addr)
		case 0x40:
			h.Display.Control = val
		case 0x41:
			h.Display.Status = val
		case 0x42:
			h.Display.ScrollY = val
		case 0x43:
			h.Display.ScrollX = val
		case 0x44:
			h.Display.CurrentLine = 0
		case 0x45:
			h.Display.CompareLine = val
		case 0x46:
			h.Display.DmaTransfer = val
		case 0x47:
			h.Display.BackgroundPalette = val
		case 0x48:
			h.Display.ObjectPalette0 = val
		case 0x49:
			h.Display.ObjectPalette1 = val
		case 0x4a:
			h.Display.WindowY = val
		case 0x4b:
			h.Display.WindowX = val
		case 0x50:
			h.Bootrom.Enabled = false
		default:
			//panic(errors.New(fmt.Sprintf("memory: writing invalid io (%#04x)", addr)))
			fmt.Printf("memory: writing invalid io (%#04x)\n", addr)
		}
	case addr >= HighRamStart && addr <= HighRamEnd:
		h.HighRam[addr-HighRamStart] = val
	case addr == 0xffff:
		// todo interrupt enable
		fmt.Printf("memory: writing interrupt enable (%#04x)\n", addr)
	default:
		panic(errors.New(fmt.Sprintf("memory: writing invalid address (%#04x)", addr)))
	}
}

func (h *Hardware) Write16(addr uint16, val uint16) {
	h.Write(addr, uint8(val&0xff))
	h.Write(addr+1, uint8(val>>8))
}
