package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
)

const msgBufSize = 1

var (
	clientPort = flag.Int("clientport", 49132, "Port to listen for client UDP requests")
	chromePort = flag.Int("chromeport", 49133, "Port to send the Chrome extension SSE events.")
	verbose = flag.Bool("verbose", false, "Turn on verbose logging")

	cmdNames = map[string]bool{
		"previous":   true,
		"playpause":  true,
		"next":       true,
		"volumeup":   true,
		"volumedown": true,
	}

	// Response headers for the SSE request.
	headers = [][2]string{
		{"Content-Type", "text/event-stream"},
		{"Cache-Control", "no-cache"},
		{"Connection", "keep-alive"},
		{"Access-Control-Allow-Origin", "*"},
	}

	commands = make(chan []byte, msgBufSize)

	mu      sync.RWMutex // protects clients
	clients = make(map[chan []byte]bool)
)

type Message struct {
	Version int    `json:"version"`
	Value   string `json:"value"`
}

// errorf writes an error to a net.Conn and also log.
func errorf(c net.Conn, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	log.Println(msg)
	c.Write([]byte(msg))
}

func commandServer(addr *net.UDPAddr) error {
	c, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(c)
	for {
		msg := &Message{}
		if err := decoder.Decode(msg); err != nil {
			errorf(c, "Invalid client request. Error: %s", err)
			continue
		}

		if msg.Version == 0 || msg.Value == "" {
			errorf(c, "Must specify protocol version and message value.")
			continue
		}
		if msg.Version != 1 {
			errorf(c, "Unhandled protocol version: %d", msg.Version)
			continue
		}
		if _, ok := cmdNames[msg.Value]; !ok {
			errorf(c, "Invalid command: '%s'", msg.Value)
			continue
		}

		j, err := json.Marshal(msg)
		if err != nil {
			errorf(c, "Error re-encoding message: %s", err)
			continue
		}

		commands <- j
	}
	return nil
}

// multicast
func processCommands() {
	for msg := range commands {
		mu.RLock()
		log.Printf("Sending message %s to %d connected client(s).", msg, len(clients))
		for c := range clients {
			c <- msg
		}
		mu.RUnlock()
	}
}

// Unification of http.ResponseWriter, http.Flusher, and http.CloseNotifier
type HTTPWriter interface {
	Header() http.Header
	Write([]byte) (int, error)
	WriteHeader(int)
	Flush()
	CloseNotify() <-chan bool
}

// Responds to requests with server-send events.
func handleBrowserListener(writer http.ResponseWriter, r *http.Request) {
	w, ok := writer.(HTTPWriter)
	if !ok {
		panic("HTTP server does not support Flusher and/or CloseNotifier needed for SSE.")
	}
	closed := w.CloseNotify()
	log.Println("Client connected.")

	c := make(chan []byte, msgBufSize)
	mu.Lock()
	clients[c] = true
	mu.Unlock()

	for _, header := range headers {
		w.Header().Set(header[0], header[1])
	}

loop:
	for {
		select {
		case msg := <-c:
			fmt.Fprintf(w, "data:%s\n\n", msg)
			w.Flush()
		case <-closed:
			log.Println("Closing client connection.")
			break loop
		}
	}

	mu.Lock()
	delete(clients, c)
	mu.Unlock()
}

func main() {
	flag.Parse()

	go processCommands()
	udpAddr := fmt.Sprintf("localhost:%d", *clientPort)
	udp, err := net.ResolveUDPAddr("udp", udpAddr)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Listening for UDP client requests on", udp)
	go func() { log.Fatal(commandServer(udp)) }()

	mux := http.NewServeMux()
	mux.HandleFunc("/", handleBrowserListener)
	sseAddr := fmt.Sprintf("localhost:%d", *chromePort)
	log.Println("Listening on for Chrome extension on", sseAddr)
	log.Fatal(http.ListenAndServe(sseAddr, mux))
}
