package main

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/EclesioMeloJunior/go-riscv/jit"
	"github.com/EclesioMeloJunior/go-riscv/jit/aarch64"
)

type Ctx struct {
	Regs [13]uint32
}

var ctx = &Ctx{
	Regs: [13]uint32{},
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	fmt.Printf("ctx: {%+v}\n", ctx)

	ctxPtr := uintptr(unsafe.Pointer(ctx))

	instructions := []jit.Instruction{
		&aarch64.BranchLink{
			Imm: aarch64.Imm(10),
		},

		// moves the ctx addr to r3
		&aarch64.Movz{
			Rd:  aarch64.R3,
			Imm: aarch64.Imm(ctxPtr),
		},
		&aarch64.Movk{
			Rd:    aarch64.R3,
			Imm:   aarch64.Imm(ctxPtr >> 16),
			Shift: 16,
		},
		&aarch64.Movk{
			Rd:    aarch64.R3,
			Imm:   aarch64.Imm(ctxPtr >> 32),
			Shift: 32,
		},
		&aarch64.Movk{
			Rd:    aarch64.R3,
			Imm:   aarch64.Imm(ctxPtr >> 48),
			Shift: 48,
		},

		// lets try changing its value
		// mov 0xf to R4
		&aarch64.Movz{
			Rd:  aarch64.R4,
			Imm: aarch64.Imm(0xff),
		},

		// storing a byte but to the field Value
		&aarch64.Strb{
			Rt:  aarch64.R4,     // value we want to place there
			Rn:  aarch64.R3,     // ctx pointer
			Imm: aarch64.Imm(0), // dont offset by any value
		},

		// storing a byte but to the field Other
		&aarch64.Strb{
			Rt:  aarch64.R4,     // value we want to place there
			Rn:  aarch64.R3,     // ctx pointer
			Imm: aarch64.Imm(4), // skip the first field and update the next
		},

		&aarch64.Ret{
			Rn: aarch64.R30, // deafult value
		},
	}

	asm := make([]byte, 0)
	for _, ins := range instructions {
		encoded := ins.Encode(aarch64.Sf64)
		asm = append(asm, encoded...)
	}

	//fmt.Println(asm)
	//fmt.Println(len(asm))

	mmapFunc, err := syscall.Mmap(
		-1,
		0,
		len(asm),
		syscall.PROT_READ|syscall.PROT_WRITE,
		// MAP_ANON is available only for darwin, for linux use syscall.MAP_ANONYMOUS
		syscall.MAP_PRIVATE|syscall.MAP_ANON,
	)

	if err != nil {
		panic(err)
	}

	copy(mmapFunc, asm)

	type execFunc func() int
	unsafeFunc := (uintptr)(unsafe.Pointer(&mmapFunc))
	f := *(*execFunc)(unsafe.Pointer(&unsafeFunc))
	MprotectRX(mmapFunc)

	_ = f()

	fmt.Printf("ctx: {%+v}\n", ctx)
}
