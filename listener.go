package main

import (
	"log"
	"net"
	"sync"
	"time"
)

type Listener struct {
	Socket     *net.UDPConn
	Target     string
	Forwarders map[string]Forwarder
	Quit       chan bool
	waitGroup  *sync.WaitGroup
	verbose    bool

	// stats
	PacketsReceived uint64
	BytesReceived   uint64
}

func NewListener(socket *net.UDPConn, options *Options) *Listener {
	l := &Listener{
		Socket:     socket,
		Target:     options.Target,
		Forwarders: make(map[string]Forwarder),
		Quit:       make(chan bool),
		waitGroup:  &sync.WaitGroup{},
		verbose:    options.Verbose,
	}
	l.waitGroup.Add(1)
	return l
}

func (l *Listener) Listen() {
	defer l.waitGroup.Done()

	// read from socket
	var buffer [1400]byte
	for {
		// Listen for a quit signal
		select {
		case <-l.Quit:
			for key, forwarder := range l.Forwarders {
				log.Printf("Stopping forwarder for %s\n", key)
				forwarder.Stop()
			}
			l.Socket.Close()
			return
		default:
		}

		// set timeout
		err := l.Socket.SetReadDeadline(time.Now().Add(3 * time.Second))
		if err != nil {
			panic(err)
		}

		// read from socket
		n, clientAddr, err := l.Socket.ReadFromUDP(buffer[0:])
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
			panic(err)
		}

		// create packet structure
		packet := &Packet{
			Timestamp: time.Now().UnixNano(),
			IP:        clientAddr.IP.String(),
			Port:      clientAddr.Port,
			Size:      n,
			Payload:   string(buffer[0:n]),
		}

		// update stats
		l.PacketsReceived += uint64(1)
		l.BytesReceived += uint64(n)

		// send to the right forwarder
		key := packet.Key()
		_, ok := l.Forwarders[key]
		if !ok {
			log.Printf("Starting new forwarder for %s -> %s", key, l.Target)
			forwarder := NewForwarder(l.Target, l.verbose)
			forwarder.Start()
			l.Forwarders[key] = forwarder
		}
		l.Forwarders[key].Channel <- packet
	}
}

func (l *Listener) Stop() {
	l.Quit <- true
	l.waitGroup.Wait()
}
