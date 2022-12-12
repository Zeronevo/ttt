// socks 实现了socks5的基本验证
// 使用方法：
// 创建Socks对象将请求放入Req进行验证之后可以在Resp得到回复内容
// Auth()方法用来确认客户端的版本以及可用方法
// Request()方法用来验证获取客户端要访问的目标
// GetAddr()获取最后的ip：port字符串
package socks

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

// 通用常量
// VER    protocol version: X'05'
// RSV    RESERVED
const VER = 0x05
const RSV = 0x00
const ZERO = 0x00

// 目前已定义方法
// X'00' NO AUTHENTICATION REQUIRED
// X'01' GSSAPI
// X'02' USERNAME/PASSWORD
// X'03' to X'7F' IANA ASSIGNED
// X'80' to X'FE' RESERVED FOR PRIVATE METHODS
// X'FF' NO ACCEPTABLE METHODS
const (
	NOAUTH   = 0x00
	GSSAPI   = 0x01
	NAMEPASS = 0x02
	NOACCEPT = 0xFF
)

// 命令字段
// CONNECT X'01'
// BIND X'02'
// UDP ASSOCIATE X'03'
const (
	CONNECT = 0x01
	BIND    = 0x02
	UDP     = 0x03
)

// 地址类型
// IP V4 address: X'01'
// DOMAINNAME: X'03'
// IP V6 address: X'04'
const (
	IPV4   = 0x01
	DOMAIN = 0x03
	IPV6   = 0x04
)

// 回复客户端
// X'00' succeeded 成功
// X'01' general SOCKS server failure 常规 SOCKS 服务故障
// X'02' connection not allowed by ruleset 规则不允许的连接
// X'03' Network unreachabl 网络不可达
// X'04' Host unreachable 主机无法访问
// X'05' Connection refused 拒绝连接
// X'06' TTL expired 连接超时
// X'07' Command not supported 不支持的命令
// X'08' Address type not supported 不支持的地址类型
// X'09' to X'FF' unassigned 未定义
const (
	SUCCEEDED = 0x00
	CMDNOSUP  = 0x07
)

// 定义socks对象
type Socks struct {
	Req  []byte
	Resp []byte
	Addr string
	Port int
}

func NewSocks() Socks {
	return Socks{
		Req:  nil,
		Resp: nil,
		Addr: "",
		Port: 0,
	}
}

// auth read request and send respose
// auth client message
// +----+----------+----------+
// |VER | NMETHODS | METHODS  |
// +----+----------+----------+
// | 1  |    1     | 1 to 255 |
// +----+----------+----------+
func (s *Socks) Auth() error {
	b := bytes.NewBuffer(s.Req)

	// validate version
	ver := int(b.Next(1)[0])
	if ver != VER {
		s.Resp = []byte{VER, NOACCEPT}
		return errors.New("[socks] protocol version not supported")
	}

	// get methods
	nMethods := int(b.Next(1)[0])
	methods := b.Next(nMethods)

	// 检查版本方法并回复客户端
	// +----+--------+
	// |VER | METHOD |
	// +----+--------+
	// | 1  |   1    |
	// +----+--------+
	acceptable := false
	// 目前只支持无认证方法
	for _, v := range methods {
		if v == NOAUTH {
			// if support then respose
			s.Resp = []byte{VER, NOAUTH}
			acceptable = true
		}
	}

	// if not support tell client
	if !acceptable {
		s.Resp = []byte{VER, NOACCEPT}
		return errors.New("[socks] protocol method not supported")
	}

	// success auth
	return nil
}

func (s *Socks) Request() error {
	// get request
	b := bytes.NewBuffer(s.Req)

	//  check version command reserved address
	// +----+-----+-------+------+----------+----------+
	// |VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
	// +----+-----+-------+------+----------+----------+
	// | 1  |  1  | X'00' |  1   | Variable |    2     |
	// +----+-----+-------+------+----------+----------+
	version := int(b.Next(1)[0])
	command := int(b.Next(1)[0])
	reserved := int(b.Next(1)[0])
	addrType := int(b.Next(1)[0])

	if version != VER {
		s.Resp = []byte{VER, NOACCEPT}
		return errors.New("[socks] protocol version not supported")
	}

	if command != CONNECT && command != BIND && command != UDP {
		s.Resp = []byte{VER, NOACCEPT}
		return errors.New("[socks] protocol command not supported")
	}

	if reserved != RSV {
		s.Resp = []byte{VER, NOACCEPT}
		return errors.New("[socks] protocol reserved error")
	}

	if addrType != IPV4 && addrType != DOMAIN && addrType != IPV6 {
		s.Resp = []byte{VER, NOACCEPT}
		return errors.New("[socks] protocol address type not supported")
	}
	// 获取地址并回复客户端
	// +----+-----+-------+------+----------+----------+
	// |VER | REP |  RSV  | ATYP | BND.ADDR | BND.PORT |
	// +----+-----+-------+------+----------+----------+
	// | 1  |  1  | X'00' |  1   | Variable |    2     |
	// +----+-----+-------+------+----------+----------+
	// read address
	var address string

	switch addrType {
	case IPV6: //IPv6
		s.Resp = []byte{VER, NOACCEPT}
		return errors.New("[socks] protocol iPv6 no supportd yet")
	case IPV4: // IPv4
		address = fmt.Sprintf("%d.%d.%d.%d", b.Next(1)[0], b.Next(1)[0], b.Next(1)[0], b.Next(1)[0])
	case DOMAIN: // DomainName
		addrLen := int(b.Next(1)[0])
		address = string(b.Next(addrLen))
	default:
		s.Resp = []byte{VER, NOACCEPT}
		return errors.New("[socks] protocol invalid address type")
	}

	// read port
	var port uint16
	port = binary.BigEndian.Uint16(b.Next(2))

	// error write responts to client status
	// 目前只支持connect 不支持 bind 和 UDP
	if command != CONNECT { // no support udp and bind 0x05, 0x07, 0x00, 0x01, 0, 0, 0, 0, 0, 0
		s.Resp = []byte{VER, CMDNOSUP, RSV, IPV4, ZERO, ZERO, ZERO, ZERO, ZERO, ZERO}
		return errors.New("[socks] protocol not support udp and bind")
	}

	// connect out network
	s.Addr = address
	s.Port = int(port)

	// success write responts to client status
	// 0x05, 0x00, 0x00, 0x01, 0, 0, 0, 0, 0, 0
	s.Resp = []byte{VER, SUCCEEDED, RSV, IPV4, ZERO, ZERO, ZERO, ZERO, ZERO, ZERO}

	return nil
}

func (s *Socks) GetAddr() (string, error) {
	if s.Addr == "" || s.Port == 0 {
		s.Resp = []byte{VER, NOACCEPT}
		return "", errors.New("[socks] protocol not get address")
	}
	return fmt.Sprintf("%s:%d", s.Addr, s.Port), nil
}
