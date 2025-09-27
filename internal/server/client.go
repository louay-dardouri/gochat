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

	c.sendWelcomeMessage()

	sc := bufio.NewScanner(c.conn)

	for sc.Scan() {
		msg := sc.Text()
		cmd, err := command.Parse(msg)

		if err != nil {
			errorMsg := fmt.Sprintf("Error: %s\n", err.Error())
			c.msg(errorMsg)
			continue
		}

		if cmd != nil {
			c.server.handleCommand(c, cmd)
		}

	}
}

func (c *Client) msg(txt string) {
	c.conn.Write([]byte(txt + "\n"))
}

func (c *Client) sendWelcomeMessage() {
	welcomeText := `
Welcome to the Go Chat Server!
---------------------------------
You are connected as: %s
Available commands are:
  /%s <username>   - Change your username.
  /%s <message>    - Send a message to everyone in the room.
  (or just type a message to send it)
  /%s 			   - Show a list of last messages sent (chat log)
  /%s 			   - Prints current username
  /%s 			   - Shows a help message (this message)
---------------------------------
`
	formattedText := fmt.Sprintf(
		welcomeText,
		c.username,
		command.CmdNick,
		command.CmdSend,
		command.CmdView,
		command.CmdWho,
		command.CmdHelp,
	)

	c.msg(formattedText)
}
