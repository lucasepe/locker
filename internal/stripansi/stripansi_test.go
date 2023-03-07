package stripansi_test

import (
	"fmt"

	"github.com/lucasepe/locker/internal/stripansi"
)

func ExampleStrip() {
	msg := "\x1b[38;5;140m foo\x1b[0m bar"
	cleanMsg := stripansi.Strip(msg)
	fmt.Println(cleanMsg)

	// Output: foo bar
}
