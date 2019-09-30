package cpu

import (
	"errors"
	"fmt"
	"github.com/adnsio/gbemu/pkg/gameboy/hardware"
)

type CPU struct {
	A  uint8
	F  *CPUFlags
	B  uint8
	C  uint8
	D  uint8
	E  uint8
	H  uint8
	L  uint8
	SP uint16
	PC uint16

	IsNextInstructionPrefixed bool
	Hardware                  *hardware.Hardware
}

func NewCPU(hwe *hardware.Hardware) *CPU {
	cpu := &CPU{
		F: &CPUFlags{
			Zero:      false,
			Subtract:  false,
			HalfCarry: false,
			Carry:     false,
		},
		Hardware: hwe,
	}

	return cpu
}

/*func (c *CPU) FetchNextInstruction() *Instruction {
	opCode := c.Hardware.ReadCycle(c.PC)
	c.PC++

	if c.IsNextInstructionPrefixed {
		c.IsNextInstructionPrefixed = false
		return c.GetPrefixedInstruction(opCode)
	} else {
		return c.GetInstruction(opCode)
	}
}*/

func (c *CPU) ExecuteNextInstruction() int {
	opCode := c.Hardware.Read(c.PC)

	var inst *Instruction
	if c.IsNextInstructionPrefixed {
		c.IsNextInstructionPrefixed = false
		inst = c.GetPrefixedInstruction(opCode)
	} else {
		inst = c.GetInstruction(opCode)
	}

	return c.ExecuteInstruction(inst)
}

func (c *CPU) ExecuteInstruction(inst *Instruction) int {
	//fmt.Printf("cpu: executing %#02x \"%s\", pc %#04x\n", inst.OpCode, inst.Description, c.PC)

	switch inst.Name {
	case "PREFIX":
		return c.PrefixInst(inst)
	case "NOP":
		return c.NopInst(inst)
	case "LD":
		return c.LdInst(inst)
	case "XOR":
		return c.XorInst(inst)
	case "BIT":
		return c.BitInst(inst)
	case "JR":
		return c.JrInst(inst)
	case "INC":
		return c.IncInst(inst)
	case "CALL":
		return c.CallInst(inst)
	case "PUSH":
		return c.PushInst(inst)
	case "RL":
		return c.RlInst(inst)
	case "RLA":
		return c.RlaInst(inst)
	case "POP":
		return c.PopInst(inst)
	case "DEC":
		return c.DecInst(inst)
	case "RET":
		return c.RetInst(inst)
	case "CP":
		return c.CpInst(inst)
	case "SUB":
		return c.SubInst(inst)
	case "ADD":
		return c.AddInst(inst)
	case "JP":
		return c.JpInst(inst)
	case "DI":
		return c.DiInst(inst)
	case "OR":
		return c.OrInst(inst)
	case "AND":
		return c.AndInst(inst)
	case "CPL":
		return c.CplInst(inst)
	case "EI":
		return c.EiInst(inst)
	case "SWAP":
		return c.SwapInst(inst)
	//case "RST": // todo
	//	return c.RstInst(inst)
	case "SRL":
		return c.SrlInst(inst)
	case "RR":
		return c.RrInst(inst)
	case "RRA":
		return c.RraInst(inst)
	//case "ADC": // todo
	//	return c.AdcInst(inst)
	default:
		panic(errors.New(fmt.Sprintf("cpu error: invalid instruction %#02x \"%s\"", inst.OpCode, inst.Description)))
	}
}

func (c *CPU) FetchImmediate8() uint8 {
	c.PC++
	return c.Hardware.Read(c.PC)
}

func (c *CPU) FetchImmediate16() uint16 {
	lo := c.FetchImmediate8()
	hi := c.FetchImmediate8()
	return (uint16(hi) << 8) | uint16(lo)
}

func (c *CPU) CheckCondition(cond string) bool {
	var res bool

	switch cond {
	case "NZ":
		res = !c.F.Zero
	case "Z":
		res = c.F.Zero
	case "NC":
		res = !c.F.Carry
	case "C":
		res = c.F.Carry
	default:
		panic(errors.New(fmt.Sprintf("cpu error: invalid condition %s", cond)))
	}

	return res
}
