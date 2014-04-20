package puzzle

import (
	"fmt"
	"strings"
)

type nickservCmdHandler func(string, []string)

type NickServ struct {
	Service
	server *Server
	handlers map[string]nickservCmdHandler
}

func NewNickserv(server *Server) *NickServ {
	ns := &NickServ{
		server: server,
	}
	ns.handlers = map[string]nickservCmdHandler{
		"REGISTER": ns.handleRegister,
		"IDENTIFY": ns.handleIdentify,
	}

	return ns
}

func (ns *NickServ) Nick() string {
	return "NickServ"
}

func (ns *NickServ) OnPrivmsg(nick, content string) {
	tokens := strings.Split(content, " ")
	command := tokens[0]
	var args []string
	if len(tokens) > 1 {
		args = tokens[1:]
	}

	h := ns.handlers[command]
	if h != nil {
		h(nick, args)
	} else {
		ns.privmsg(nick, fmt.Sprintf("No such command: %s", command))
	}
}

func (ns *NickServ) handleIdentify(nick string, args []string) {
	ns.privmsg(
		nick,
		fmt.Sprintf("IDENTIFY with arguments: %s", strings.Join(args, ",")),
	)
}

func (ns *NickServ) handleRegister(nick string, args []string) {
	ns.privmsg(
		nick,
		fmt.Sprintf("REGISTER with arguments: %s", strings.Join(args, ",")),
	)
}

func (ns *NickServ) privmsg(recip, message string) {
	line := fmt.Sprintf(
		":%s PRIVMSG %s :%s\r\n",
		ns.Nick(),
		recip,
		message,
	)
	ns.server.write(line)
}