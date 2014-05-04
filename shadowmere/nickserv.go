package shadowmere

import (
	"fmt"
	"strings"
	"time"
	"math/rand"
	"strconv"
)

type nickservCmdHandler func(string, []string)

type NickServ struct {
	Nick string
	server *Server

	handlers map[string]nickservCmdHandler

	repo *RegisteredNickRepo
	rand *rand.Rand
}

func NewNickserv(server *Server) *NickServ {
	ns := &NickServ{
		Nick: "NickServ",
		server: server,
		repo: server.datastore.RegisteredNicks,
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
	ns.server.svsmode(ns.Nick, newNick, "-r")

	ns.handleRegisteredNick(newNick)
}

func (ns *NickServ) handleRegister(nick string, args []string) {
	if len(args) != 2 {
		ns.notice(nick, "Syntax: REIGSTER <email> <password>")
		return
	}

	rn, err := ns.repo.GetByNick(nick)
	if err != nil {
		// TODO - Handle error better
		fmt.Printf("***ERROR*** %v\n", err.Error())
		return
	}
	if rn != nil {
		ns.notice(nick, "That nickname is already registered.")
		return
	}

	newRn := &RegisteredNick{
		Nick: nick,
		Email: args[0],
		Passwd: args[1],
	}
	err = ns.repo.Register(newRn)
	if err != nil {
		ns.notice(nick, "Error registering nickname")
		fmt.Printf("***ERROR*** %s\n", err.Error())
		return
	}

	ns.notice(nick, "Your nickname is now registered!")
	ns.identifyUser(nick)
}

func (ns *NickServ) handleIdentify(nick string, args []string) {
	if len(args) != 1{
		ns.notice(nick, "You must specify a password")
		return
	}
	if ns.server.curstate.IsIdentified(nick) {
		ns.notice(nick, "You are already identified")
		return
	}

	rn, err := ns.repo.GetByNick(nick)
	if err != nil {
		// TODO - Handle error better
		fmt.Printf("***ERROR*** %v\n", err.Error())
		return
	}
	if rn == nil {
		ns.notice(nick, "This nickname is not registered")
		return
	}
	
	validPassword := ns.repo.Authenticate(rn, args[0])
	if validPassword {
		ns.identifyUser(nick)
	} else {
		ns.notice(nick, "Invalid password for this nick")
	}	
}

func (ns *NickServ) handleRegisteredNick(nick string) {
	rn, err := ns.repo.GetByNick(nick)
	if err != nil {
		// TODO - Handle error better
		fmt.Printf("***ERROR*** %v\n", err.Error())
		return
	}
	if rn == nil {
		return
	}

	ns.notice(nick, "Your nickname is registered. You have 60 seconds to identify or change your nick.")

	ns.server.curstate.NewNick(nick)
	ns.enforceIdentifiedNick(nick)
}

func (ns *NickServ) enforceIdentifiedNick(nick string) {
	// TODO - Locking ?
	go func() {
		time.Sleep(60 * time.Second)

		if ns.server.curstate.IsNew(nick) {
			newNick := ns.svsnick(nick)
			ns.notice(newNick, "Your nickname has been changed because you did not identify")
		}
	}()
}

func (ns *NickServ) svsnick(nick string) string {
	newNick := "User" + strconv.Itoa(ns.rand.Int())
	ns.server.svsnick(ns.Nick, nick, newNick)
	return newNick
}

func (ns *NickServ) identifyUser(nick string) {
	ns.server.curstate.Identify(nick)
	ns.server.svsmode(ns.Nick, nick, "+r")
	ns.server.chghost(ns.Nick, nick, "registered." + nick)
	ns.notice(nick, "You are now identified")
}

func (ns *NickServ) unidentifyUser(nick string) {
	ns.server.curstate.Unidentify(nick)
}

func (ns *NickServ) notice(recip, message string) {
	ns.server.notice(ns.Nick, recip, message)
}
