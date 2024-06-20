package main

import "net"

type Message struct {
	sender  net.Conn
	payload []byte
}
