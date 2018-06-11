package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"solver"
	"solver/helper"
	"strings"

	"github.com/gorilla/websocket"
)

// number of microseconds between update events
//const delay = 200

var upgrader = websocket.Upgrader{}

func handler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	log.Print("Request for URI: ", path)
	if path == "/puzzle" {
		w.Header().Set("Content-Type", "application/json")
		grid, err := solver.LoadPuzzle()
		if err != nil {
			outputError(w, err)
			return
		}
		fmt.Fprintf(w, "{\"puzzle\":\"%s\"}", grid)
		return
	}

	if strings.HasPrefix(path, "/solve/") {
		puzzle := path[len("/solve/"):]
		if len(puzzle) != 81 {
			outputError(w, fmt.Errorf("puzzle did not have the expected length of 81 - received %d instead", len(puzzle)))
			return
		}
		grid, err := solver.NewGridFromString(puzzle)
		if err != nil {
			outputError(w, fmt.Errorf("could not convert %s to Grid object: %v", puzzle, err))
			return
		}
		handleSolveRequest(w, r, grid)
		return
	}

	if path == "/" || path == "/index.html" {
		helper.StaticHTML(w)
		return
	}

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Not Found"))
}

func handleSolveRequest(w http.ResponseWriter, r *http.Request, grid solver.Grid) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		outputError(w, fmt.Errorf("could not upgrade to websocket: %v", err))
		return
	}
	defer c.Close()
	ch := make(chan solver.UpdateEvent)
	go grid.Solve(ch)
	message := make([]byte, 2, 2)
	for event, ok := <-ch; ok; event, ok = <-ch {
		message[0] = byte(event.Index)
		message[1] = byte(event.Value)
		c.WriteMessage(websocket.BinaryMessage, message)
		//time.Sleep(time.Duration(delay * time.Microsecond))
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
	port := 8080
	flag.IntVar(&port, "port", 8080, "HTTP listener port")
	flag.Parse()

	log.Print("Listening on port ", port)
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
