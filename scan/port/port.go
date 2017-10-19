package port

import (
	"net"
	"time"
)

const (
	identifier = "port"
)

var (
	Doer = &PortDoer{}
)

type PortDoer struct {
}

func (p *PortDoer) Identifier() string {
	return identifier
}

func (p *PortDoer) Do(addr string) (bool, []string) {
	conn, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		return false, nil
	}
	if conn != nil {
		conn.Close()
	}
	return true, nil
}
