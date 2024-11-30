package http

import (
	"encoding/binary"
	"fmt"
	"net"
)

func Ip2long(ipStr string) (uint32, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return 0, fmt.Errorf("invalid IP address: %s", ipStr)
	}

	ip = ip.To4()
	if ip == nil {
		return 0, fmt.Errorf("not an IPv4 address: %s", ipStr)
	}

	return binary.BigEndian.Uint32(ip), nil
}

func Long2ip(ipLong uint32) string {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, ipLong)
	return ip.String()
}
