package bits

func GetHigher(value uint16) uint8 {
	return uint8(value >> 8)
}

func GetLower(value uint16) uint8 {
	return uint8(value)
}

func SetHigher(source uint16, value uint8) uint16 {
	lower := source & 0x00f0
	return (uint16(value) << 8) | lower
}

func SetLower(source uint16, value uint8) uint16 {
	higher := source >> 8
	return (higher << 8) | uint16(value)
}
