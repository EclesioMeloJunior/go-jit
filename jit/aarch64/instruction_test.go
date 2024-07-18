package aarch64

import (
	"testing"

	"github.com/EclesioMeloJunior/go-riscv/jit"
)

func TestAddImm64Encoding(t *testing.T) {
	instructions := []jit.Instruction{
		&Movz{
			Rd:  R1,
			Imm: Imm(0x02),
		},
		&AddImm{
			Rd:  R0,
			Rn:  R1,
			Imm: Imm(0x05),
		},
	}

	asm := make([]byte, 0)
	for _, ins := range instructions {
		encoded := ins.Encode(Sf64)
		asm = append(asm, encoded...)
	}
}
