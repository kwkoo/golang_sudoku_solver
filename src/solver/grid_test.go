package solver

import (
	"strconv"
	"testing"
)

func TestSolve(t *testing.T) {
	s := "009060000040010000050700320890400070000507000002009180400000002005000760060200400"
	grid := Grid{}
	r := []rune(s)
	for i := 0; i < 81; i++ {
		value, _ := strconv.Atoi(string(r[i]))
		grid[i] = value
	}
	grid.Solve(nil)
}
