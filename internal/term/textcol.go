package term

import (
	"fmt"
	"io"
	"math"
	"strings"
	"unicode/utf8"

	"github.com/lucasepe/locker/internal/stripansi"
	"golang.org/x/term"
)

func PrintColumns(wri io.Writer, strs *[]string, margin int) {
	// get the longest string the columns need to contain
	maxLength := 0
	marginStr := strings.Repeat(" ", margin)
	// also keep track of each individual length to easily calculate padding
	lengths := []int{}
	for _, str := range *strs {
		colorless := stripansi.Strip(str)
		// len() is insufficient here, as it counts emojis as 4 characters each
		length := utf8.RuneCountInString(colorless)
		maxLength = max(maxLength, length)
		lengths = append(lengths, length)
	}

	// see how wide the terminal is
	width := getTermWidth()
	// calculate the dimensions of the columns
	numCols, numRows := calculateTableSize(width, margin, maxLength, len(*strs))

	// if we're forced into a single column, fall back to simple printing (one per line)
	if numCols == 1 {
		for _, str := range *strs {
			fmt.Fprintln(wri, str)
		}
		return
	}

	// `i` will be a left-to-right index. this will need to get converted to a top-to-bottom index
	for i := 0; i < numCols*numRows; i++ {
		// treat output like a "table" with (x, y) coordinates as an intermediate representation
		// first calculate (x, y) from i
		x, y := rowIndexToTableCoords(i, numCols)
		// then convery (x, y) to `j`, the top-to-bottom index
		j := tableCoordsToColIndex(x, y, numRows)

		// try to access the array, but the table might have more cells than array elements, so only try to access if within bounds
		strLen := 0
		str := ""
		if j < len(lengths) {
			strLen = lengths[j]
			str = (*strs)[j]
		}

		// calculate the amount of padding required
		numSpacesRequired := maxLength - strLen
		spaceStr := strings.Repeat(" ", numSpacesRequired)

		// print the item itself
		fmt.Fprint(wri, str)

		// if we're at the last column, print a line break
		if x+1 == numCols {
			fmt.Fprintf(wri, "\n")
		} else {
			fmt.Fprint(wri, spaceStr)
			fmt.Fprint(wri, marginStr)
		}
	}
}

func getTermWidth() int {
	if !term.IsTerminal(0) {
		return 80
	}

	width, _, err := term.GetSize(0)
	check(err)
	return width
}

func calculateTableSize(width, margin, maxLength, numCells int) (int, int) {
	numCols := (width + margin) / (maxLength + margin)
	if numCols == 0 {
		numCols = 1
	}
	numRows := int(math.Ceil(float64(numCells) / float64(numCols)))
	return numCols, numRows
}

func rowIndexToTableCoords(i, numCols int) (int, int) {
	x := i % numCols
	y := i / numCols
	return x, y
}

func tableCoordsToColIndex(x, y, numRows int) int {
	return y + numRows*x
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
