package puzzle

import (
	"fmt"
)

type NickServ struct {
	Service
	server *Server
}

func NewNickserv(server *Server) *NickServ {
	return &NickServ{
		server: server,
	}
}

func (ns *NickServ) Nick() string {
	return "NickServ"
}

func (ns *NickServ) OnPrivmsg(nick, content string) {
	fmt.Printf("NickServ <%s> %s\n", nick, content)
}