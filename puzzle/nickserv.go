package puzzle

import (
	"fmt"
	"strings"
	"time"
	"math/rand"
	"strconv"
	"sync"
)

type nickservCmdHandler func(string, []string)

type NickServ struct {
	Nick string
	server *Server
	handlers map[string]nickservCmdHandler

	rand *rand.Rand

	// Store a list, by nick, of users that have identified to NickServ.
	// In a production system, we would probably keep track of all users
	// on the network and whether they are identified. However, in this
	// case, simply tracking identified users is sufficient.
	identified map[string]bool

	// Locking for nicks that are awaiting being killed
	killLocks map[string]bool
	killLocksLock sync.Mutex
}

func NewNickserv(server *Server) *NickServ {
	ns := &NickServ{
		Nick: "NickServ",
		server: server,
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
		identified: make(map[string]bool),
		killLocks: make(map[string]bool),
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
		ns.notice(nick, fmt.Sprintf("No such command: %s", command))
	}
}

func (ns *NickServ) OnQuit(nick, quitMessage string) {
	ns.userUnIdentified(nick)
}

func (ns *NickServ) OnNewNick(nick string) {
	ns.handleRegisteredNick(nick)
}

func (ns *NickServ) OnNickChange(oldNick, newNick string) {
	ns.userUnIdentified(oldNick)
	ns.server.svsmode(ns.Nick, newNick, "-r")

	ns.handleRegisteredNick(newNick)
}

func (ns *NickServ) handleRegister(nick string, args []string) {
	if len(args) != 2 {
		ns.notice(nick, "Syntax: REIGSTER <email> <password>")
		return
	}

	rn, err := ns.server.datastore.GetRegisteredNick(nick)
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
	err = ns.server.datastore.Register(newRn)
	if err != nil {
		ns.notice(nick, "Error registering nickname")
		fmt.Printf("***ERROR*** %s\n", err.Error())
		return
	}

	ns.notice(nick, "Your nickname is now registered!")
	ns.userIdentified(nick)
}

func (ns *NickServ) handleIdentify(nick string, args []string) {
	if len(args) != 1{
		ns.notice(nick, "You must specify a password")
		return
	}
	if ns.isIdentified(nick) {
		ns.notice(nick, "You are already identified")
		return
	}

	rn, err := ns.server.datastore.GetRegisteredNick(nick)
	if err != nil {
		// TODO - Handle error better
		fmt.Printf("***ERROR*** %v\n", err.Error())
		return
	}
	if rn == nil {
		ns.notice(nick, "This nickname is not registered")
		return
	}
	
	validPassword := ns.server.datastore.Authenticate(rn, args[0])
	if validPassword {
		ns.userIdentified(nick)
	} else {
		ns.notice(nick, "Invalid password for this nick")
	}	
}

func (ns *NickServ) handleRegisteredNick(nick string) {
	rn, err := ns.server.datastore.GetRegisteredNick(nick)
	if err != nil {
		// TODO - Handle error better
		fmt.Printf("***ERROR*** %v\n", err.Error())
		return
	}
	if rn == nil {
		return
	}

	ns.notice(nick, "Your nickname is registered. You have 60 seconds to identify or change your nick.")
	ns.enforceIdentifiedNick(nick)
}

func (ns *NickServ) enforceIdentifiedNick(nick string) {
	ns.killLocksLock.Lock()
	defer ns.killLocksLock.Unlock()
	if ns.killLocks[nick] == true {
		// This nick already has a goroutine spawned waiting to svsnick
		// it, so don't spawn a new one.
		return
	}
	ns.killLocks[nick] = true

	go func() {
		time.Sleep(60 * time.Second)

		ns.killLocksLock.Lock()
		defer ns.killLocksLock.Unlock()

		if ns.identified[nick] != true {
			newNick := ns.svsnick(nick)
			ns.notice(newNick, "Your nickname has been changed because you did not identify")
		}

		delete(ns.killLocks, nick)
	}()
}

func (ns *NickServ) svsnick(nick string) string {
	newNick := "User" + strconv.Itoa(ns.rand.Int())
	ns.server.svsnick(ns.Nick, nick, newNick)
	return newNick
}

func (ns *NickServ) userIdentified(nick string) {
	ns.identified[nick] = true
	ns.server.svsmode(ns.Nick, nick, "+r")
	ns.server.chghost(ns.Nick, nick, "registered." + nick)
	ns.notice(nick, "You are now identified")
}

func (ns *NickServ) userUnIdentified(nick string) {
	delete(ns.identified, nick)
}

func (ns *NickServ) isIdentified(nick string) bool {
	return ns.identified[nick] == true
}

func (ns *NickServ) notice(recip, message string) {
	ns.server.notice(ns.Nick, recip, message)
}
