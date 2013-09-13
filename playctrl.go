package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

type Message struct {
	Version int    `json:"version"`
	Value   string `json:"value"`
}

var CmdNames = map[string]bool{
	"previous":   true,
	"playpause":  true,
	"next":       true,
	"volumeup":   true,
	"volumedown": true,
}

const (
	ProtocolVersion = 1
	msgBufSize      = 1
)

var (
	port = flag.Int("port", 49133, "Port to send the Chrome extension SSE events.")
	sock string

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

func init() {
	user := os.Getenv("USER")
	if user == "" {
		user = "everyone"
	}
	sock = filepath.Join(os.TempDir(), fmt.Sprintf("playctrl-daemon-%s.sock", user))
}

func launchDaemon() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	path := os.Args[0]
	if !filepath.IsAbs(path) {
		path = filepath.Join(cwd, path)
	}
	args := []string{os.Args[0], "daemon", "-port", strconv.Itoa(*port)}
	procattr := &os.ProcAttr{
		Dir:   cwd,
		Env:   os.Environ(),
		Files: []*os.File{nil, nil, nil},
	}
	p, err := os.StartProcess(path, args, procattr)
	if err != nil {
		return err
	}
	return p.Release()
}

func runClient(command string) {
	client, err := rpc.Dial("unix", sock)
	if err != nil {
		// Maybe the server isn't started.
		os.Remove(sock) // Remove the socket if it exists (maybe the server exited uncleanly).
		if err := launchDaemon(); err != nil {
			fatal(err)
		}
		for i := 0; i < 1000; i++ {
			time.Sleep(10 * time.Millisecond)
			client, err = rpc.Dial("unix", sock)
			if err == nil {
				break
			}
		}
		if err != nil {
			fatal("Could not start daemon", err)
		}
		fmt.Println("Starting daemon...")
		// This delay needs to be pretty long to give Chrome a chance to connect.
		// TODO: figure out a better way to do this. We can't just queue up the requests because we don't want the
		// daemon to have anything in its queue if Play isn't actually running.
		time.Sleep(3 * time.Second)
	}

	result := &Nothing{}
	if command == "stop-daemon" {
		if err := client.Call("Server.Shutdown", "", &result); err != nil {
			fatal(err)
		}
		return
	}

	if _, ok := CmdNames[command]; !ok {
		fatal("no such command:", command)
	}
	if err := client.Call("Server.Do", command, &result); err != nil {
		fatal(err)
	}
}

type Server struct {
	quit chan bool
}

type Nothing struct{}

func (s *Server) Shutdown(arg string, reply *Nothing) error {
	s.quit <- true
	return nil
}

func (s *Server) Do(arg string, reply *Nothing) error {
	fmt.Println("command:", arg)
	msg := &Message{
		Version: ProtocolVersion,
		Value:   arg,
	}
	j, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	commands <- j
	return nil
}

func runServer() {
	go processCommands()
	go runHTTPServer()

	s := new(Server)
	rpc.Register(s)
	l, err := net.Listen("unix", sock)
	defer os.Remove(sock)
	if err != nil {
		fatal(err)
	}
	s.quit = make(chan bool)
	conns := make(chan net.Conn)
	go func() {
		for {
			c, err := l.Accept()
			if err != nil {
				fatal(err)
			}
			conns <- c
		}
	}()
	for {
		select {
		case c := <-conns:
			go rpc.ServeConn(c)
		case <-s.quit:
			fmt.Println("Quitting.")
			// Give shutdown RPC time to return normally.
			time.Sleep(10 * time.Millisecond)
			return
		}
	}
}

// multicast
func processCommands() {
	for msg := range commands {
		mu.RLock()
		fmt.Printf("Sending message %s to %d connected client(s).\n", msg, len(clients))
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
	fmt.Println("Client connected.")

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
			fmt.Println("Closing client connection.")
			break loop
		}
	}

	mu.Lock()
	delete(clients, c)
	mu.Unlock()
}

func runHTTPServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleBrowserListener)
	sseAddr := fmt.Sprintf("localhost:%d", *port)
	fmt.Println("Listening on for Chrome extension on", sseAddr)
	fatal(http.ListenAndServe(sseAddr, mux))
}

func usage() {
	fmt.Printf("Usage:\n    $ %s [OPTIONS] COMMAND\nwhere OPTIONS are\n", os.Args[0])
	flag.PrintDefaults()
	commands := []string{"daemon", "stop-daemon"}
	for c := range CmdNames {
		commands = append(commands, c)
	}
	fmt.Printf("and CMD is one of: %v\n", commands)
}

func fatal(args ...interface{}) {
	fmt.Println(args...)
	os.Exit(1)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() < 1 {
		usage()
		fatal("no command provided")
	}
	if flag.Arg(0) == "daemon" {
		runServer()
	} else {
		runClient(flag.Arg(0))
	}
}
