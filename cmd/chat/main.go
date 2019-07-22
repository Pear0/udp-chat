package main

import (
	"bytes"
	"github.com/Pear0/udp-chat/ptypes"
	"github.com/golang/protobuf/proto"
	"log"
	"net"
	"os"
	"time"
)

var magic = []byte("CHTR")

type App struct {
	sendConn     *net.UDPConn
	recvMessages chan *ptypes.BasicMessage
	Name         string
}

func (a *App) Send(msg *ptypes.BasicMessage) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	var payload bytes.Buffer
	_, _ = payload.Write(magic)
	_, _ = payload.Write(data)

	_, err = a.sendConn.Write(payload.Bytes())
	return err
}

func main() {

	// see: https://www.iana.org/assignments/multicast-addresses/multicast-addresses.xhtml#multicast-addresses-3

	// local: 224.0.0.157
	// adhoc 1: 224.0.21.137

	var addrString = "224.0.21.137:22355"
	if addr := os.Getenv("CHAT_ADDR"); addr != "" {
		addrString = addr
	}

	addr, err := net.ResolveUDPAddr("udp4", addrString)
	if err != nil {
		log.Panicln(err)
	}

	var myIfi *net.Interface
	if v := os.Getenv("CHAT_IF"); v != "" {
		myIfi, err = net.InterfaceByName("en0")
		if err != nil {
			log.Panicln(err)
		}
	} else {
		log.Println("Using default interface")
	}

	conn, err := net.ListenMulticastUDP("udp4", myIfi, addr)
	if err != nil {
		log.Panicln(err)
	}

	sendConn, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		log.Panicln(err)
	}

	recvMessages := make(chan *ptypes.BasicMessage)

	go func() {
		var rawBuf [4096]byte

		for {
			bufN, _, err := conn.ReadFromUDP(rawBuf[:])
			if err != nil {
				log.Panicln(err)
			}

			buf := rawBuf[:bufN]
			if !bytes.HasPrefix(buf, magic) {
				continue
			}

			var msg ptypes.BasicMessage
			err = proto.Unmarshal(buf[len(magic):], &msg)
			if err != nil {
				log.Println("failed to parse message", err)
				continue
			}

			recvMessages <- &msg
		}
	}()

	time.Sleep(500 * time.Millisecond)

	guiMain(&App{
		sendConn:     sendConn,
		recvMessages: recvMessages,
		Name:         "Unnamed",
	})

}
