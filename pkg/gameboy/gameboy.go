package gameboy

import (
	"fmt"
	"github.com/adnsio/gbemu/pkg/gameboy/bits"
	"github.com/adnsio/gbemu/pkg/gameboy/cpu"
	"github.com/adnsio/gbemu/pkg/gameboy/hardware"
	"github.com/adnsio/gbemu/pkg/gameboy/hardware/display"
)

type Config struct {
	Bootrom   []uint8
	Cartridge []uint8
}

type GameBoy struct {
	ClockSpeed    int
	CPU           *cpu.CPU
	Hardware      *hardware.Hardware
	DisplayCycles int
	Paused        bool
	ForcedPause   bool
}

func NewGameBoy(cfg Config) *GameBoy {
	hwe := hardware.NewHardware()
	cpu := cpu.NewCPU(hwe)

	gb := &GameBoy{
		ClockSpeed: 4194304,
		Hardware:   hwe,
		CPU:        cpu,
	}

	if cfg.Bootrom != nil {
		hwe.Bootrom.Load(cfg.Bootrom)
		hwe.Bootrom.Enabled = true
	} else {
		cpu.WriteAF(0x01b0)
		cpu.WriteBC(0x0013)
		cpu.WriteDE(0x00d8)
		cpu.WriteHL(0x014d)
		cpu.SP = 0xfffe
		cpu.PC = 0x0100
		hwe.Write(hardware.IO_LCDC, 0x91)
		hwe.Write(hardware.IO_BGP, 0xfc)
		hwe.Write(hardware.IO_OBP0, 0xff)
		hwe.Write(hardware.IO_OBP1, 0xff)
	}

	if cfg.Cartridge != nil {
		hwe.Cartrdige.Load(cfg.Cartridge)
	}

	return gb
}

func (gb *GameBoy) RunFrame() {
	if gb.Paused || gb.ForcedPause {
		return
	}

	maxCyclesPerFrame := gb.ClockSpeed / 60
	frameCycles := 0

	for frameCycles < maxCyclesPerFrame {
		cycles := gb.CPU.ExecuteNextInstruction()
		frameCycles += cycles

		gb.Hardware.Timer.Update(cycles)
		gb.UpdateDisplay(cycles) // todo move to Display
		// todo run interrupts

		if gb.CPU.PC == 0x008f {
			fmt.Println("end bootrom scroll")
			//gb.Paused = true
			//gb.ForcedPause = true
		}

		if gb.CPU.PC == 0x00fe {
			fmt.Println("end bootrom")
			//gb.Paused = true
			//gb.ForcedPause = true
		}

		if gb.CPU.PC == 0x0254 { // tetris write to cartridge
			fmt.Println("debug")
			//gb.Paused = true
			//gb.ForcedPause = true
		}

		if gb.CPU.PC == 0x0286 { // tetris write to not usable ram
			fmt.Println("debug")
			//gb.Paused = true
			//gb.ForcedPause = true
		}
	}

	//fmt.Printf("frame cycles %d\n", frameCycles)
}

func (gb *GameBoy) UpdateDisplay(cycles int) {
	isDisplayEnabled := bits.Test(gb.Hardware.Display.Control, display.ControlDisplayEnabled)

	if !isDisplayEnabled {
		gb.DisplayCycles = 456
		gb.Hardware.Display.CurrentLine = 0
		// mode 1
		gb.Hardware.Display.Status &= 252
		gb.Hardware.Display.Status = bits.Set(gb.Hardware.Display.Status, 0)
		return
	}

	mode := gb.Hardware.Display.Status & 0x3
	nextMode := uint8(0)
	requestInterrupt := false
	mode2Cycles := 456 - 80
	mode3Cycles := mode2Cycles - 172

	switch {
	case gb.Hardware.Display.CurrentLine >= 144:
		nextMode = 1
		gb.Hardware.Display.Status = bits.Set(gb.Hardware.Display.Status, 0)
		gb.Hardware.Display.Status = bits.Clear(gb.Hardware.Display.Status, 1)

		requestInterrupt = bits.Test(gb.Hardware.Display.Status, hardware.STAT_MODE1_VBLANK_INTERRUPT)
	case gb.DisplayCycles >= mode2Cycles:
		// mode 2
		nextMode = 2
		gb.Hardware.Display.Status = bits.Clear(gb.Hardware.Display.Status, 0)
		gb.Hardware.Display.Status = bits.Set(gb.Hardware.Display.Status, 1)

		requestInterrupt = bits.Test(gb.Hardware.Display.Status, hardware.STAT_MODE2_OAM_INTERRUPT)
	case gb.DisplayCycles >= mode3Cycles:
		nextMode = 3
		gb.Hardware.Display.Status = bits.Set(gb.Hardware.Display.Status, 0)
		gb.Hardware.Display.Status = bits.Set(gb.Hardware.Display.Status, 1)

		if nextMode != mode {
			gb.Hardware.Display.DrawLine()
		}
	default:
		nextMode = 0
		gb.Hardware.Display.Status = bits.Clear(gb.Hardware.Display.Status, 0)
		gb.Hardware.Display.Status = bits.Clear(gb.Hardware.Display.Status, 1)

		requestInterrupt = bits.Test(gb.Hardware.Display.Status, hardware.STAT_MODE0_HBLANK_INTERRUPT)

		if nextMode != mode {
			// todo hdma transfer
		}
	}

	if requestInterrupt && mode != nextMode {
		// todo request interrupt 1
	}

	if gb.Hardware.Display.CurrentLine == gb.Hardware.Display.CompareLine {
		gb.Hardware.Display.Status = bits.Set(gb.Hardware.Display.Status, hardware.STAT_COINCIDENCE_FLAG)

		if bits.Test(gb.Hardware.Display.Status, hardware.STAT_COINCIDENCE_INTERRUPT) {
			// todo request interrupt 1
		}
	} else {
		gb.Hardware.Display.Status = bits.Clear(gb.Hardware.Display.Status, hardware.STAT_COINCIDENCE_FLAG)
	}

	gb.DisplayCycles -= cycles

	if gb.DisplayCycles <= 0 {
		gb.Hardware.Display.CurrentLine++

		if gb.Hardware.Display.CurrentLine > 153 {
			gb.Hardware.Display.CurrentLine = 0
		}

		gb.DisplayCycles = 456

		if gb.Hardware.Display.CurrentLine == 144 {
			// todo request interrupt 0
		}
	}
}
