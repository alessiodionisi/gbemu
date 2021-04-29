package cpu_test

import (
	"testing"

	"github.com/adnsio/gbemu/pkg/gameboy/cpu"
	"github.com/adnsio/gbemu/pkg/gameboy/cpu/register"
	"github.com/stretchr/testify/assert"
)

func TestGetAndSetRegister16Bit(t *testing.T) {
	dummyRegister := register.SixteenBit(-1)

	c := cpu.New()

	c.SetRegister8Bit(register.A, 0x01)
	c.SetRegister8Bit(register.F, 0x02)

	c.SetRegister8Bit(register.B, 0x03)
	c.SetRegister8Bit(register.C, 0x04)

	c.SetRegister8Bit(register.D, 0x05)
	c.SetRegister8Bit(register.E, 0x06)

	c.SetRegister8Bit(register.H, 0x07)
	c.SetRegister8Bit(register.L, 0x08)

	assert.Equal(t, uint16(0x0102), c.GetRegister16Bit(register.AF))
	assert.Equal(t, uint16(0x0304), c.GetRegister16Bit(register.BC))
	assert.Equal(t, uint16(0x0506), c.GetRegister16Bit(register.DE))
	assert.Equal(t, uint16(0x0708), c.GetRegister16Bit(register.HL))
	assert.Equal(t, uint16(0x0000), c.GetRegister16Bit(dummyRegister))
}

func TestGetAndSetRegister8Bit(t *testing.T) {
	dummyRegister := register.EightBit(-1)

	c := cpu.New()

	c.SetRegister8Bit(register.A, 0x01)
	c.SetRegister8Bit(register.F, 0x02)

	c.SetRegister8Bit(register.B, 0x03)
	c.SetRegister8Bit(register.C, 0x04)

	c.SetRegister8Bit(register.D, 0x05)
	c.SetRegister8Bit(register.E, 0x06)

	c.SetRegister8Bit(register.H, 0x07)
	c.SetRegister8Bit(register.L, 0x08)

	assert.Equal(t, uint8(0x01), c.GetRegister8Bit(register.A))
	assert.Equal(t, uint8(0x02), c.GetRegister8Bit(register.F))
	assert.Equal(t, uint8(0x03), c.GetRegister8Bit(register.B))
	assert.Equal(t, uint8(0x04), c.GetRegister8Bit(register.C))
	assert.Equal(t, uint8(0x05), c.GetRegister8Bit(register.D))
	assert.Equal(t, uint8(0x06), c.GetRegister8Bit(register.E))
	assert.Equal(t, uint8(0x07), c.GetRegister8Bit(register.H))
	assert.Equal(t, uint8(0x08), c.GetRegister8Bit(register.L))
	assert.Equal(t, uint8(0x00), c.GetRegister8Bit(dummyRegister))
}
