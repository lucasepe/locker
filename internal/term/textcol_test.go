package term_test

import (
	"os"

	"github.com/lucasepe/locker/internal/term"
)

func ExamplePrintColumns() {
	items := []string{
		"📂 folder thing",
		"won't get tripped up by emojis",
		"or even color codes",
		"here",
		"are",
		"some",
		"shorter",
		"lines",
		"running out of stuff",
		"foo bar",
	}
	// pass pointer to array of strings and a margin value. this will ensure at least 4 spaces appear to the right of each cell
	term.PrintColumns(os.Stdout, &items, 4)

	// Output: 📂 folder thing                    some
	// won't get tripped up by emojis    shorter
	// or even color codes               lines
	// here                              running out of stuff
	// are                               foo bar
}
