package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/jessevdk/go-flags"
	"github.com/mitchellh/colorstring"
)

const (
	version = "0.0.1"
)

type Options struct {
	Bind    string `short:"b" long:"bind" value-name:"IP" description:"Address to bind to" default:"0.0.0.0"`
	Port    int    `short:"p" long:"port" value-name:"PORT" description:"Port to listen on" default:"auto"`
	Target  string `short:"t" long:"target" value-name:"ADDRESS" description:"API endpoints to forward logs to" default:"http://logs.madcity.gg/"`
	Verbose bool   `short:"v" long:"verbose" description:"Verbose output"`
	Version bool   `long:"version" description:"Show version and exit"`
	// LogDirectory string `short:"l" long:"logdir" value-name:"DIRECTORY" description:"Directory to write logs to"`
}

func main() {
	// Parse options
	var options Options
	var parser = flags.NewParser(&options, flags.Default)
	if _, err := parser.Parse(); err != nil {
		os.Exit(1)
	}

	// Show version and exit
	if options.Version {
		colorstring.Println(version)
		os.Exit(0)
	}

	// Test target
	client := &http.Client{}
	resp, err := client.Post(strings.TrimLeft(options.Target, "/")+"/ping", "application/json", nil)
	if err != nil {
		colorstring.Printf("[red]ERROR: Could not ping %s, got: %s\n", options.Target, err.Error())
		os.Exit(1)
	}
	resp.Body.Close()

	// Open socket
	address := fmt.Sprintf("%s:%d", options.Bind, options.Port)
	resolvedAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		colorstring.Printf("[red]ERROR: Could not resolve %s\n", address)
		os.Exit(1)
	}
	socket, err := net.ListenUDP("udp", resolvedAddr)
	if err != nil {
		colorstring.Printf("[red]ERROR: Could not bind to %s\n", address)
		colorstring.Printf("[red]%s\n", err.Error())
		os.Exit(1)
	}
	defer socket.Close()
	socket.SetReadBuffer(1048576)

	// Tell user where to send stuff
	socketAddr := socket.LocalAddr().(*net.UDPAddr)
	log.Printf("Listening on %s:%d\n", options.Bind, socketAddr.Port)

	// Start listener
	listener := NewListener(socket, &options)
	go listener.Listen()

	// Handle SIGINT and SIGTERM.
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	s := <-c
	log.Println("Got signal:", s)
	log.Println("Attempting graceful shutdown")

	// Stop the service gracefully.
	listener.Stop()
	log.Println("Done")
}
