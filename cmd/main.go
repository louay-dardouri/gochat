package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

type Client struct {
	username string
	conn     net.Conn
	server   *Server
}

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

func NewClient(conn net.Conn, s *Server) *Client {
	return &Client{
		username: "anon",
		conn:     conn,
		server:   s,
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

func (c *Client) readLoop() {
	defer func() {
		c.server.leaving <- c
		c.conn.Close()
	}()

	sc := bufio.NewScanner(c.conn)

	for sc.Scan() {
		msg := sc.Text()
		formattedMsg := fmt.Sprintf("%s: %s", c.username, msg)

		c.server.msgch <- &Message{
			from:    c,
			payload: []byte(formattedMsg),
			time:    time.Now().Format("2006-01-02 15:04:05"),
		}
	}
}

func main() {
	server := NewServer(":8080")
	log.Fatal(server.Start())
}
