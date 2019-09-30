package cpu

import (
	"errors"
	"fmt"
)

func (c *CPU) ReadParameter(param string) uint8 {
	switch param {
	case "A":
		return c.A
	case "B":
		return c.B
	case "C":
		return c.C
	case "D":
		return c.D
	case "E":
		return c.E
	case "H":
		return c.H
	case "L":
		return c.L
	case "(HL)":
		return c.Hardware.Read(c.ReadHL())
	case "u8":
		return c.FetchImmediate8()
	case "(DE)":
		return c.Hardware.Read(c.ReadDE())
	case "(HL-)":
		hl := c.ReadHL()
		res := c.Hardware.Read(hl)
		c.WriteHL(hl - 1)
		return res
	case "(HL+)":
		hl := c.ReadHL()
		res := c.Hardware.Read(hl)
		c.WriteHL(hl + 1)
		return res
	case "(u16)":
		return c.Hardware.Read(c.FetchImmediate16())
	case "(BC)":
		return c.Hardware.Read(c.ReadBC())
	case "(FF00+u8)":
		return c.Hardware.Read(0xff00 + uint16(c.FetchImmediate8()))
	case "(FF00+C)":
		return c.Hardware.Read(0xff00 + uint16(c.C))
	default:
		panic(errors.New(fmt.Sprintf("cpu error: cannot read parameter 8 %s", param)))
	}
}

func (c *CPU) WriteParameter(param string, val uint8) {
	switch param {
	case "A":
		c.A = val
	case "B":
		c.B = val
	case "C":
		c.C = val
	case "D":
		c.D = val
	case "E":
		c.E = val
	case "H":
		c.H = val
	case "L":
		c.L = val
	case "(HL)":
		c.Hardware.Write(c.ReadHL(), val)
	case "(HL-)":
		hl := c.ReadHL()
		c.Hardware.Write(hl, val)
		c.WriteHL(hl - 1)
	case "(HL+)":
		hl := c.ReadHL()
		c.Hardware.Write(hl, val)
		c.WriteHL(hl + 1)
	case "(FF00+C)":
		c.Hardware.Write(uint16(0xff00)+uint16(c.C), val)
	case "(FF00+u8)":
		c.Hardware.Write(uint16(0xff00)+uint16(c.FetchImmediate8()), val)
	case "(BC)":
		c.Hardware.Write(c.ReadBC(), val)
	case "(u16)":
		u16 := c.FetchImmediate16()
		c.Hardware.Write(u16, val)
	case "(DE)":
		c.Hardware.Write(c.ReadDE(), val)
	default:
		panic(errors.New(fmt.Sprintf("cpu error: cannot write parameter 8 %s", param)))
	}
}

func (c *CPU) ReadParameter16(param string) uint16 {
	switch param {
	case "AF":
		return c.ReadAF()
	case "BC":
		return c.ReadBC()
	case "DE":
		return c.ReadDE()
	case "HL":
		return c.ReadHL()
	case "SP":
		return c.SP
	case "u16":
		return c.FetchImmediate16()
	case "SP+i8":
		return c.SP + uint16(int8(c.FetchImmediate8()))
	case "i8":
		return uint16(int8(c.FetchImmediate8()))
	default:
		panic(errors.New(fmt.Sprintf("cpu error: cannot read parameter 16 %s", param)))
	}
}

func (c *CPU) WriteParameter16(param string, val uint16) {
	switch param {
	case "AF":
		c.WriteAF(val)
	case "BC":
		c.WriteBC(val)
	case "DE":
		c.WriteDE(val)
	case "HL":
		c.WriteHL(val)
	case "SP":
		c.SP = val
	case "(u16)":
		c.Hardware.Write16(c.FetchImmediate16(), val)
	default:
		panic(errors.New(fmt.Sprintf("cpu error: cannot write parameter 16 %s", param)))
	}
}
