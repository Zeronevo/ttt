package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"io"
	"log"
	"net"
	"zero/breaker/socks"
)

func serverConn(conn net.Conn) {
	defer conn.Close()

	// 设置加密传输对象
	key := []byte(secureKey)
	iv := sha256.Sum256(key)

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Println("<ClientConn> NewCipher error: ", err)
		return
	}

	stream := cipher.NewCTR(block, iv[:aes.BlockSize])

	reader := &cipher.StreamReader{S: stream, R: conn}
	writer := &cipher.StreamWriter{S: stream, W: conn}

	// 初始化socks5
	sock := socks.NewSocks()
	buf := make([]byte, 1024)

	// begin auth
	n, err := reader.Read(buf)
	if err != nil {
		log.Println("<ServerConn> read date for auth error: ", err)
		return
	}

	sock.Req = buf[:n]
	err = sock.Auth()
	if err != nil {
		log.Println("<ServerConn> auth error: ", err)
		writer.Write(sock.Resp)
		return
	}
	writer.Write(sock.Resp)

	// begin response
	n, err = reader.Read(buf)
	if err != nil {
		log.Println("<ServerConn> read date for request error: ", err)
		return
	}

	sock.Req = buf[:n]
	err = sock.Request()
	if err != nil {
		log.Println("<ServerConn> request error: ", err)
		writer.Write(sock.Resp)
		return
	}
	writer.Write(sock.Resp)

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

	// 开始交换数据
	go io.Copy(writer, tConn)
	io.Copy(tConn, reader)
}
