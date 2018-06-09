package solver

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

const (
	timeout = 30
	url     = "http://davidbau.com/generated/sudoku.txt"
)

// Grid represents the Sudoku grid, with 0 representing an empty cell.
type Grid [9][9]int

// LoadPuzzle loads a new grid from a puzzle generator.
func LoadPuzzle() (Grid, error) {
	grid := Grid{}

	client := http.Client{Timeout: time.Duration(timeout * time.Second)}
	resp, err := client.Get(url)
	if err != nil {
		return grid, fmt.Errorf("could not load puzzle from %s: %v", url, err)
	}

	defer resp.Body.Close()
	scanner := bufio.NewScanner(resp.Body)
	var value int
	for i := 0; i < 18 && scanner.Scan(); i++ {
		// we're only interested in odd lines
		if i%2 == 0 {
			continue
		}
		runes := []rune(scanner.Text())
		if len(runes) < 36 {
			return grid, fmt.Errorf("expected a line of length of at least 36 runes - got %d instead", len(runes))
		}
		for j := 2; j < 9*4; j += 4 {
			r := runes[j]
			if r == ' ' {
				value = 0
			} else {
				value, err = strconv.Atoi(string(r))
				if err != nil {
					return grid, fmt.Errorf("%v is not a digit: %v", r, err)
				}
			}
			grid[(j-2)/4][(i-1)/2] = value
		}
	}

	return grid, nil
}

// Print prints a grid to the writer.
func (grid Grid) Print(w io.Writer) {
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			value := grid[j][i]
			if value == 0 {
				fmt.Fprint(w, ".")
			} else {
				fmt.Fprint(w, value)
			}
		}
		fmt.Fprintln(w, "")
	}
}
