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

func (ns *NickServ) handleRegister(nick string, args []string) {
	ns.privmsg(
		nick,
		fmt.Sprintf("REGISTER with arguments: %s", strings.Join(args, ",")),
	)
}

func (ns *NickServ) handleIdentify(nick string, args []string) {
	if len(args) != 1{
		ns.privmsg(nick, "You must specify a password")
		return
	}

	rn, err := ns.server.datastore.GetRegisteredNick(nick)
	if err != nil {
		// TODO - Handle error better
		fmt.Errorf(err.Error())
	}
	if rn == nil {
		ns.privmsg(nick, "This nickname is not registered")
		return
	}
	
	validPassword := ns.server.datastore.Authenticate(rn, args[0])
	if validPassword {
		ns.privmsg(nick, "You are now identified")
	} else {
		ns.privmsg(nick, "Invalid password for this nick")
	}	
}

func (ns *NickServ) privmsg(recip, message string) {
	ns.server.privmsg(ns.Nick(), recip, message)
}
