package bits

func Test(value uint8, bit uint8) bool {
	return (value>>bit)&1 == 1
}

func Get(value uint8, bit uint8) uint8 {
	return (value >> bit) & 1
}

func Set(value uint8, bit uint8) uint8 {
	return value | (1 << bit)
}

func Clear(value uint8, bit uint8) uint8 {
	return value & ^(1 << bit)
}
