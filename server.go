package main

import (
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"net"
	"sync"
)

type Server struct {
	listenAddress string
	listener      net.Listener
	players       map[*Player]bool
	msgch         chan Message
	scores        map[*Player]uint

	mu sync.Mutex
}

type Player struct {
	net.Conn
	target int
}

func NewServer(listenAddress string) *Server {
	return &Server{
		listenAddress: listenAddress,
		msgch:         make(chan Message),
		players:       make(map[*Player]bool),
		scores:        make(map[*Player]uint),
		mu:            sync.Mutex{},
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
			s.HandleGuess(&msg)
		}
	}
}

func (s *Server) HandleGuess(msg *Message) {
	var guess int

	_, err := fmt.Sscanf(string(msg.payload), "%d", &guess)
	if err != nil {
		fmt.Printf("[%s] has written a non integer value\n", msg.player.RemoteAddr())
		msg.player.Write([]byte("Non integer value :/ Please enter number from 1 to 10 -_-\n"))
		return
	}

	if guess < 1 || guess > 9 {
		fmt.Printf("[%s] has entered number which is not in range from 1 to 10 (%d)\n", msg.player.RemoteAddr(), guess)
		msg.player.Write([]byte("Unlucky :/ Only 1 to 10 numbers are available\n"))
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if guess == msg.player.target {
		fmt.Printf("[%s] got a point!\n", msg.player.RemoteAddr())
		msg.player.Write([]byte("Nice! You got a point! Continue guessing!\n"))

		v := s.scores[msg.player]
		s.scores[msg.player] = v + 1
		msg.player.Write([]byte(fmt.Sprintf("Your balance: %d\n", s.scores[msg.player])))
	} else {
		fmt.Printf("[%s] didn't get a point\n", msg.player.RemoteAddr())
		msg.player.Write([]byte("Unfortunately, you're wrong :/ Upset? Go on!\n"))
	}

	msg.player.target = rand.Intn(10) + 1
	msg.player.Write([]byte(fmt.Sprintf("num: %d\n", msg.player.target)))
}

func (s *Server) AcceptLoop() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			slog.Error("accept loop", "error", err)
		}
		defer conn.Close()

		conn.Write([]byte("Hello, player! Enjoy the game! Try to guess number from 1 to 10! Good luck!\n"))
		player := &Player{
			Conn:   conn,
			target: rand.Intn(10) + 1,
		}
		s.mu.Lock()
		s.players[player] = true
		s.mu.Unlock()
		slog.Info("new connection", "addr", conn.RemoteAddr())

		go s.HandleConnection(player)
	}
}

func (s *Server) HandleConnection(player *Player) {
	buf := make([]byte, 1024)
	for {
		n, err := player.Read(buf)
		if err != nil {
			if err == io.EOF {
				slog.Info("connection lost", "addr", player.RemoteAddr())
				s.mu.Lock()
				delete(s.players, player)
				s.mu.Unlock()
				return
			}

			slog.Error("handle connection", "error", err)
		}

		s.msgch <- Message{
			player:  player,
			payload: buf[:n],
		}
	}
}
