package main

import (
	"encoding/gob"
	"log"
	"net"
	"zero/break/server"
	"zero/break/socks"
	"zero/break/transfers"
)

func serverConn(conn net.Conn) {
	defer conn.Close()

	// init
	var buf []byte
	enc := gob.NewEncoder(conn) // Will write to network.
	dec := gob.NewDecoder(conn) // Will read from network.
	sock := socks.NewSocks()
	trans, err := transfers.NewTransfers(server.CFG.SecureKey)
	if err != nil {
		log.Println("<ServerConn> create transfers error: ", err)
		return
	}

	// read crypto data from client request
	err = dec.Decode(&buf)
	if err != nil {
		log.Println("<ServerConn> read date for auth error: ", err)
		return
	}

	// decrypto data
	obuf, err := trans.DeData(buf)
	if err != nil {
		log.Println("<ServerConn> Decode date for auth error: ", err)
		return
	}

	// begin auth
	sock.Req = obuf
	err = sock.Auth()
	if err != nil {
		log.Println("<ServerConn> auth error: ", err)
		buf, _ = trans.EnData(sock.Resp)
		enc.Encode(buf)
		return
	}

	// auth right and then response client
	buf, err = trans.EnData(sock.Resp)
	if err != nil {
		log.Println("<ServerConn> Eecode date for response auth error: ", err)
		return
	}

	err = enc.Encode(buf)
	if err != nil {
		log.Println("<ServerConn> response auth data error: ", err)
		return
	}

	// read data 2 request
	err = dec.Decode(&buf)
	if err != nil {
		log.Println("<ServerConn> read date for request error: ", err)
		return
	}

	// decrypto data
	obuf, err = trans.DeData(buf)
	if err != nil {
		log.Println("<ServerConn> Decode date for request error: ", err)
		return
	}

	// begin request
	sock.Req = obuf
	err = sock.Request()
	if err != nil {
		log.Println("<ServerConn> request error: ", err)
		buf, _ = trans.EnData(sock.Resp)
		enc.Encode(buf)
		return
	}

	// response client
	buf, err = trans.EnData(sock.Resp)
	if err != nil {
		log.Println("<ServerConn> Eecode date for request error: ", err)
		return
	}

	err = enc.Encode(buf)
	if err != nil {
		log.Println("<ServerConn> response request data error: ", err)
		return
	}

	// request ok get address port and dial
	addr, err := sock.GetAddr()
	if err != nil {
		log.Println("<ServerConn> can not get target error: ", err)
	}

	tConn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println("<ServerConn> connect target error: ", err)
		return
	}
	defer tConn.Close()

	// socks5 done begin transfers data
	go trans.DataCopyEn(conn, tConn)
	trans.DataCopyDe(tConn, conn)
}
