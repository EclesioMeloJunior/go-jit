package main

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/EclesioMeloJunior/go-riscv/jit"
	"github.com/EclesioMeloJunior/go-riscv/jit/aarch64"
)

func main_v1() {
	hostHelloFn := func(a, b int) {
		fmt.Printf("hello from host! %d\n", a+b)
	}

	hostHelloFnPtr := funcAddr(hostHelloFn)
	//fmt.Printf("0x%x\n", hostHelloFnPtr)

	callGoFnPtr := funcAddr(entrypoint)
	//mov := movUintPtrToReg(aarch64.R0, hostHelloFnPtr)

	instructions := []jit.Instruction{
		&aarch64.Movz{
			Rd:  aarch64.R0,
			Imm: aarch64.Imm(hostHelloFnPtr),
		},
		&aarch64.Movk{
			Rd:    aarch64.R0,
			Imm:   aarch64.Imm(hostHelloFnPtr >> 16),
			Shift: 16,
		},
		&aarch64.Movk{
			Rd:    aarch64.R0,
			Imm:   aarch64.Imm(hostHelloFnPtr >> 32),
			Shift: 32,
		},
		&aarch64.Movk{
			Rd:    aarch64.R0,
			Imm:   aarch64.Imm(hostHelloFnPtr >> 48),
			Shift: 48,
		},
		&aarch64.Movz{
			Rd:  aarch64.R1,
			Imm: aarch64.Imm(0x0a),
		},
		&aarch64.Movz{
			Rd:  aarch64.R2,
			Imm: aarch64.Imm(0x0a),
		},
		&aarch64.BranchLink{
			Imm: aarch64.Imm(0),
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
	pc := uintptr(unsafe.Pointer(&mmapFunc[16]))
	offset := (callGoFnPtr - pc) >> 2
	//fmt.Printf("calculated offset: 0x%x\n", offset)

	bl := &aarch64.BranchLink{
		Imm: aarch64.Imm(offset),
	}
	copy(mmapFunc[24:], bl.Encode(aarch64.Sf64)[:])

	type execFunc func() int
	unsafeFunc := (uintptr)(unsafe.Pointer(&mmapFunc))
	f := *(*execFunc)(unsafe.Pointer(&unsafeFunc))
	MprotectRX(mmapFunc)
	value := f()
	fmt.Printf("0x%x\n", value)
}
