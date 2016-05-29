package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	ForwarderBuffer = 10000
)

type Forwarder struct {
	Target  string
	Channel chan *Packet
	Quit    chan bool
	Client  *http.Client
	verbose bool

	// stats
	PacketsReceived  uint64
	BytesReceived    uint64
	PacketsForwarded uint64
	BytesForwarded   uint64
}

func NewForwarder(target string, verbose bool) Forwarder {
	return Forwarder{
		Target:  target,
		Channel: make(chan *Packet, ForwarderBuffer),
		Quit:    make(chan bool),
		Client:  &http.Client{},
		verbose: verbose,
	}
}

func (f *Forwarder) Start() {
	go func() {
		for {
			select {
			case packet := <-f.Channel:
				// update receive stats
				f.PacketsReceived += uint64(1)
				f.BytesReceived += uint64(packet.Size)

				// build request
				req, _ := http.NewRequest("POST", f.Target, bytes.NewBuffer(packet.JSON()))
				req.Header.Set("Content-Type", "application/json")

				// do request
				resp, err := f.Client.Do(req)
				if err != nil {
					log.Printf("ERROR: send packet %+v to %s: %s\n", packet, f.Target, err.Error())
					continue
				}
				defer resp.Body.Close()

				if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusCreated {
					body, _ := ioutil.ReadAll(resp.Body)
					log.Printf("ERROR: Received unexpected status: %s\n", resp.Status)
					log.Printf(" > Headers: %+v\n", resp.Header)
					log.Printf(" > Body: %s\n", string(body))
					continue
				}

				// update send stats
				f.PacketsForwarded += uint64(1)
				f.BytesForwarded += uint64(packet.Size)

				// print status
				if f.verbose {
					msg := "%d/%d packets received/forwarded"
					stats := fmt.Sprintf(msg, f.PacketsReceived, f.PacketsForwarded)
					log.Printf("Sent %d byte packet from %s:%d to %s (%s)\n", packet.Size, packet.IP, packet.Port, f.Target, stats)
				}
			case <-f.Quit:
				return
			}
		}
	}()
}

func (f *Forwarder) Stop() {
	go func() {
		f.Quit <- true
	}()
}

func (f *Forwarder) Backlog() int {
	return len(f.Channel)
}
