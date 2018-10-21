package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"solver"
	"solver/helper"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const defaultDelay = 1
const bufferSize = 50
const staticBufferSize = 4096
const staticFilename = "debug.html"

var upgrader = websocket.Upgrader{}
var debug = false

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
		staticContent(w)
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
	updatech := make(chan solver.UpdateEvent, bufferSize-1)
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
		var message [bufferSize * 2]byte
		keepgoing := true
		count := 1
		for keepgoing {
			select {
			case event, ok := <-updatech:
				if !ok || event.Index == -1 {
					log.Print("got -1 at 1")
					keepgoing = false
					break
				}
				message[0] = byte(event.Index)
				message[1] = byte(event.Value)
				count = 1
				// more messages in the queue - let's fill it up
				depth := len(updatech)
				if depth > 0 {
					for i := count; i <= depth; i++ {
						event, ok := <-updatech
						if !ok || event.Index == -1 {
							log.Print("got -1 at 2")
							keepgoing = false
							break
						} else {
							count++
							message[i*2] = byte(event.Index)
							message[i*2+1] = byte(event.Value)
						}
					}
				}
				c.WriteMessage(websocket.BinaryMessage, message[:count*2])
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

func staticContent(w io.Writer) {
	if !debug {
		helper.StaticHTML(w)
		return
	}

	// We're in debug mode - dump contents of debug.html.
	f, err := os.Open(staticFilename)
	if err != nil {
		log.Printf("Error opening static HTML %s: %v", staticFilename, err)
		fmt.Fprintf(w, "Error opening static HTML %s: %v", staticFilename, err)
		return
	}
	defer f.Close()
	buf := make([]byte, staticBufferSize, staticBufferSize)
	for {
		n, err := f.Read(buf)
		if n > 0 {
			w.Write(buf[:n])
		}
		if err != nil {
			if err != io.EOF {
				log.Printf("IO Error while reading static HTML %s: %v", staticFilename, err)
			}
			break
		}
	}
}

func main() {
	port := 0
	portenv := os.Getenv("PORT")
	if len(portenv) > 0 {
		port, _ = strconv.Atoi(portenv)
	}
	if port == 0 {
		port = 8080
		flag.IntVar(&port, "port", port, "HTTP listener port")
	}

	debugenv := os.Getenv("DEBUG")
	if len(debugenv) > 0 {
		debug = true
	} else {
		flag.BoolVar(&debug, "debug", debug, "Use debug.html instead of the hardcoded statichtml.go file.")
	}

	flag.Parse()

	if debug {
		log.Print("Debug mode on")
	}
	log.Print("Listening on port ", port)
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
