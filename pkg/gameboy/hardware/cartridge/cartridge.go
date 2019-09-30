package cartridge

import (
	"bytes"
	"fmt"
)

const (
	Start = 0x0000
	End   = 0x7fff

	BankStart = 0x0000
	BankEnd   = 0x3fff
	BankSize  = BankEnd - BankStart + 1

	SwitchableBankStart = 0x4000
	SwitchableBankEnd   = 0x7fff
	SwitchableBankSize  = SwitchableBankEnd - SwitchableBankStart + 1

	RamStart = 0xa000
	RamEnd   = 0xbfff
	RamSize  = RamEnd - RamStart + 1
)

type Cartridge struct {
	Bank            [BankSize]uint8
	SwitchableBank0 [SwitchableBankSize]uint8
	Ram             [RamSize]uint8
	Title           string
}

func NewCartridge() *Cartridge {
	crt := &Cartridge{}

	return crt
}

func (c *Cartridge) Load(data []uint8) {
	titleVal := data[0x0134:0x0143]
	cgbFlag := data[0x0143]
	sgbFlag := data[0x0146]
	carType := data[0x0147]
	romSize := data[0x0148]
	ramSize := data[0x0149]

	c.Title = string(bytes.Trim(titleVal, "\x00"))

	fmt.Printf("cartridge: title %s, cgb %#02x, sgb %#02x, type %#02x, rom %#02x, ram %#02x\n", c.Title, cgbFlag, sgbFlag, carType, romSize, ramSize)

	if carType != 0 {
		panic(fmt.Errorf("cartridge: unimplemented type %#02x", carType))
	}

	if romSize != 0 {
		panic(fmt.Errorf("cartridge: unimplemented rom %#02x", romSize))
	}

	if ramSize != 0 {
		panic(fmt.Errorf("cartridge: unimplemented ram %#02x", ramSize))
	}

	for i := BankStart; i <= BankEnd; i++ {
		c.Bank[i] = data[i]
	}

	for i := SwitchableBankStart; i <= SwitchableBankEnd; i++ {
		c.SwitchableBank0[i-SwitchableBankStart] = data[i]
	}
}

func (c *Cartridge) Read(addr uint16) uint8 {
	switch {
	case addr >= BankStart && addr <= BankEnd:
		return c.Bank[addr]
	case addr >= SwitchableBankStart && addr <= SwitchableBankEnd:
		return c.SwitchableBank0[addr-SwitchableBankStart]
	case addr >= RamStart && addr <= RamEnd:
		return c.Ram[addr-RamStart]
	default:
		panic(fmt.Errorf("cartridge: invalid address %#04x", addr))
	}
}
