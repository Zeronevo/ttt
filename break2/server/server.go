package server

import (
	"fmt"
	"log"
	"net"
)

type handleConn func(net.Conn)

type Config struct {
	Name       string
	Ipv        string
	Ip         string
	Port       int
	RemoteIp   string
	RemotePort int
	SecureKey  string
}

var CFG Config

type Server struct {
	name   string
	ipv    string
	ip     string
	port   int
	handle handleConn
}

func NewServer() Server {
	return Server{
		name:   CFG.Name,
		ipv:    CFG.Ipv,
		ip:     CFG.Ip,
		port:   CFG.Port,
		handle: nil,
	}
}

func (s *Server) Run() {
	s.start()
}

func (s *Server) start() {
	fmt.Printf("[Server] %s is starting...\n", s.name)
	addr := fmt.Sprintf("%s:%d", s.ip, s.port)

	server, err := net.Listen(s.ipv, addr)
	if err != nil {
		fmt.Println(err)

	}
	defer server.Close()
	fmt.Printf("[Server] %s start on -- %s:%d\n", s.name, s.ip, s.port)

	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		if s.handle == nil {
			log.Println("<Server> Not register handle function.")
			break
		}

		go s.handle(conn)
	}
}

func (s *Server) Handle(fun handleConn) {
	s.handle = fun
}
