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
	services []Service

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
	}

	server.services = []Service{
		NewNickserv(server),
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
	for _, sv := range srv.services {
		IdentifyService(sv, srv)
	}
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
	command, origin, args, err := parseMessage(line)
	if err != nil {
		fmt.Errorf("handleLine(): %s", err.Error())
	}

	h := srv.handlers[command]
	if h != nil {
		h(origin, args)
	}
}

func (srv *Server) handlePing(origin string, args []string) {
	if len(args) == 0 {
		fmt.Errorf("handlePing(): Malformed PING!")
	}

	srv.write(fmt.Sprintf("PONG :%s\r\n", args[0]))
}

func (srv *Server) handlePrivmsg(origin string, args []string) {
	if len(args) < 2 {
		fmt.Errorf("handlePing(): Malformed PRIVMSG!")
	}

	to := args[0]
	msg := args[1]
	
	for _, sv := range srv.services {
		if sv.Nick() == to {
			sv.OnPrivmsg(origin, msg)
		}
	}
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