package shadowmere

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

var MalformedMessage = errors.New("Malformed message")

func (srv *Connection) parseMessage(line string) (string, string, []string, error) {
	var command string
	var origin string
	var args []string

	line = strings.Trim(line, " \r\n")
	tokens := srv.splitMessage(line)

	if line[0] == ':' {
		if len(tokens) < 3 {
			return "", "", nil, MalformedMessage
		}
		origin = strings.TrimPrefix(tokens[0], ":")
		command = tokens[1]
		args = tokens[2:]
	} else {
		command = strings.ToUpper(tokens[0])
		if len(tokens) > 1 {
			args = tokens[1:]
		}
	}

	return command, origin, args, nil
}

func (srv *Connection) splitMessage(line string) []string {
	tokens := strings.Split(line, " ")

	rightTokenIdx := -1
	for idx, t := range tokens {
		if strings.HasPrefix(t, ":") && idx != 0 {
			rightTokenIdx = idx
			break
		}
	}

	if rightTokenIdx != -1 {
		leftTokens := tokens[:rightTokenIdx]
		rightToken := strings.Join(tokens[rightTokenIdx:], " ")
		rightToken = strings.TrimPrefix(rightToken, ":")

		tokens = leftTokens
		tokens = append(tokens, rightToken)
	}

	return tokens
}

func (srv *Connection) privmsg(from, to, message string) {
	cmd := fmt.Sprintf(":%s PRIVMSG %s :%s\r\n", from, to, message)
	srv.write(cmd)
}

func (srv *Connection) notice(from, to, message string) {
	cmd := fmt.Sprintf(":%s NOTICE %s :%s\r\n", from, to, message)
	srv.write(cmd)
}

func (srv *Connection) nick(nick, user, host, real string) {
	// NICK <nick> <hops> <ts> <modes> <user> <host> <server> :<real>
	cmd := fmt.Sprintf(
		"NICK %s 1 %v +i %s %s %s :%s\r\n",
		nick,
		time.Now().Unix(),
		user,
		host,
		srv.name,
		real,
	)
	srv.write(cmd)
}

func (srv *Connection) pong(origin string) {
	srv.write(fmt.Sprintf("PONG :%s\r\n", origin))
}

func (srv *Connection) svsmode(origin, nick, modes string) {
	srv.write(fmt.Sprintf(":%s SVS2MODE %s :%s\r\n", origin, nick, modes))
}

func (srv *Connection) svsnick(origin, old, new string) {
	cmd := fmt.Sprintf(
		":%s SVSNICK %s %s :%v\r\n",
		origin,
		old,
		new,
		time.Now().Unix(),
	)
	srv.write(cmd)
}

func (srv *Connection) svskill(origin, nick, reason string) {
	srv.write(fmt.Sprintf(":%s h %s :%s\r\n", origin, nick, reason))
}

func (srv *Connection) chghost(origin, nick, vhost string) {
	srv.write(fmt.Sprintf(":%s AL %s :%s\r\n", origin, nick, vhost))
}
