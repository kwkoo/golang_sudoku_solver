package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"solver"
)

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/puzzle" {
		w.Header().Set("Content-Type", "application/json")
		grid, err := solver.LoadPuzzle()
		question := grid.Clone()
		if err != nil {
			outputError(w, err)
			return
		}
		if !grid.Solve() {
			outputError(w, errors.New("could not solve puzzle"))
			return
		}

		fmt.Fprintf(w, "{\"question\":\"%s\", \"answer\":\"%s\"}", question, grid)

	}
}

func outputError(w http.ResponseWriter, err error) {
	fmt.Fprint(w, `{"error":`)
	output, _ := json.Marshal(err)
	buffer := bytes.NewBufferString(string(output))
	buffer.WriteTo(w)
	fmt.Fprintln(w, "}")
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
