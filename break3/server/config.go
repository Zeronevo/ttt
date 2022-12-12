package server

const (
	DefaultRemoteIp   = "0.0.0.0"
	DefaultRemotePort = 9999
	DefaultSecureKey  = "1w2q3r4e5y6t7i8u"
)

type Config struct {
	config
	RemoteIp   string
	RemotePort int
	SecureKey  string
}

var CFG = Config{
	config:     cfg,
	RemoteIp:   DefaultRemoteIp,
	RemotePort: DefaultListenPort,
	SecureKey:  DefaultSecureKey,
}
