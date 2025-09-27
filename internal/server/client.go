package server

import (
	"bufio"
	"fmt"
	"gochat/internal/command"
	"net"
)

type Client struct {
	username string
	conn     net.Conn
	server   *Server
}

func NewClient(conn net.Conn, s *Server) *Client {
	return &Client{
		username: "anon",
		conn:     conn,
		server:   s,
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
		cmd, err := command.Parse(msg)

		if err != nil {
			errorMsg := fmt.Sprintf("Error: %s\n", err.Error())
			c.conn.Write([]byte(errorMsg))
			continue
		}

		if cmd != nil {
		}

	}
}
