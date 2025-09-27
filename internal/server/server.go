package server

import (
	"fmt"
	"net"
)

type Message struct {
	from    *Client
	payload []byte
	time    string
}

type Server struct {
	listenAddr string
	ln         net.Listener
	quitch     chan struct{}
	msgch      chan *Message

	clients map[*Client]bool
	joinch  chan *Client
	leaving chan *Client
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		quitch:     make(chan struct{}),
		msgch:      make(chan *Message, 256),
		clients:    make(map[*Client]bool),
		joinch:     make(chan *Client),
		leaving:    make(chan *Client),
	}
}

func (s *Server) hub() {
	for {
		select {
		case client := <-s.joinch:
			s.clients[client] = true
			fmt.Printf("New client joined: %s (%s)\n", client.username, client.conn.RemoteAddr())
		case client := <-s.leaving:
			fmt.Printf("Client left: %s (%s)\n", client.username, client.conn.RemoteAddr())
			delete(s.clients, client)

		case msg := <-s.msgch:
			fmt.Printf("%s | %s\n", msg.time, string(msg.payload))
		}
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()
	s.ln = ln

	go s.hub()

	go s.acceptLoop()

	<-s.quitch
	close(s.msgch)

	return nil
}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			fmt.Println("err: ", err)
			continue
		}
		client := NewClient(conn, s)
		s.joinch <- client

		go client.readLoop()
	}
}
