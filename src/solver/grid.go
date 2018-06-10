package solver

import (
	"fmt"
	"io"
	"strconv"
	"strings"
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
			fmt.Fprintln(w, "---+---+----")
		}
	}
}

func (grid Grid) String() string {
	var b strings.Builder
	for _, value := range grid {
		b.WriteString(strconv.Itoa(value))
	}
	return b.String()
}

// Solve will keep running till it finds a solution to the puzzle. Returns true if successful, false if there is a problem.
func (grid *Grid) Solve() bool {
	index := grid.nextEmptyCellFromIndex(0)
	if index == -1 {
		return true
	}

	s := newStack()
	c := newCellContext(index, grid.nextEmptyCellFromIndex(index+1), grid.candidatesForCell(index))
	s.push(c)

	for s.hasMore() {
		c, _ = s.pop()
		if c.hasMoreCandidates() {
			candidate := c.nextCandidate()
			grid[c.index] = candidate
			if c.nextEmpty == -1 {
				return true
			}
			s.push(c)
			s.push(newCellContext(c.nextEmpty, grid.nextEmptyCellFromIndex(c.nextEmpty+1), grid.candidatesForCell(c.nextEmpty)))
		} else {
			// unsuccessful - so we'll reset the cell to empty
			grid[c.index] = 0
		}
	}

	return false
}

// Clone produces a copy of grid.
func (grid Grid) Clone() Grid {
	target := grid
	return target
}

// returns -1 if there are no empty cells left
func (grid Grid) nextEmptyCellFromIndex(index int) int {
	for i := index; i < 81; i++ {
		if grid[i] == 0 {
			return i
		}
	}
	return -1
}

func (grid Grid) candidatesForCell(index int) []int {
	if index > 80 {
		return []int{}
	}
	return intersectingCandidates(
		candidatesFromDigits(grid.digitsInRow(index)),
		candidatesFromDigits(grid.digitsInColumn(index)),
		candidatesFromDigits(grid.digitsInBox(index)))
}

func (grid Grid) digitsInRow(index int) []int {
	var digits []int
	var value int
	y := index / 9
	for i := y * 9; i < (y+1)*9; i++ {
		value = grid[i]
		if value != 0 {
			digits = append(digits, value)
		}
	}

	return digits
}

func (grid Grid) digitsInColumn(index int) []int {
	var digits []int
	var value int
	x := index % 9
	for i := x; i < 81; i += 9 {
		value = grid[i]
		if value != 0 {
			digits = append(digits, value)
		}
	}

	return digits
}

func (grid Grid) digitsInBox(index int) []int {
	var digits []int
	var value int
	row := index / 9
	column := index % 9

	// find index of top-left cell in box
	i := (row/3)*27 + (column/3)*3

	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			value = grid[i]
			i++
			if value != 0 {
				digits = append(digits, value)
			}
		}
		i += 6
	}

	return digits
}

// returns inverse of argument
func candidatesFromDigits(digits []int) []int {
	var taken [9]bool
	var candidates []int
	for _, i := range digits {
		taken[i-1] = true
	}
	for i, v := range taken {
		if !v {
			candidates = append(candidates, i+1)
		}
	}

	return candidates
}

func intersectingCandidates(x, y, z []int) []int {
	var count [9]int
	for _, value := range x {
		count[value-1]++
	}
	for _, value := range y {
		count[value-1]++
	}
	for _, value := range z {
		count[value-1]++
	}

	var candidates []int
	for i, v := range count {
		if v == 3 {
			candidates = append(candidates, i+1)
		}
	}

	return candidates
}
