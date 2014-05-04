package shadowmere

import (
	"fmt"
	"strings"
)

// TODO - Extract common between NickServ/ChanServ

type chanservCmdHandler func(string, []string)

type ChanServ struct {
	Nick string
	mere *Services

	handlers map[string]chanservCmdHandler
}

func NewChanserv(mere *Services) *ChanServ {
	cs := &ChanServ{
		Nick: "ChanServ",
		mere: mere,
	}
	cs.handlers = map[string]chanservCmdHandler{
		"REGISTER": cs.handleRegister,
		"OP": cs.handleOp,
		"DEOP": cs.handleDeop,
	}

	return cs
}

func (cs *ChanServ) OnPrivmsg(nick, content string) {
	nick = strings.ToLower(nick)

	tokecs := strings.Split(content, " ")
	command := strings.ToUpper(tokecs[0])
	var args []string
	if len(tokecs) > 1 {
		args = tokecs[1:]
	}

	h := cs.handlers[command]
	if h != nil {
		h(nick, args)
	} else {
		cs.notice(nick, fmt.Sprintf("No such command: %s", command))
	}
}

func (cs *ChanServ) handleRegister(nick string, args []string) {
	cs.notice(nick, "Register")
}

func (cs *ChanServ) handleOp(nick string, args []string) {
	cs.notice(nick, "op")
}

func (cs *ChanServ) handleDeop(nick string, args []string) {
	cs.notice(nick, "deop")
}

func (cs *ChanServ) notice(recip, message string) {
	cs.c().notice(cs.Nick, recip, message)
}

func (cs *ChanServ) c() *Connection {
	return cs.mere.connection
}

func (cs *ChanServ) r() *RegisteredNickRepo {
	return cs.mere.datastore.RegisteredNicks
}