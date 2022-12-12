package main

import (
	"flag"
	"log"
	"zero/break/server"
)

const (
	defaultServerName = "Break server"
	defaultClientName = "Breaker"
	defaultListenIp   = "0.0.0.0"
	defaultIpv        = "tcp"
	defaultServerPort = 9999
	defaultClientPort = 6666
	defaultKey        = "1w2q3r4e5y6t7i8u"
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
	flag.IntVar(&serverPort, "p", defaultClientPort, "set server port")
	flag.StringVar(&remoteIp, "remote", defaultListenIp, "connect remote server ip")
	flag.IntVar(&remotePort, "port", defaultServerPort, "connect remote server port")
	flag.StringVar(&secureKey, "key", defaultKey, "change secure key 16 24 32 bit")
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
	server.CFG.Name = defaultClientName
	server.CFG.Ipv = defaultIpv
	server.CFG.Ip = defaultListenIp
	server.CFG.Port = serverPort
	server.CFG.RemoteIp = remoteIp
	server.CFG.RemotePort = remotePort
	server.CFG.SecureKey = secureKey

	// 初始化connection回调函数
	connection := clientConn

	// 判断启动的是客户端还是服务端
	if isServer {
		server.CFG.Name = defaultServerName
		if serverPort == defaultClientPort {
			server.CFG.Port = defaultServerPort
		}
		connection = serverConn
	}

	// 创建服务器
	s := server.NewServer()
	s.Handle(connection)
	s.Run()
}
