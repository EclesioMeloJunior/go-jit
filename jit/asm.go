package jit

type Instruction interface {
	Encode(sf byte) []byte
}
