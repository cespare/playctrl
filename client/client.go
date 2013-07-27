package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"

	"playctrl"
)

var port = flag.Int("port", 49132, "Port to send UDP request")

func usage() {
	fmt.Printf("Usage:\n   $ %s [OPTIONS] CMD\nwhere OPTIONS are\n", os.Args[0])
	flag.PrintDefaults()
	commands := []string{}
	for c := range playctrl.CmdNames {
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
	if !playctrl.CmdNames[cmd] {
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

	msg := &playctrl.Message{
		Version: playctrl.ProtocolVersion,
		Value:   cmd,
	}
	encoder := json.NewEncoder(c)
	if err := encoder.Encode(msg); err != nil {
		fatal(err)
	}
}
