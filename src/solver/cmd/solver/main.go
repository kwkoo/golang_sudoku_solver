package main

import (
	"fmt"
	"log"
	"os"
	"solver"
)

func main() {
	grid, err := solver.LoadPuzzle()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Before")
	fmt.Println("------")
	fmt.Println("")
	grid.Print(os.Stdout)
	fmt.Println("")

	if !grid.Solve() {
		log.Fatal("Could not solve puzzle")
	}

	fmt.Println("After")
	fmt.Println("-----")
	fmt.Println("")
	grid.Print(os.Stdout)
}
