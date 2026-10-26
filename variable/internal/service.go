package internal

import (
	"fmt"
	"os"
)

var (
	// Stdin is an open file, pointing to the standard input
	Stdin *os.File
	// PrintFormatted is used to print a formatted string to standard output
	PrintFormatted func(format string, a ...any) (n int, err error)
)

func init() {
	Stdin = os.Stdin
	PrintFormatted = fmt.Printf
}
