package cpu

import (
	"errors"
	"fmt"
	"strconv"
)

func (c *CPU) PrefixInst(inst *Instruction) int {
	c.IsNextInstructionPrefixed = true
	c.PC++

	return inst.CyclesBranch
}

func (c *CPU) NopInst(inst *Instruction) int {
	c.PC++

	return inst.CyclesBranch
}

func (c *CPU) LdInst(inst *Instruction) int {
	switch inst.Bits {
	case 8:
		srcVal := c.ReadParameter(inst.Parameters[1])
		c.WriteParameter(inst.Parameters[0], srcVal)
	case 16:
		srcVal := c.ReadParameter16(inst.Parameters[1])
		c.WriteParameter16(inst.Parameters[0], srcVal)

		if inst.Parameters[1] == "SP+i8" {
			c.F.Zero = false
			c.F.Subtract = false
			c.F.HalfCarry = srcVal>>3 != 0x0
			c.F.Carry = srcVal>>7 != 0x0
		}
	default:
		panic(errors.New(fmt.Sprintf("cpu error: invalid bits %d", inst.Bits)))
	}

	c.PC++

	return inst.CyclesBranch
}

func (c *CPU) XorInst(inst *Instruction) int {
	srcVal := c.ReadParameter(inst.Parameters[1])
	res := c.A ^ srcVal

	c.F.Zero = res == 0
	c.F.Subtract = false
	c.F.HalfCarry = false
	c.F.Carry = false

	c.A = res

	c.PC++

	return inst.CyclesBranch
}

func (c *CPU) BitInst(inst *Instruction) int {
	bit, err := strconv.Atoi(inst.Parameters[0])
	if err != nil {
		panic(err)
	}

	ubit := uint8(bit)
	srcVal := c.ReadParameter(inst.Parameters[1])
	test := (srcVal >> ubit) & 0x1

	c.F.Zero = test == 0x0
	c.F.Subtract = false
	c.F.HalfCarry = true

	c.PC++

	return inst.CyclesBranch
}

func (c *CPU) JrInst(inst *Instruction) int {
	var cond bool

	if len(inst.Parameters) == 1 {
		cond = true
	} else {
		cond = c.CheckCondition(inst.Parameters[0])
	}

	if cond {
		var srcParam string

		if len(inst.Parameters) == 1 {
			srcParam = inst.Parameters[0]
		} else {
			srcParam = inst.Parameters[1]
		}

		c.PC += c.ReadParameter16(srcParam)
	} else {
		c.PC++
	}

	c.PC++

	return inst.CyclesBranch
}

func (c *CPU) IncInst(inst *Instruction) int {
	switch inst.Bits {
	case 8:
		srcVal := c.ReadParameter(inst.Parameters[0])
		res := srcVal + 1

		c.F.Zero = res == 0x0
		c.F.Subtract = false
		c.F.HalfCarry = res&0xf == 0xf

		c.WriteParameter(inst.Parameters[0], res)
	case 16:
		srcVal := c.ReadParameter16(inst.Parameters[0])
		res := srcVal + 1

		c.WriteParameter16(inst.Parameters[0], res)
	default:
		panic(errors.New(fmt.Sprintf("cpu error: invalid bits %d", inst.Bits)))
	}

	c.PC++

	return inst.CyclesBranch
}

func (c *CPU) CallInst(inst *Instruction) int {
	var cond bool

	if len(inst.Parameters) == 1 {
		cond = true
	} else {
		cond = c.CheckCondition(inst.Parameters[0])
	}

	if cond {
		var srcParam string

		if len(inst.Parameters) == 1 {
			srcParam = inst.Parameters[0]
		} else {
			srcParam = inst.Parameters[1]
		}

		srcVal := c.ReadParameter16(srcParam)
		nextPC := c.PC + 1

		c.Hardware.Write(c.SP-2, uint8(nextPC&0xff))
		c.Hardware.Write(c.SP-1, uint8(nextPC>>8))

		c.SP -= 2
		c.PC = srcVal

		return inst.CyclesBranch
	} else {
		c.PC++

		return inst.CyclesNoBranch
	}
}

