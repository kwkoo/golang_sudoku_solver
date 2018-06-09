package main

import (
	"log"
	"os"
	"solver"
)

func main() {
	grid, err := solver.LoadPuzzle()
	if err != nil {
		log.Fatal(err)
	}
	grid.Print(os.Stdout)
}
