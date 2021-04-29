package register

type EightBit int
type SixteenBit int

const (
	A EightBit = iota
	F
	B
	C
	D
	E
	H
	L
)

const (
	AF SixteenBit = iota
	BC
	DE
	HL
	PC
	SP
)
