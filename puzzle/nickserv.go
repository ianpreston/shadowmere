package puzzle

import (
	"fmt"
	"strings"
)

type nickservCmdHandler func(string, []string)

type NickServ struct {
	Nick string
	server *Server
	handlers map[string]nickservCmdHandler

	// Store a list, by nick, of users that have identified to NickServ.
	// In a production system, we would probably keep track of all users
	// on the network and whether they are identified. However, in this
	// case, simply tracking identified users is sufficient.
	identified map[string]bool
}

func NewNickserv(server *Server) *NickServ {
	ns := &NickServ{
		Nick: "NickServ",
		server: server,
		identified: make(map[string]bool),
	}
	ns.handlers = map[string]nickservCmdHandler{
		"REGISTER": ns.handleRegister,
		"IDENTIFY": ns.handleIdentify,
	}

	return ns
}

func (ns *NickServ) OnPrivmsg(nick, content string) {
	nick = strings.ToLower(nick)

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

func (ns *NickServ) OnQuit(nick, quitMessage string) {
	ns.unsetIdentified(nick)
}

func (ns *NickServ) OnNickChange(oldNick, newNick string) {
	if ns.getIdentified(oldNick) {
		ns.privmsg(newNick, "You are no longer identified")
	}

	ns.unsetIdentified(oldNick)
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
	if ns.identified[nick] == true {
		ns.privmsg(nick, "You are already identified")
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
		ns.setIdentified(nick)
		ns.privmsg(nick, "You are now identified")
	} else {
		ns.privmsg(nick, "Invalid password for this nick")
	}	
}

func (ns *NickServ) privmsg(recip, message string) {
	ns.server.privmsg(ns.Nick, recip, message)
}

func (ns *NickServ) setIdentified(nick string) {
	ns.identified[nick] = true
}

func (ns *NickServ) unsetIdentified(nick string) {
	delete(ns.identified, nick)
}

func (ns *NickServ) getIdentified(nick string) bool {
	return ns.identified[nick] == true
}