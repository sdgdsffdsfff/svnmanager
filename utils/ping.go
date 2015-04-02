package utils

import (
	"io"
	"net"
	"time"
)

func Ping(ip string) (bool, error) {
	one := []byte{}
	conn, err := net.DialTimeout("tcp", ip, time.Second*3)
	if err == nil {
		conn.SetReadDeadline(time.Now())
		if _, err := conn.Read(one); err == io.EOF {
			conn.Close()
			conn = nil
			return false, err
		}
		return true, nil
	}
	return false, err
}

func GetLocalIp() (string, error) {
	return "", nil
}
