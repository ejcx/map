package port

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

const (
	Workers = 20
)

type Report struct {
}

type IPEnumerator struct {
	iprange *net.IPNet
	current net.IP
}

type Cidr struct {
	IP    net.IP
	IPNet *net.IPNet
}

type PortScan struct {
	CIDR []*Cidr

	Ports []int
}

func NewPortScan(c *Cidr) *PortScan {
	cidr := []*Cidr{c}
	return &PortScan{
		CIDR: cidr,
	}
}

func New() *PortScan {
	return &PortScan{}
}

func (p *PortScan) SetCIDR(c *Cidr) {
	p.CIDR = []*Cidr{c}
}

// ParseSetCIDR will parse a CIDR string that is in the format
// described in RFC 4632 and RFC 4291, and return a PortScan
// that contains the set ParseSetCIDR
func (p *PortScan) ParseSetCIDR(cidr string) error {
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

func (p *PortScan) Do() {
	fmt.Println(p)
	var wg sync.WaitGroup

	for _, n := range p.CIDR {
		for _, port := range p.Ports {
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
						conn, err := net.DialTimeout("tcp", fmt.Sprintf(fmter, ip.String(), port), time.Second)
						if err != nil {
							fmt.Println(err)
							continue
						}
						fmt.Println("Actual conneciton!")
						conn.Close()
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
