package shadowmere

import (
	"../kenny"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type nickservCmdHandler func(string, []string)

type NickServ struct {
	Nick string
	mere *Services

	handlers map[string]nickservCmdHandler

	rand *rand.Rand
}

func NewNickserv(mere *Services) *NickServ {
	ns := &NickServ{
		Nick: "NickServ",
		mere: mere,
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
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
	command := strings.ToUpper(tokens[0])
	var args []string
	if len(tokens) > 1 {
		args = tokens[1:]
	}

	h := ns.handlers[command]
	if h != nil {
		h(nick, args)
	} else {
		ns.notice(nick, fmt.Sprintf("No such command: %s", command))
	}
}

func (ns *NickServ) OnQuit(nick, quitMessage string) {
	ns.unidentifyUser(nick)
}

func (ns *NickServ) OnNewNick(nick string) {
	ns.handleRegisteredNick(nick)
}

func (ns *NickServ) OnNickChange(oldNick, newNick string) {
	ns.unidentifyUser(oldNick)
	ns.c().svsmode(ns.Nick, newNick, "-r")

	ns.handleRegisteredNick(newNick)
}

func (ns *NickServ) handleRegister(nick string, args []string) {
	if len(args) != 2 {
		ns.notice(nick, "Syntax: REIGSTER <email> <password>")
		return
	}

	rn, err := ns.r().GetByNick(nick)
	if err != nil {
		kenny.ErrorErr(err)
		return
	}
	if rn != nil {
		ns.notice(nick, "That nickname is already registered.")
		return
	}

	newRn := &RegisteredNick{
		Nick:   nick,
		Email:  args[0],
		Passwd: args[1],
	}
	err = ns.r().Register(newRn)
	if err != nil {
		ns.notice(nick, "Error registering nickname")
		kenny.ErrorErr(err)
		return
	}

	ns.notice(nick, "Your nickname is now registered!")
	ns.identifyUser(nick)
}

func (ns *NickServ) handleIdentify(nick string, args []string) {
	if len(args) != 1 {
		ns.notice(nick, "You must specify a password")
		return
	}
	if ns.mere.curstate.IsIdentified(nick) {
		ns.notice(nick, "You are already identified")
		return
	}

	rn, err := ns.r().GetByNick(nick)
	if err != nil {
		kenny.ErrorErr(err)
		return
	}
	if rn == nil {
		ns.notice(nick, "This nickname is not registered")
		return
	}

	validPassword := ns.r().Authenticate(rn, args[0])
	if validPassword {
		ns.identifyUser(nick)
	} else {
		ns.notice(nick, "Invalid password for this nick")
	}
}

func (ns *NickServ) handleRegisteredNick(nick string) {
	rn, err := ns.r().GetByNick(nick)
	if err != nil {
		kenny.ErrorErr(err)
		return
	}
	if rn == nil {
		return
	}

	ns.notice(nick, "Your nickname is registered. You have 60 seconds to identify or change your nick.")

	ns.mere.curstate.NewNick(nick)
	ns.enforceIdentifiedNick(nick)
}

func (ns *NickServ) enforceIdentifiedNick(nick string) {
	// TODO - Locking ?
	go func() {
		time.Sleep(60 * time.Second)

		if ns.mere.curstate.IsNew(nick) {
			newNick := ns.assignRandomNick(nick)
			ns.notice(newNick, "Your nickname has been changed because you did not identify")
		}
	}()
}

func (ns *NickServ) assignRandomNick(nick string) string {
	newNick := "User" + strconv.Itoa(ns.rand.Int())
	ns.c().svsnick(ns.Nick, nick, newNick)
	return newNick
}

func (ns *NickServ) identifyUser(nick string) {
	ns.mere.curstate.Identify(nick)
	ns.c().svsmode(ns.Nick, nick, "+r")
	ns.c().chghost(ns.Nick, nick, "registered."+nick)
	ns.notice(nick, "You are now identified")
}

func (ns *NickServ) unidentifyUser(nick string) {
	ns.mere.curstate.Unidentify(nick)
}

func (ns *NickServ) notice(recip, message string) {
	ns.c().notice(ns.Nick, recip, message)
}

func (ns *NickServ) c() *Connection {
	return ns.mere.connection
}

func (ns *NickServ) r() *RegisteredNickRepo {
	return ns.mere.datastore.RegisteredNicks
}
