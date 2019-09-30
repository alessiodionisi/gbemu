package cpu

import "github.com/adnsio/gbemu/pkg/gameboy/bits"

type CPUFlags struct {
	Zero      bool
	Subtract  bool
	HalfCarry bool
	Carry     bool
}

func (f *CPUFlags) Read() uint8 {
	res := uint8(0x00)

	if f.Zero {
		res = bits.Set(res, 7)
	}

	if f.Subtract {
		res = bits.Set(res, 6)
	}

	if f.HalfCarry {
		res = bits.Set(res, 5)
	}

	if f.Carry {
		res = bits.Set(res, 4)
	}

	return res
}

func (f *CPUFlags) Write(val uint8) {
	f.Zero = bits.Test(val, 7)
	f.Subtract = bits.Test(val, 6)
	f.HalfCarry = bits.Test(val, 5)
	f.Carry = bits.Test(val, 4)
}

func (c *CPU) ReadAF() uint16 {
	return (uint16(c.A) << 8) | uint16(c.F.Read())
}

func (c *CPU) WriteAF(val uint16) {
	c.A = uint8(val >> 8)
	c.F.Write(uint8(val))
}

func (c *CPU) ReadBC() uint16 {
	return (uint16(c.B) << 8) | uint16(c.C)
}

func (c *CPU) WriteBC(val uint16) {
	c.B = uint8(val >> 8)
	c.C = uint8(val)
}

func (c *CPU) ReadDE() uint16 {
	return (uint16(c.D) << 8) | uint16(c.E)
}

func (c *CPU) WriteDE(val uint16) {
	c.D = uint8(val >> 8)
	c.E = uint8(val)
}

func (c *CPU) ReadHL() uint16 {
	return (uint16(c.H) << 8) | uint16(c.L)
}

func (c *CPU) WriteHL(val uint16) {
	c.H = uint8(val >> 8)
	c.L = uint8(val)
}
