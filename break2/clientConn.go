package main

import (
	"fmt"
	"log"
	"net"
	"zero/break/server"
	"zero/break/transfers"
)

func clientConn(conn net.Conn) {
	defer conn.Close()

	addr := fmt.Sprintf("%s:%d", server.CFG.RemoteIp, server.CFG.RemotePort)
	rconn, err := net.Dial(server.CFG.Ipv, addr)
	if err != nil {
		log.Println("<ClientConn> can not connect remote server error: ", err)
		return
	}
	defer rconn.Close()

	trans, err := transfers.NewTransfers(server.CFG.SecureKey)
	if err != nil {
		log.Println("<ClientConn> create transfers error: ", err)
		return
	}

	go trans.DataCopyDe(conn, rconn)
	trans.DataCopyEn(rconn, conn)
}
