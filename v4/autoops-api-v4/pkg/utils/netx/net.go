package netx

import "net"

// GetLocalIP 获取本机IP
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "127.0.0.1"
}

// WrapSshConn 包装SSH连接(暂不实现)
func WrapSshConn(conn net.Conn) net.Conn {
	return conn
}
