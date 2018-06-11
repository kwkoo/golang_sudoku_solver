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
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const defaultDelay = 1

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
	updatech := make(chan solver.UpdateEvent)
	delaych := make(chan int)
	go func(ch chan int, c *websocket.Conn) {
		for {
			_, payload, err := c.ReadMessage()
			if err != nil {
				break
			}
			if len(payload) == 0 {
				continue
			}
			p := payload[0]
			if p < 11 {
				delay := int(p)
				ch <- delay
			}
		}
		log.Println("Terminating delay goroutine")
	}(delaych, c)

	var wg sync.WaitGroup
	wg.Add(1)
	go func(ch chan solver.UpdateEvent, c *websocket.Conn) {
		delay := defaultDelay
		message := make([]byte, 2, 2)
		keepgoing := true
		for keepgoing {
			select {
			case event, ok := <-updatech:
				if !ok || event.Index == -1 {
					keepgoing = false
					break
				}
				message[0] = byte(event.Index)
				message[1] = byte(event.Value)
				c.WriteMessage(websocket.BinaryMessage, message)
				time.Sleep(time.Duration(delay*100) * time.Duration(time.Microsecond))
			case delayEvent, ok := <-delaych:
				if !ok {
					keepgoing = false
					break
				}
				delay = delayEvent
				log.Printf("Setting delay to %d", delay)
			}
		}
		log.Println("Terminating update goroutine")
		wg.Done()
	}(updatech, c)

	grid.Solve(updatech)
	log.Println("Done solving the puzzle")
	updatech <- solver.UpdateEvent{Index: -1}

	wg.Wait()
	close(updatech)
	close(delaych)
	log.Println("Request done")
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
