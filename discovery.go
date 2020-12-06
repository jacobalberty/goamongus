package goamongus

import (
	"encoding/hex"
	"fmt"
	"net"
	"strconv"
	"time"
)

type Discovery struct {
	Players chan int
	Name    chan string
	AChan   chan int
	players string
	name    string
	pc      net.PacketConn
	addr    *net.UDPAddr
}

func (d *Discovery) Open() {
	var err error
	d.AChan = make(chan int)
	d.Name = make(chan string, 1)
	d.Players = make(chan int, 1)
	d.pc, err = net.ListenPacket("udp4", ":47777")
	if err != nil {
		panic(err)
	}
	d.addr, err = net.ResolveUDPAddr("udp4", "255.255.255.255:47777")
	if err != nil {
		panic(err)
	}
}

func (d *Discovery) Close() {
	d.pc.Close()
}

func (d *Discovery) Do() {
	var err error
	for {
		select {
		case <-d.AChan:
			return
		case p := <-d.Players:
			d.players = strconv.Itoa(p)
		case n := <-d.Name:
			d.name = n
		default:
			_, err = d.pc.WriteTo(d.genMessage(d.name, d.players), d.addr)
			if err != nil {
				panic(err)
			}
			time.Sleep(time.Second * 2)
		}
	}
}

func (d *Discovery) genMessage(name, players string) []byte {
	msg := fmt.Sprintf("0402%x7e4f70656e7e%02x7e", name, players)

	msgBytes, err := hex.DecodeString(msg)
	if err != nil {
		panic(err)
	}
	return msgBytes
}
