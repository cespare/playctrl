package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
)

var (
	port = flag.Int("port", 49132, "Port to send UDP request")

	cmdNames = map[string]bool{
		"previous":   true,
		"playpause":  true,
		"next":       true,
		"volumeup":   true,
		"volumedown": true,
	}
)

type Message struct {
	Version int    `json:"version"`
	Value   string `json:"value"`
}

func usage() {
	fmt.Printf("Usage:\n   $ %s [OPTIONS] CMD\nwhere OPTIONS are\n", os.Args[0])
	flag.PrintDefaults()
	commands := []string{}
	for c := range cmdNames {
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
	if flag.NArg() != 1 {
		usage()
		os.Exit(1)
	}
	cmd := flag.Arg(0)
	if !cmdNames[cmd] {
		usage()
		os.Exit(1)
	}

	addr := fmt.Sprintf("localhost:%d", *port)
	udp, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		fatal(err)
	}
	c, err := net.DialUDP("udp", nil, udp)
	if err != nil {
		fatal(err)
	}
	defer c.Close()

	msg := &Message{
		Version: 1,
		Value: cmd,
	}
	encoder := json.NewEncoder(c)
	if err := encoder.Encode(msg); err != nil {
		fatal(err)
	}
}
