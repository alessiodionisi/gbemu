package cpu

import (
	"github.com/adnsio/gbemu/pkg/gameboy/hardware"
)

func NewTestCPU(opCode uint8) *CPU {
	hwe := hardware.NewHardware()

	hwe.Cartrdige.Bank[0x0000] = opCode
	hwe.Cartrdige.Bank[0x0001] = 0x01
	hwe.Cartrdige.Bank[0x0002] = 0x02

	cpu := NewCPU(hwe)

	cpu.A = 0x01
	cpu.B = 0x02
	cpu.C = 0x03
	cpu.D = 0x04
	cpu.E = 0x05
	cpu.H = 0x06
	cpu.L = 0x07
	cpu.SP = 0x0102

	return cpu
}
