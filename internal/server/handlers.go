package server

import (
	"fmt"
	"gochat/internal/command"
	"strings"
	"time"
)

func (s *Server) handleCommand(c *Client, cmd *command.Command) {
	switch cmd.Name {
	case command.CmdNick:
		s.handleNickCommand(c, cmd)
	case command.CmdSend:
		s.handleSendCommand(c, cmd)
	case command.CmdView:
		s.handleViewCommand(c)
	case command.CmdHelp:
		s.handleHelpCommand(c)
	case command.CmdWho:
		s.handleWhoCommand(c)
	}
}

func (s *Server) handleWhoCommand(c *Client) {
	msg := fmt.Sprintf("Your username is: %s", c.username)
	c.msg(msg)
}

func (s *Server) handleHelpCommand(c *Client) {
	c.sendWelcomeMessage()
}

func (s *Server) handleNickCommand(c *Client, cmd *command.Command) {
	if len(cmd.Args) != 1 {
		errMsg := fmt.Sprintf("Usage: %s%s <username>\n", command.Prefix, command.CmdNick)
		c.msg(errMsg)
		return
	}
	oldUsername := c.username
	newUsername := cmd.Args[0]
	for cl := range s.clients {
		if cl.username == newUsername {
			c.msg(fmt.Sprintf("Username \"%s\" is already taken. Please choose another", newUsername))
		}
	}
	c.username = newUsername
	feedback := fmt.Sprintf("Your username is now %s", newUsername)
	c.msg(feedback)
	s.msgch <- &Message{
		payload: []byte(fmt.Sprintf("%s (%s) is now %s", c.conn.RemoteAddr().String(), oldUsername, newUsername)),
		time:    time.Now(),
		from:    c,
	}
}

func (s *Server) handleSendCommand(c *Client, cmd *command.Command) {
	if len(cmd.Args) < 1 {
		errMsg := fmt.Sprintf("Usage: %s%s <message>\n", command.Prefix, command.CmdSend)
		c.msg(errMsg)
		return
	}
	message := &Message{
		from:    c,
		time:    time.Now(),
		payload: []byte(strings.Join(cmd.Args, " ")),
	}

	s.msgch <- message

	s.histMutex.Lock()
	if len(s.msgHist) >= msgHistSize {
		s.msgHist = s.msgHist[1:]
	}
	s.msgHist = append(s.msgHist, message)
	s.histMutex.Unlock()

}

func (s *Server) handleViewCommand(c *Client) {
	s.histMutex.Lock()
	toSend := make([]*Message, len(s.msgHist))
	copy(toSend, s.msgHist)
	s.histMutex.Unlock()

	if len(toSend) == 0 {
		c.msg("No messages in history yet.")
		return
	}

	c.msg(fmt.Sprintf("--- Displaying last %d messages ---", len(toSend)))

	for _, msg := range toSend {
		formattedMsg := fmt.Sprintf(
			"%s | %s: %s",
			msg.time.Format("02-01 15:04:05"),
			msg.from.username,
			string(msg.payload),
		)
		c.msg(formattedMsg)
	}

	c.msg("--- End of history ---")
}
