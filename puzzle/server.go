package puzzle

import (
	"net"
	"fmt"
	"bufio"
)

type handler func(string, []string)

type Server struct {
	conn net.Conn
	reader *bufio.Reader

	handlers map[string]handler
	
	nickserv *NickServ

	// TODO - This really doesn't belong here
	datastore *Datastore

	name string
	addr string
	pass string
}

func NewServer(name, addr, pass string, ds *Datastore) (*Server, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(conn)

	server := &Server{
		conn: conn,
		reader: reader,

		datastore: ds,

		name: name,
		addr: addr,
		pass: pass,
	}

	server.handlers = map[string]handler{
		"PING": server.handlePing,
		"PRIVMSG": server.handlePrivmsg,
		"QUIT": server.handleQuit,
		"NICK": server.handleNickChange,
	}

	return server, nil
}

func (srv *Server) Start() {
	srv.authenticateTS5()
	srv.initializeServices()
	srv.listenLoop()
}

func (srv *Server) authenticateTS5() {
	// Implements authentication that is compliant with the TS5
	// protocol - but NOT compliant with RFC1459/RFC2813.
	srv.write(fmt.Sprintf("PASS %s TS 5\r\n", srv.pass))
	srv.write(fmt.Sprintf("SERVER %s 1 1 :%s\r\n", srv.name, srv.name))
}

func (srv *Server) initializeServices() {
	ns := NewNickserv(srv)
	srv.nick(ns.Nick, ns.Nick, srv.name, "Services")
	srv.nickserv = ns
}

func (srv *Server) listenLoop() {
	for {
		line, err := srv.read()
		if err != nil {
			fmt.Print("***ERROR*** ", err.Error())
			return
		}
		
		srv.handleLine(line)
	}
}

func (srv *Server) handleLine(line string) {
	command, origin, args, err := srv.parseMessage(line)
	if err != nil {
		fmt.Errorf("handleLine(): %s", err.Error())
		return
	}

	h := srv.handlers[command]
	if h != nil {
		h(origin, args)
	}
}

func (srv *Server) handlePing(origin string, args []string) {
	if len(args) == 0 {
		fmt.Errorf("handlePing(): Malformed PING!")
		return
	}

	srv.pong(args[0])
}

func (srv *Server) handlePrivmsg(origin string, args []string) {
	if len(args) < 2 {
		fmt.Errorf("handlePing(): Malformed PRIVMSG!")
		return
	}

	to := args[0]
	msg := args[1]
	if srv.nickserv.Nick == to {
		srv.nickserv.OnPrivmsg(origin, msg)
	}
}

func (srv *Server) handleQuit(origin string, args []string) {
	var msg string
	if len(args) > 0 {
		msg = args[0]
	}

	srv.nickserv.OnQuit(origin, msg)
}

func (srv *Server) handleNickChange(origin string, args []string) {
	if origin == "" {
		// If not origin is set, this is a server-NICK, introducting a
		// new user. We're not interested in these, only the other kind
		// of NICK, which is a nickchange.
		return
	}
	if len(args) < 1 {
		fmt.Errorf("handleNickChange(): Malformed NICK!")
		return
	}

	oldNick := origin
	newNick := args[0]
	srv.nickserv.OnNickChange(oldNick, newNick)
}

func (srv *Server) read() (string, error) {
	s, err := srv.reader.ReadString('\n')

	fmt.Printf("<-%s", s)
	return s, err
}

func (srv *Server) write(s string) {
	fmt.Printf("->%s", s)
	fmt.Fprint(srv.conn, s)
}