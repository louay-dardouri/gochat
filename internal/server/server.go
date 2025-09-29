package server

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

type Message struct {
	from    *Client
	payload []byte
	time    time.Time
}

type Server struct {
	listenAddr string
	ln         net.Listener
	quitch     chan struct{}
	msgch      chan *Message

	clients      map[*Client]bool
	clientsMutex sync.RWMutex
	joinch       chan *Client
	leaving      chan *Client

	msgHist   []*Message
	histMutex sync.Mutex
}

const (
	msgHistSize = 256
)

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		quitch:     make(chan struct{}),
		msgch:      make(chan *Message, 256),
		clients:    make(map[*Client]bool),
		joinch:     make(chan *Client),
		leaving:    make(chan *Client),

		msgHist: make([]*Message, 0, msgHistSize),
	}
}

func (s *Server) hub() {
	for {
		select {
		case client := <-s.joinch:
			s.clientsMutex.Lock()
			s.clients[client] = true
			s.clientsMutex.Unlock()
			fmt.Printf("New client joined: %s (%s)\n", client.username, client.conn.RemoteAddr().String())

		case client := <-s.leaving:
			fmt.Printf("Client left: %s (%s)\n", client.username, client.conn.RemoteAddr())
			s.clientsMutex.Lock()
			delete(s.clients, client)
			s.clientsMutex.Unlock()

		case msg := <-s.msgch:
			formattedMsg := fmt.Sprintf(
				"%s | %s: %s\n",
				msg.time.Format("2006-01-02 15:04:05"),
				msg.from.username,
				string(msg.payload),
			)

			fmt.Print(formattedMsg)

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

	log.Printf("TCP Chat Server started, listening on %s", s.listenAddr)

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
