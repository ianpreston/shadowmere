package puzzle

import (
	"fmt"
	"time"
)

type Service interface {
	Nick() string
	OnPrivmsg(string, string)
}

func IdentifyService(sv Service, server *Server) {
	// TODO - There has to be a better way to do this
	// TODO - Perhaps move this method into Server?
	// NICK <nick> <hops> <ts> <modes> <user> <host> <server> :<real>
	line := fmt.Sprintf(
		"NICK %s 1 %v +i %s services %s :services\r\n",
		sv.Nick(),
		time.Now().Unix(),
		sv.Nick(),
		server.name,
	)

	server.write(line)
}