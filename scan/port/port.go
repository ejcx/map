package port

import (
	"net"
	"time"
)

var (
	Doer = &PortDoer{}
)

type PortDoer struct {
}

func (p *PortDoer) Do(addr string) bool {
	conn, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		return false
	}
	if conn != nil {
		conn.Close()
	}
	return true
}
