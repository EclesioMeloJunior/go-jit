package aarch64

import "errors"

var (
	ErrInvalidImmediate = errors.New("invalid immediate")
)

// Imm should not go beyond 12 bits
// using uint16 as the most acurate type
type Imm uint32

const (
	Imm9BitMask  uint32 = (^uint32(0) >> 23) // 0b111111111
	Imm12BitMask uint32 = (^uint32(0) >> 20) // 0b111111111111
	Imm16BitMask uint32 = (^uint32(0) >> 16) // 0b1111111111111111
	Imm26BitMask uint32 = (^uint32(0) >> 6)  // 0b11111111111111111111111111
)

type (
	ZeroRegister byte
	Register     byte
)

func (ZeroRegister) Encode() byte {
	return 0b00011111
}

func (r Register) Encode() byte {
	return byte(r)
}

const Register5BitMask byte = 0b00011111
const (
	R0 Register = iota
	R1
	R2
	R3
	R4
	R5
	R6
	R7
	R8
	R9
	R10
	R11
	R12
	R13
	R14
	R15
	R16
	R17
	R18
	R19
	R20
	R21
	R22
	R23
	R24
	R25
	R26
	R27
	R28
	R29
	R30
)
