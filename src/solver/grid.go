package solver

import (
	"fmt"
	"io"
)

// Grid represents the Sudoku grid, with 0 representing an empty cell.
type Grid [81]int

// Print prints a grid to the writer.
func (grid Grid) Print(w io.Writer) {
	var value int
	for i := 0; i < 81; i++ {
		value = grid[i]
		if value == 0 {
			fmt.Fprint(w, ".")
		} else {
			fmt.Fprint(w, value)
		}
		if i%3 == 2 && i%9 != 8 {
			fmt.Fprint(w, "|")
		}
		if i%9 == 8 {
			fmt.Fprintln(w, "")
		}
		if i%27 == 26 && i != 80 {
			fmt.Fprintln(w, "------------")
		}
	}
}
