package main

import (
	"fmt"
	"log/slog"
	"net"
)

type Server struct {
	listenAddress string
	listener      net.Listener
	players       map[string]net.Conn
	msgch         chan Message
}

func NewServer(listenAddress string) *Server {
	return &Server{
		listenAddress: listenAddress,
	}
}

func (s *Server) Start() error {
	var err error
	s.listener, err = net.Listen("tcp", s.listenAddress)
	if err != nil {
		slog.Error("server start", "error", err)
		return fmt.Errorf("server start error: %s", err)
	}

	go s.AcceptLoop()

	for {
		select {
		case msg := <-s.msgch:
			fmt.Println(string(msg.payload))
		}
	}

	return nil
}

func (s *Server) AcceptLoop() error {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			slog.Error("accept loop", "error", err)
			return fmt.Errorf("accept loop error: %s", err)
		}

		fmt.Println("new connection", conn.RemoteAddr())

		go s.HandleConnection(conn)
	}
}

func (s *Server) HandleConnection(conn net.Conn) error {
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			slog.Error("handle connection", "error", err)
			return fmt.Errorf("handle connection error: %s", err)
		}

		s.msgch <- Message{
			sender:  conn,
			payload: buf[:n],
		}
	}
}
