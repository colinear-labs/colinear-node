package nodeutil

import "strconv"

// function that converts a string to a [32]byte array
func StringToByte32(str string) [32]byte {
	var b [32]byte
	copy(b[:], str)
	return b
}

// function that converts hex string to uint32
func HexStringToUint32(hex string) uint32 {
	i, _ := strconv.ParseUint(hex, 16, 32)
	return uint32(i)
}
