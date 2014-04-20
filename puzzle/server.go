package puzzle

import (
	"net"
	"fmt"
	"bufio"
	"strings"
)

type Server struct {
	conn net.Conn
	reader *bufio.Reader
	handlers map[string]func([]string)

	name string
	addr string
	pass string
}

func NewServer() (*Server, error) {
	// TODO - Load configuration values from a file
	name := "noveria.0x-1.com"
	addr := "localhost:6667"
	pass := "foo"

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(conn)

	server := &Server{
		conn: conn,
		reader: reader,

		name: name,
		addr: addr,
		pass: pass,
	}

	server.handlers = map[string]func([]string){
		"PING": server.handlePing,
	}

	return server, nil
}

func (srv *Server) Start() {
	srv.authenticateTS5()
	srv.listenLoop()
}

func (srv *Server) authenticateTS5() {
	// Implements authentication that is compliant with the TS5
	// protocol - but NOT compliant with RFC1459/RFC2813.
	srv.write(fmt.Sprintf("PASS %s TS 5\r\n", srv.pass))
	srv.write(fmt.Sprintf("SERVER %s 1 1 :%s\r\n", srv.name, srv.name))
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
	line = strings.Replace(line, "\r", "", -1)
	line = strings.Replace(line, "\n", "", -1)

	command, args := srv.parseServerLine(line)

	handler := srv.handlers[command]
	if handler != nil {
		handler(args)
	}
}

func (srv *Server) handlePing(args []string) {
	// TODO - It is possible a server could send a "PING" with
	// no server argument.
	srv.write(fmt.Sprintf("PONG %s\r\n", args[0]))
}

func (srv *Server) parseServerLine(line string) (string, []string) {
	// TODO - Handle colon-seperated multi-space last arguments
	tokens := strings.Split(line, " ")

	var command string
	var args []string

	command = strings.ToUpper(tokens[0])
	if len(tokens) > 1 {
		args = tokens[1:]
	}
	
	return command, args
}

func (srv *Server) read() (string, error) {
	s, err := srv.reader.ReadString('\n')

	fmt.Print("<-", s)
	return s, err
}

func (srv *Server) write(s string) {
	fmt.Print("->", s)
	fmt.Fprint(srv.conn, s)
}