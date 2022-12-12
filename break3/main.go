package main

import (
	"flag"
	"log"
	"zero/breaker/server"
)

var ( // 需要从命令行获取的变量
	isServer   bool
	serverPort int
	remoteIp   string
	remotePort int
	secureKey  string
)

func init() { // 解析命令行参数
	flag.BoolVar(&isServer, "s", false, "start server")
	flag.IntVar(&serverPort, "p", server.DefaultListenPort, "set server port")
	flag.StringVar(&remoteIp, "remote", server.DefaultRemoteIp, "connect remote server ip")
	flag.IntVar(&remotePort, "port", server.DefaultRemotePort, "connect remote server port")
	flag.StringVar(&secureKey, "key", server.DefaultSecureKey, "change secure key 16 24 32 bit")
}
func main() {
	flag.Parse()

	// 判断密钥字节正确
	keylen := len([]byte(secureKey))
	if keylen != 16 && keylen != 24 && keylen != 32 {
		log.Printf("[Config] key %d bit, length of key must be 16 24 32 bit", keylen)
		return
	}

	// 初始化配置
	server.CFG.IpVersion = server.DefaultIpVersion
	server.CFG.ListenIp = server.DefaultListenIp
	server.CFG.SecureKey = secureKey

	// 初始化配置客户端
	server.CFG.ServerName = "Break"
	server.CFG.ListenPort = serverPort
	server.CFG.RemoteIp = remoteIp
	server.CFG.RemotePort = remotePort
	connection := clientConn

	// 初始化配置服务端
	if isServer {
		server.CFG.ServerName = "Break Server"
		server.CFG.ListenPort = remotePort
		connection = serverConn
	}

	// 创建并启动服务器
	s := server.NewServer()
	s.Handle(connection)
	s.Run()
}
