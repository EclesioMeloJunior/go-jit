package aarch64

import (
	"encoding/binary"
	"fmt"
)

const (
	Sf32 byte = 0
	Sf64 byte = 1

	AddImm_OpCode               = 0b00100010
	BranchLink_OpCode           = 0b100101
	BranchReg_OpCode            = 0b1101011000011111000000
	BranchLinkReg_OpCode        = 0b1101011000111111000000
	LoadRegImm_UnsOffset_OpCode = 0b11100101
	Movz_OpCode                 = 0b10100101
	Movk_OpCode                 = 0b11100101
	Ret_Opcode                  = 0b1101011001011111000000
)

type (
	AddImm struct {
		Rd  Register
		Rn  Register
		Imm Imm
	}

	BranchLink struct {
		Imm Imm
	}

	BranchReg struct {
		Rn Register
	}

	BranchLinkReg struct {
		Rn Register
	}

	LoadRegImm_UnsOffset struct {
		Imm Imm
		Rn  Register
		Rt  Register
	}

	Movz struct {
		Rd    Register
		Shift uint32
		Imm   Imm
	}

	Movk struct {
		Rd    Register
		Shift uint32
		Imm   Imm
	}

	Ret struct {
		Rn Register
	}
)

func encodeInstruction(inst uint32) []byte {
	encodedInst := make([]byte, 4)
	binary.LittleEndian.PutUint32(encodedInst, inst)
	return encodedInst
}

func (a *AddImm) Encode(sf byte) []byte {
	inst := uint32(0)

	inst |= (uint32(sf) << 31)
	inst |= (AddImm_OpCode << 23)

	// shift not implemented yet
	// but should be a single bit placed at pos 22
	inst |= (uint32(a.Imm&Imm(Imm12BitMask)) << 10)
	inst |= (uint32(a.Rn&Register(Register5BitMask)) << 5)
	inst |= (uint32(a.Rd & Register(Register5BitMask)))

	return encodeInstruction(inst)
}

func (b *BranchReg) Encode(_ byte) []byte {
	inst := uint32(0)
	inst |= BranchReg_OpCode << 10
	inst |= uint32(b.Rn&Register(Register5BitMask)) << 5
	return encodeInstruction(inst)
}

func (b *BranchLink) Encode(_ byte) []byte {
	inst := uint32(0)
	inst |= (BranchLink_OpCode << 26)
	inst |= uint32(b.Imm & Imm(Imm26BitMask))
	return encodeInstruction(inst)
}

func (b *BranchLinkReg) Encode(_ byte) []byte {
	inst := uint32(0)
	inst |= BranchLinkReg_OpCode << 10
	inst |= uint32(b.Rn&Register(Register5BitMask)) << 5
	return encodeInstruction(inst)
}

func (l *LoadRegImm_UnsOffset) Encode(sf byte) []byte {
	inst := uint32(0)

	inst |= (1 << 31)
	inst |= (uint32(sf) << 30)
	inst |= LoadRegImm_UnsOffset_OpCode << 22
	inst |= (uint32(l.Imm) & Imm12BitMask) << 10
	inst |= (uint32(l.Rn&Register(Register5BitMask)) << 5)
	inst |= uint32(l.Rt & Register(Register5BitMask))

	return encodeInstruction(inst)
}

func genericMov(sf byte, opcode uint32, rd Register, imm Imm, shift uint32) []byte {
	inst := uint32(0)

	inst |= (uint32(sf) << 31)
	inst |= (opcode << 23)

	var encodedShift uint32 = 0
	switch shift {
	case 0:
	case 16:
		encodedShift = 0b01
	case 32:
		if sf != Sf64 {
			panic(fmt.Sprintf("shift %v only for 64-bit", shift))
		}
		encodedShift = 0b10
	case 48:
		if sf != Sf64 {
			panic(fmt.Sprintf("shift %v only for 64-bit", shift))
		}
		encodedShift = 0b11
	default:
		panic(fmt.Sprintf("shift not supported: %v", shift))
	}

	inst |= ((encodedShift & 0b11) << 21)

	// shift not implemented yet
	// but should be a single bit placed at pos 22
	inst |= uint32(imm&Imm(Imm16BitMask)) << 5
	inst |= uint32(rd & Register(Register5BitMask))

	return encodeInstruction(inst)
}

func (m *Movz) Encode(sf byte) []byte {
	return genericMov(sf, Movz_OpCode, m.Rd, m.Imm, m.Shift)
}

func (m *Movk) Encode(sf byte) []byte {
	return genericMov(sf, Movk_OpCode, m.Rd, m.Imm, m.Shift)
}

func (r *Ret) Encode(_ byte) []byte {
	inst := uint32(0)
	inst |= (Ret_Opcode << 10)

	inst |= uint32(r.Rn&Register(Register5BitMask)) << 5
	return encodeInstruction(inst)
}
