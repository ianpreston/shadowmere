package shadowmere

import (
	"net"
	"fmt"
	"bufio"
	"strings"
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
		"NICK": server.handleNick,
	}

	return server, nil
}

func (srv *Server) Start() {
	srv.authenticateUnreal()
	srv.initializeServices()
	srv.listenLoop()
}

func (srv *Server) authenticateUnreal() {
	// Implements UnrealIRCd-compatible aurhentication
	srv.write(fmt.Sprintf("PASS :%s\r\n", srv.pass))
	srv.write(fmt.Sprintf("PROTOCTL %s\r\n", "SJ3 NICKv2 NOQUIT"))
	srv.write(fmt.Sprintf("SERVER %s 1 :%s\r\n", srv.name, "Services"))
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
		fmt.Println("handleLine(): %s", err.Error())
		return
	}

	h := srv.handlers[command]
	if h != nil {
		h(origin, args)
	}
}

func (srv *Server) handlePing(origin string, args []string) {
	if len(args) == 0 {
		fmt.Println("handlePing(): Malformed PING!")
		return
	}

	srv.pong(args[0])
}

func (srv *Server) handlePrivmsg(origin string, args []string) {
	if len(args) < 2 {
		fmt.Println("handlePing(): Malformed PRIVMSG!")
		return
	}

	to := args[0]
	msg := args[1]
	if strings.ToLower(srv.nickserv.Nick) == strings.ToLower(to) {
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

func (srv *Server) handleNick(origin string, args []string) {
	if origin == "" {
		// Server introducing a new user
		// nick hopcount timestamp	username hostname server servicestamp +usermodes virtualhost :realname
		srv.handleNewNick(args)
	} else {
		// User changing their nick
		// :old nick new timestamp
		srv.handleNickChange(origin, args)
	}
}

func (srv *Server) handleNewNick(args []string) {
	if len(args) < 1 {
		fmt.Println("handleNewNick(): Malformed NICKv2!")
		return
	}

	newNick := args[0]
	srv.nickserv.OnNewNick(newNick)
}

func (srv *Server) handleNickChange(origin string, args []string) {
	if len(args) < 1 {
		fmt.Println("handleNickChange(): Malformed NICK!")
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