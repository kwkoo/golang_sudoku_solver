package solver

import (
	"bufio"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const (
	timeout = 30
	url     = "http://davidbau.com/generated/sudoku.txt"
)

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
	i := 0
	for line := 0; line < 18 && scanner.Scan(); line++ {
		// we're only interested in odd lines
		if line%2 == 0 {
			continue
		}
		runes := []rune(scanner.Text())
		if len(runes) < 36 {
			return grid, fmt.Errorf("expected a line of length of at least 36 runes - got %d instead", len(runes))
		}
		for col := 2; col < 9*4; col += 4 {
			r := runes[col]
			if r == ' ' {
				value = 0
			} else {
				value, err = strconv.Atoi(string(r))
				if err != nil {
					return grid, fmt.Errorf("%v is not a digit: %v", r, err)
				}
			}
			grid[i] = value
			i++
		}
	}

	return grid, nil
}
