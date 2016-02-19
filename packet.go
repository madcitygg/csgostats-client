package main

import (
	"encoding/json"
	"fmt"
)

type Packet struct {
	Timestamp int64  `json:"timestamp"`
	IP        string `json:"ip"`
	Port      int    `json:"port"`
	Size      int    `json:"size"`
	Payload   string `json:"Payload"`
}

func (p *Packet) Key() string {
	return fmt.Sprintf("%s:%d", p.IP, p.Port)
}

func (p *Packet) JSON() []byte {
	res, _ := json.Marshal(p)
	return res
}
