package cpu

import (
	"github.com/adnsio/gbemu/pkg/gameboy/cpu/register"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type CPU struct {
	GetMemoryFunc func(address uint16) uint8
	SetMemoryFunc func(address uint16, value uint8)

	log zerolog.Logger

	// Accumulator, an 8-bit register for storing data and the results of arithmetic and logical operations.
	a uint8

	// Flag register, consists of 4 flags that are set and reset according to the results of instruction execution.
	// Flags CY and Z are tested by various conditional branch instructions.
	//
	// Z: Set to 1 when the result of an operation is 0; otherwise reset.
	//
	// N: Set to 1 following execution of the substruction instruction, regardless of the result.
	//
	// H: Set to 1 when an operation results in carrying from or borrowing to bit 3.
	//
	// CY: Set to 1 when an operation results in carrying from or borrowing to bit 7.
	f uint8

	// Auxiliary registers, these serve as auxiliary registers to the accumulator.
	// As register pairs (BC, DE, HL), they are 8-bit registers that function as data pointers.
	b, c, d, e, h, l uint8

	// Program counter, a 16-bit register that holds the address data of the program to be executed next.
	// Usually incremented automatically according to the byte count of the fetched instructions.
	// When an instruction with branching is executed, however, immediate data and register contents are loaded.
	pc uint16

	// Stack pointer, a 16-bit register that holds the starting address of the stack area of memory.
	// The contents of the stack pointer are decremented when a subroutine CALL instruction or PUSH instruction is executed or when an interrupt occurs and incremented when a return instruction or pop instruction is executed.
	sp uint16

	prefix bool
}

// func (c *CPU) getA() uint8 {
// 	return bits.GetHigher(c.af)
// }

// func (c *CPU) getF() uint8 {
// 	return uint8(c.af)
// }

// func (c *CPU) getB() uint8 {
// 	return bits.GetHigher(c.bc)
// }

// func (c *CPU) getC() uint8 {
// 	return uint8(c.bc)
// }

// func (c *CPU) getD() uint8 {
// 	return bits.GetHigher(c.de)
// }

// func (c *CPU) getE() uint8 {
// 	return uint8(c.de)
// }

// func (c *CPU) getH() uint8 {
// 	return bits.GetHigher(c.hl)
// }

// func (c *CPU) getL() uint8 {
// 	return uint8(c.hl)
// }

// func (c *CPU) setA(value uint8) {
// 	c.af = bits.SetHigher(c.af, value)
// }

// func (c *CPU) setF(value uint8) {
// 	c.af = bits.SetLower(c.af, value)
// }

// func (c *CPU) setB(value uint8) {
// 	c.bc = bits.SetHigher(c.bc, value)
// }

// func (c *CPU) setC(value uint8) {
// 	c.bc = bits.SetLower(c.bc, value)
// }

// func (c *CPU) setD(value uint8) {
// 	c.de = bits.SetHigher(c.de, value)
// }

// func (c *CPU) setE(value uint8) {
// 	c.de = bits.SetLower(c.de, value)
// }

// func (c *CPU) setH(value uint8) {
// 	c.hl = bits.SetHigher(c.hl, value)
// }

// func (c *CPU) setL(value uint8) {
// 	c.hl = bits.SetLower(c.hl, value)
// }

func (c *CPU) SetFlagZ() {
	// c.af |= 0x0080
}

func (c *CPU) UnsetFlagZ() {
	// c.af &= 0xff7f
}

func (c *CPU) SetFlagZ8Bit(value uint8) {
	if value == 0 {
		c.SetFlagZ()
	} else {
		c.UnsetFlagZ()
	}
}

func (c *CPU) SetFlagN() {
	// c.af |= 0x0040
}

func (c *CPU) UnsetFlagN() {
	// c.af &= 0xffbf
}

func (c *CPU) SetFlagH() {
	// c.af |= 0x0020
}

func (c *CPU) UnsetFlagH() {
	// c.af &= 0xffdf
}

func (c *CPU) SetFlagH8Bit(value uint8) {
	if value&0x10 != 0 {
		c.SetFlagH()
	} else {
		c.UnsetFlagH()
	}
}

func (c *CPU) SetFlagH16Bit(value uint16) {
	if value&0x1000 != 0 {
		c.SetFlagH()
	} else {
		c.UnsetFlagH()
	}
}

func (c *CPU) GetFlagC() bool {
	// if c.af&0x0010 != 0 {
	// 	return true
	// }
	return false
}

func (c *CPU) SetFlagC() {
	// c.af |= 0x0010
}

func (c *CPU) UnsetFlagC() {
	// c.af &= 0xffef
}

func (c *CPU) SetFlagCBool(value bool) {
	if value {
		c.SetFlagC()
	} else {
		c.UnsetFlagC()
	}
}

func (c *CPU) SetFlagC8Bit(value uint8) {
	if value&0x10 != 0 {
		c.SetFlagC()
	} else {
		c.UnsetFlagC()
	}
}

func (c *CPU) SetFlagC16Bit(value uint16) {
	if value&0x1000 != 0 {
		c.SetFlagC()
	} else {
		c.UnsetFlagC()
	}
}

func (c *CPU) GetRegister8Bit(reg register.EightBit) uint8 {
	switch reg {
	case register.A:
		return c.a
	case register.F:
		return c.f
	case register.B:
		return c.b
	case register.C:
		return c.c
	case register.D:
		return c.d
	case register.E:
		return c.e
	case register.H:
		return c.h
	case register.L:
		return c.l
	default:
		return 0
	}
}

func (c *CPU) SetRegister8Bit(reg register.EightBit, value uint8) {
	switch reg {
	case register.A:
		c.a = value
	case register.F:
		c.f = value
	case register.B:
		c.b = value
	case register.C:
		c.c = value
	case register.D:
		c.d = value
	case register.E:
		c.e = value
	case register.H:
		c.h = value
	case register.L:
		c.l = value
	}
}

func (c *CPU) GetRegister16Bit(reg register.SixteenBit) uint16 {
	switch reg {
	case register.AF:
		return uint16(c.a)<<8 | uint16(c.f)
	case register.BC:
		return uint16(c.b)<<8 | uint16(c.c)
	case register.DE:
		return uint16(c.d)<<8 | uint16(c.e)
	case register.HL:
		return uint16(c.h)<<8 | uint16(c.l)
	default:
		return 0
	}
}

func (c *CPU) SetRegister16Bit(reg register.SixteenBit, value uint16) {
	// switch reg {
	// case register.AF:
	// 	c.af = value
	// case register.BC:
	// 	c.bc = value
	// case register.DE:
	// 	c.de = value
	// case register.HL:
	// 	c.hl = value
	// case register.SP:
	// 	c.sp = value
	// }
}

func (c *CPU) Fetch16Bit() uint16 {
	return 0
}

func (c *CPU) Fetch8Bit() uint8 {
	return 0
}

func (c *CPU) GetMemory8Bit(address uint16) uint8 {
	return 0
}

func (c *CPU) SetMemory8Bit(address uint16, value uint8) {

}

func (c *CPU) ExecuteOp(op uint16) {
	opFunc, ok := opUnprefixedMap[op]
	if !ok {
		panic("!ok")
	}

	length, cycles := opFunc(c)

	c.pc += uint16(length)
	_ = cycles
}

func New() *CPU {
	return &CPU{
		log: log.With().Str("package", "cpu").Logger(),
	}
}