func (c *CPU) PushInst(inst *Instruction) int {
	srcVal := c.ReadParameter16(inst.Parameters[0])

	c.SP -= 2
	c.Hardware.Write16(c.SP, srcVal)

	c.PC++

	return inst.CyclesBranch
}

func (c *CPU) RlInst(inst *Instruction) int {
	srcVal := c.ReadParameter(inst.Parameters[0])

	val := srcVal << 1

	if c.F.Carry {
		val = val | 0x1
	}

	c.F.Zero = val == 0x0
	c.F.Subtract = false
	c.F.HalfCarry = false
	c.F.Carry = srcVal>>7 != 0x0

	c.WriteParameter(inst.Parameters[0], val)

	c.PC++

	return inst.CyclesBranch
}

func (c *CPU) RlaInst(inst *Instruction) int {
	val := c.A << 1

	if c.F.Carry {
		val = val | 0x1
	}

	c.F.Zero = false
	c.F.Subtract = false
	c.F.HalfCarry = false
	c.F.Carry = c.A>>7 != 0x0

	c.A = val

	c.PC++

	return inst.CyclesBranch
}

func (c *CPU) PopInst(inst *Instruction) int {
	val := c.Hardware.Read16(c.SP)

	c.WriteParameter16(inst.Parameters[0], val)
	c.SP += 2

	c.PC++

	return inst.CyclesBranch
}

func (c *CPU) DecInst(inst *Instruction) int {
	switch inst.Bits {
	case 8:
		srcVal := c.ReadParameter(inst.Parameters[0])
		res := srcVal - 1

		c.F.Zero = res == 0x0
		c.F.Subtract = true
		c.F.HalfCarry = res&0xf == 0xf

		c.WriteParameter(inst.Parameters[0], res)
	case 16:
		srcVal := c.ReadParameter16(inst.Parameters[0])
		res := srcVal - 1

		c.WriteParameter16(inst.Parameters[0], res)
	default:
		panic(errors.New(fmt.Sprintf("cpu error: invalid bits %d", inst.Bits)))
	}

	c.PC++

	return inst.CyclesBranch
}

func (c *CPU) RetInst(inst *Instruction) int {
	var cond bool

	if len(inst.Parameters) == 1 {
		cond = c.CheckCondition(inst.Parameters[0])
	} else {
		cond = true
	}

	if cond {
		c.PC = c.Hardware.Read16(c.SP)
		c.SP += 2

		return inst.CyclesBranch
	} else {
		c.PC++

		return inst.CyclesNoBranch
	}
}

func (c *CPU) CpInst(inst *Instruction) int {
	tarVal := c.ReadParameter(inst.Parameters[0])
	srcVal := c.ReadParameter(inst.Parameters[1])
	res := tarVal - srcVal

	c.F.Zero = res == 0x0
	c.F.Subtract = true
	c.F.HalfCarry = res&0xf == 0xf
	c.F.Carry = res>>7 != 0x0

	c.PC++

	return inst.CyclesBranch
}

func (c *CPU) SubInst(inst *Instruction) int {
	tarVal := c.ReadParameter(inst.Parameters[0])
	srcVal := c.ReadParameter(inst.Parameters[1])
	res := tarVal - srcVal

	c.F.Zero = res == 0x0
	c.F.Subtract = true
	c.F.HalfCarry = res&0xf < 0xf
	c.F.Carry = res < 0

	c.WriteParameter(inst.Parameters[0], res)

	c.PC++

	return inst.CyclesBranch
}

func (c *CPU) AddInst(inst *Instruction) int {
	switch inst.Bits {
	case 8:
		tarVal := c.ReadParameter(inst.Parameters[0])
		srcVal := c.ReadParameter(inst.Parameters[1])
		res := tarVal + srcVal

		c.F.Zero = res == 0x0
		c.F.Subtract = false
		c.F.HalfCarry = res&0xf > 0xf
		c.F.Carry = res > 0xff

		c.WriteParameter(inst.Parameters[0], res)
	case 16:
		tarVal := c.ReadParameter16(inst.Parameters[0])
		srcVal := c.ReadParameter16(inst.Parameters[1])
		res := tarVal + srcVal

		c.F.Zero = res == 0x0
		c.F.Subtract = false
		c.F.HalfCarry = res&0xf > 0xf
		c.F.Carry = res > 0xff

		c.WriteParameter16(inst.Parameters[0], res)
	default:
		panic(errors.New(fmt.Sprintf("cpu error: invalid bits %d", inst.Bits)))
	}

	c.PC++

	return inst.CyclesBranch
}

