package shadowmere

import (
	"net"
	"fmt"
	"bufio"
	"strings"
)

type handler func(string, []string)

type Connection struct {
	mere *Services

	conn net.Conn
	reader *bufio.Reader

	handlers map[string]handler
	
	name string
	addr string
	pass string
}

func NewConnection(mere *Services, name, addr, pass string) (*Connection, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(conn)

	srv := &Connection{
		mere: mere,

		conn: conn,
		reader: reader,

		name: name,
		addr: addr,
		pass: pass,
	}
	srv.handlers = map[string]handler{
		"PING": srv.handlePing,
		"PRIVMSG": srv.handlePrivmsg,
		"QUIT": srv.handleQuit,
		"NICK": srv.handleNick,
	}

	return srv, nil
}

func (srv *Connection) Start() {
	srv.authenticateUnreal()
	srv.initializeServices()
	srv.listenLoop()
}

func (srv *Connection) authenticateUnreal() {
	// Implements UnrealIRCd-compatible aurhentication
	srv.write(fmt.Sprintf("PASS :%s\r\n", srv.pass))
	srv.write(fmt.Sprintf("PROTOCTL %s\r\n", "SJ3 NICKv2 NOQUIT"))
	srv.write(fmt.Sprintf("SERVER %s 1 :%s\r\n", srv.name, "Services"))
}

func (srv *Connection) initializeServices() {
	ns := srv.mere.nickserv
	srv.nick(ns.Nick, ns.Nick, srv.name, "Services")
}

func (srv *Connection) listenLoop() {
	for {
		line, err := srv.read()
		if err != nil {
			fmt.Print("***ERROR*** ", err.Error())
			return
		}
		
		srv.handleLine(line)
	}
}

func (srv *Connection) handleLine(line string) {
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

func (srv *Connection) handlePing(origin string, args []string) {
	if len(args) == 0 {
		fmt.Println("handlePing(): Malformed PING!")
		return
	}

	srv.pong(args[0])
}

func (srv *Connection) handlePrivmsg(origin string, args []string) {
	if len(args) < 2 {
		fmt.Println("handlePing(): Malformed PRIVMSG!")
		return
	}

	to := args[0]
	msg := args[1]
	if strings.ToLower(srv.mere.nickserv.Nick) == strings.ToLower(to) {
		srv.mere.nickserv.OnPrivmsg(origin, msg)
	}
}

func (srv *Connection) handleQuit(origin string, args []string) {
	var msg string
	if len(args) > 0 {
		msg = args[0]
	}

	srv.mere.nickserv.OnQuit(origin, msg)
}

func (srv *Connection) handleNick(origin string, args []string) {
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

func (srv *Connection) handleNewNick(args []string) {
	if len(args) < 1 {
		fmt.Println("handleNewNick(): Malformed NICKv2!")
		return
	}

	newNick := args[0]
	srv.mere.nickserv.OnNewNick(newNick)
}

func (srv *Connection) handleNickChange(origin string, args []string) {
	if len(args) < 1 {
		fmt.Println("handleNickChange(): Malformed NICK!")
		return
	}

	oldNick := origin
	newNick := args[0]
	srv.mere.nickserv.OnNickChange(oldNick, newNick)
}

func (srv *Connection) read() (string, error) {
	s, err := srv.reader.ReadString('\n')

	fmt.Printf("<-%s", s)
	return s, err
}

func (srv *Connection) write(s string) {
	fmt.Printf("->%s", s)
	fmt.Fprint(srv.conn, s)
}