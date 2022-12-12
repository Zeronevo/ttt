package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net"
	"zero/breaker/server"
)

func clientConn(conn net.Conn) {
	defer conn.Close()

	// 连接远程服务器
	addr := fmt.Sprintf("%s:%d", remoteIp, remotePort)
	rconn, err := net.Dial(server.CFG.IpVersion, addr)
	if err != nil {
		log.Println("<ClientConn> connect remote server error: ", err)
		return
	}
	defer rconn.Close()

	// 设置加密传输对象
	key := []byte(secureKey)
	iv := sha256.Sum256(key)

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Println("<ClientConn> NewCipher error: ", err)
		return
	}

	stream := cipher.NewCTR(block, iv[:aes.BlockSize])

	reader := &cipher.StreamReader{S: stream, R: rconn}
	writer := &cipher.StreamWriter{S: stream, W: rconn}

	// 开始交换数据
	go io.Copy(writer, conn)
	io.Copy(conn, reader)
}