func (c *CPU) JpInst(inst *Instruction) int {
	var cond bool

	if len(inst.Parameters) == 1 {
		cond = true
	} else {
		cond = c.CheckCondition(inst.Parameters[0])
	}

	if cond {
		var srcParam string

		if len(inst.Parameters) == 1 {
			srcParam = inst.Parameters[0]
		} else {
			srcParam = inst.Parameters[1]
		}

		srcVal := c.ReadParameter16(srcParam)

		c.PC = srcVal

		return inst.CyclesBranch
	} else {
		c.PC++

		return inst.CyclesNoBranch
	}
}

func (c *CPU) DiInst(inst *Instruction) int {
	// todo disable interrupts

	c.PC++

	return inst.CyclesBranch
}

func (c *CPU) OrInst(inst *Instruction) int {
	srcVal := c.ReadParameter(inst.Parameters[1])
	res := c.A | srcVal

	c.F.Zero = res == 0
	c.F.Subtract = false
	c.F.HalfCarry = false
	c.F.Carry = false

	c.A = res

	c.PC++

	return inst.CyclesBranch
}

func (c *CPU) AndInst(inst *Instruction) int {
	srcVal := c.ReadParameter(inst.Parameters[1])
	res := c.A | srcVal

	c.F.Zero = res == 0
	c.F.Subtract = false
	c.F.HalfCarry = true
	c.F.Carry = false

	c.A = res

	c.PC++

	return inst.CyclesBranch
}

func (c *CPU) CplInst(inst *Instruction) int {
	c.A = ^c.A

	c.F.Subtract = true
	c.F.HalfCarry = true

	c.PC++

	return inst.CyclesBranch
}

func (c *CPU) EiInst(inst *Instruction) int {
	// todo enable interrupts

	c.PC++

	return inst.CyclesBranch
}

func (c *CPU) SwapInst(inst *Instruction) int {
	val := c.ReadParameter(inst.Parameters[0])
	res := (val&0xF0)>>4 ^ ((val & 0x0F) << 4)

	c.WriteParameter(inst.Parameters[0], val)

	c.F.Zero = res == 0
	c.F.Subtract = false
	c.F.HalfCarry = false
	c.F.Carry = false

	c.PC++

	return inst.CyclesBranch
}

func (c *CPU) RstInst(inst *Instruction) int {
	// todo
	switch inst.Parameters[0] {
	case "28h":

	}

	c.PC++

	return inst.CyclesBranch
}

func (c *CPU) SrlInst(inst *Instruction) int {
	val := c.ReadParameter(inst.Parameters[0])
	res := val >> 1

	c.WriteParameter(inst.Parameters[0], res)

	c.F.Zero = res == 0
	c.F.Subtract = false
	c.F.HalfCarry = false
	c.F.Carry = res&0x01 != 0

	c.PC++

	return inst.CyclesBranch
}

func (c *CPU) RrInst(inst *Instruction) int {
	val := c.ReadParameter(inst.Parameters[0])

	var ci uint8
	if c.F.Carry {
		ci = 1
	} else {
		ci = 0
	}

	res := (val >> 1) | (ci << 7)

	c.WriteParameter(inst.Parameters[0], res)

	c.F.Zero = res == 0
	c.F.Subtract = false
	c.F.HalfCarry = false
	c.F.Carry = val&0x01 != 0

	c.PC++

	return inst.CyclesBranch
}

func (c *CPU) RraInst(inst *Instruction) int {
	val := c.A

	var ci uint8
	if c.F.Carry {
		ci = 1
	} else {
		ci = 0
	}

	res := (val >> 1) | (ci << 7)

	c.A = res

	c.F.Zero = false
	c.F.Subtract = false
	c.F.HalfCarry = false
	c.F.Carry = val&0x01 != 0

	c.PC++

	return inst.CyclesBranch
}

func (c *CPU) AdcInst(inst *Instruction) int {
	// todo

	c.PC++

	return inst.CyclesBranch
}
