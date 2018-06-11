package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"solver"
	"solver/helper"
)

func handler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	log.Print("Request for URI: ", path)
	if path == "/puzzle" {
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
		return
	}

	if path == "/" || path == "/index.html" {
		fmt.Fprint(w, helper.StaticHTML())
		return
	}

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Not Found"))
}

func outputError(w http.ResponseWriter, err error) {
	fmt.Fprint(w, `{"error":`)
	output, _ := json.Marshal(err)
	buffer := bytes.NewBufferString(string(output))
	buffer.WriteTo(w)
	fmt.Fprintln(w, "}")
}

func main() {
	port := 8080
	flag.IntVar(&port, "port", 8080, "HTTP listener port")
	flag.Parse()

	log.Print("Listening on port ", port)
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
