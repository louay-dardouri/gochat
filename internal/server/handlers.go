package server

import (
	"fmt"
	"gochat/internal/command"
	"strings"
	"time"
)

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
			msg.time.Format("01-02 15:04:05"),
			msg.from.username,
			string(msg.payload),
		)
		c.msg(formattedMsg)
	}

	c.msg("--- End of history ---")
}

func (s *Server) handleCommand(from *Client, cmd *command.Command) {
	switch cmd.Name {
	case command.CmdNick:
		if len(cmd.Args) != 1 {
			errMsg := fmt.Sprintf("Usage: %s%s <username>\n", command.Prefix, command.CmdNick)
			from.msg(errMsg)
			return
		}
		oldUsername := from.username
		newUsername := cmd.Args[0]
		from.username = newUsername
		feedback := fmt.Sprintf("Your username is now %s", newUsername)
		from.msg(feedback)
		s.msgch <- &Message{
			payload: []byte(fmt.Sprintf("%s (%s) is now %s", from.conn.RemoteAddr().String(), oldUsername, newUsername)),
			time:    time.Now(),
			from:    from,
		}

	case command.CmdSend:
		if len(cmd.Args) < 1 {
			errMsg := fmt.Sprintf("Usage: %s%s <message>\n", command.Prefix, command.CmdSend)
			from.msg(errMsg)
			return
		}
		message := &Message{
			from:    from,
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

	case command.CmdView:
		s.handleViewCommand(from)
	}
}
