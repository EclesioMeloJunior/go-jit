package elf

import (
	"fmt"
	"io"
)

type ELFPolkaVM struct{}

func FromPolkaVm(r io.Reader) (*ELFPolkaVM, error) {
	const headerBytesSize = 16
	header := make([]byte, headerBytesSize)
	n, err := r.Read(header)
	if err != nil {
		return nil, fmt.Errorf("while reading ELF header: %w", err)
	}

	if n != headerBytesSize {
		return nil, fmt.Errorf("expected %d bytes got %d bytes", headerBytesSize, n)
	}

	fmt.Printf("%04x\n", header)
	return nil, nil
}
