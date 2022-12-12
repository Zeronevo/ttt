package server

import (
	"fmt"
	"log"
	"net"
)

const (
	DefaultServerName = "Servers"
	DefaultIpVersion  = "tcp"
	DefaultListenIp   = "0.0.0.0"
	DefaultListenPort = 6666
)

type config struct {
	ServerName string
	IpVersion  string
	ListenIp   string
	ListenPort int
}

var cfg = config{
	ServerName: DefaultServerName,
	IpVersion:  DefaultIpVersion,
	ListenIp:   DefaultListenIp,
	ListenPort: DefaultListenPort,
}

type handleConn func(net.Conn)

type server struct {
	name   string
	ipv    string
	ip     string
	port   int
	handle handleConn
}

func NewServer() *server {
	return &server{
		name:   CFG.ServerName,
		ipv:    CFG.IpVersion,
		ip:     CFG.ListenIp,
		port:   CFG.ListenPort,
		handle: nil,
	}
}

func (s *server) Run() {
	if s.handle == nil {
		s.handle = handleConnection
	}

	err := s.start()
	if err != nil {
		log.Printf("<server> Con not listen: %v", err)
		return
	}
}

func (s *server) start() error {
	fmt.Printf("[Server] %s --> starting <--\n", s.name)
	addr := fmt.Sprintf("%s:%d", s.ip, s.port)

	server, err := net.Listen(s.ipv, addr)
	if err != nil {
		return err
	}
	defer server.Close()
	fmt.Printf("[Server] %s start on -- %s:%d\n", s.name, s.ip, s.port)

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Printf("<server> Accept error: %v\n", err)
			continue
		}

		go s.handle(conn)
	}
}

func (s *server) Handle(fun handleConn) {
	s.handle = fun
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("The connection is ok.")
	fmt.Println("But you have not define handle.")
}
