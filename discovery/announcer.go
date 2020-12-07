package discovery

import (
	"encoding/hex"
	"fmt"
	"net"
	"strconv"
	"time"
)

type Announcer struct {
	Players chan int
	Name    chan string
	AChan   chan int
	players string
	name    string
	pc      net.PacketConn
	addr    *net.UDPAddr
}

func (a *Announcer) Open() {
	var err error
	a.AChan = make(chan int)
	a.Name = make(chan string, 1)
	a.Players = make(chan int, 1)
	a.pc, err = net.ListenPacket("udp4", ":47777")
	if err != nil {
		panic(err)
	}
	a.addr, err = net.ResolveUDPAddr("udp4", "255.255.255.255:47777")
	if err != nil {
		panic(err)
	}
}

func (d *Announcer) Close() {
	a.pc.Close()
}

func (d *Announcer) Do() {
	var err error
	for {
		select {
		case <-a.AChan:
			return
		case p := <-a.Players:
			a.players = strconv.Itoa(p)
		case n := <-a.Name:
			a.name = n
		default:
			_, err = a.pc.WriteTo(a.genMessage(a.name, a.players), a.addr)
			if err != nil {
				panic(err)
			}
			time.Sleep(time.Second * 2)
		}
	}
}

func (a Announcer) genMessage(name, players string) []byte {
	msg := fmt.Sprintf("0402%x7e4f70656e7e%02x7e", name, players)

	msgBytes, err := hex.DecodeString(msg)
	if err != nil {
		panic(err)
	}
	return msgBytes
}
