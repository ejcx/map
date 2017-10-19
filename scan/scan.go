package scan

import (
	"fmt"
	"net"
	"strings"
	"sync"
)

const (
	Workers = 20
)

type Scan interface {
	Do(s *NetScan)
	Identifier() string
}

func RegisterDoer(s *Scan) {

}

type Connector interface {
	Do(addr string) bool
}

type IPEnumerator struct {
	iprange *net.IPNet
	current net.IP
}

type NetScan struct {
	CIDR  []*Cidr
	Ports []int
	Scans []string
}

type Cidr struct {
	IP    net.IP
	IPNet *net.IPNet
}

// ParseSetCIDR will parse a CIDR string that is in the format
// described in RFC 4632 and RFC 4291, and return a NetScan
// that contains the set ParseSetCIDR
func (p *NetScan) ParseSetCIDR(cidr string) error {
	i, c, err := net.ParseCIDR(cidr)
	if err != nil {
		return err
	}
	p.CIDR = []*Cidr{
		&Cidr{
			IP:    i,
			IPNet: c,
		},
	}
	return nil
}
func (p *NetScan) SetCIDR(c *Cidr) {
	p.CIDR = []*Cidr{c}
}

func Do(p *NetScan, conFunc []Connector) {
	var wg sync.WaitGroup

	for _, c := range conFunc {
		for _, n := range p.CIDR {
			for _, xport := range p.Ports {
				ipList := ipEnumerator(n)
				q := make(chan net.IP)
				for i := 0; i < Workers; i++ {
					wg.Add(1)
					go func() {
						for ip := range q {
							fmter := "%s:%d"
							// Handle IPv6 Addresses.
							if strings.Contains(ip.String(), ":") {
								fmter = "[%s]:%d"
							}
							addr := fmt.Sprintf(fmter, ip.String(), xport)
							yes := c.Do(addr)
							if yes {
								fmt.Println(addr)
							}
						}
						defer wg.Done()
						return
					}()
				}
				for i := 0; i < len(ipList); i++ {
					q <- ipList[i]
				}
				close(q)
			}
		}
	}
	wg.Wait()
}

func ipEnumerator(c *Cidr) []net.IP {
	var ipList []net.IP
	for ip := c.IP.Mask(c.IPNet.Mask); c.IPNet.Contains(ip); inc(ip) {
		newip := make(net.IP, len(ip))
		copy(newip, ip)
		ipList = append(ipList, newip)
	}
	return ipList
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
