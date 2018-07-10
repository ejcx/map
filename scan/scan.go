package scan

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

const (
	Workers = 20
	TAG     = "map"
)

var (
	entryList []Entry
)

type Scan interface {
	Do(s *NetScan)
}

type Entry struct {
	Addr    string
	Success bool
}

type NetScanReport struct {
	Workers int
	Cidrs   []string
}

type ScanInfo struct {
	Info   *NetScanReport
	Start  time.Time
	Finish time.Time
}

type Report struct {
	Entries []Entry
	Scan    ScanInfo
}

type Connector interface {
	Do(addr string) (bool, []string)
	Identifier() string
}

type IPEnumerator struct {
	iprange *net.IPNet
	current net.IP
}

type NetScan struct {
	CIDR    []*Cidr
	Ports   []int
	Scans   []string
	Workers int
	Verbose bool
}

type Cidr struct {
	IP    net.IP
	IPNet *net.IPNet

	Original string
}

func (c *Cidr) String() string {
	return c.Original
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

func (p *NetScan) Report() *NetScanReport {
	var cidrs []string
	for _, c := range p.CIDR {
		cidrs = append(cidrs, c.Original)
	}
	return &NetScanReport{
		Cidrs:   cidrs,
		Workers: p.Workers,
	}
}
func Do(p *NetScan, conFunc []Connector) {

	var (
		// Make sure we wait for all goroutines to gracefully exit.
		// before reporting.
		wg sync.WaitGroup

		// Net:Port format. In IPv6 make sure you add brackets.
		addrFormat = "%s:%d"

		// Keep track of the ones we scan.
		scannedList = make(map[string]bool)
	)

	// For now, there is only one conFunc that should
	// be passed to this at a time. Instead of having
	// a loop and this code being future proof, just
	// take out the conFunc and run the one.
	if len(conFunc) != 1 {
		log.Fatalf("Invalid number of arguments passed to Do: %d", len(conFunc))
	}
	c := conFunc[0]

	// Let's get the party started. Start the clock.
	start := time.Now()
	for _, n := range p.CIDR {
		for _, xport := range p.Ports {
			ipList := ipEnumerator(n)
			q := make(chan string)
			for i := 0; i < p.Workers; i++ {
				wg.Add(1)
				go func() {
					for addr := range q {
						// Whatever the result is, reportit will handle it.
						success, _ := c.Do(addr)
						if success {
							entryList = append(entryList, Entry{Addr: addr, Success: success})
							if p.Verbose {
								log.Printf("[%s] %s running %s", TAG, addr, c.Identifier())
							}
						}
					}
					defer wg.Done()
					return
				}()
			}

			for i := 0; i < len(ipList); i++ {
				ip := ipList[i]
				// Handle IPv6 Addresses.
				if strings.Contains(ip.String(), ":") {
					addrFormat = "[%s]:%d"
				}
				addr := fmt.Sprintf(addrFormat, ip.String(), xport)
				// Check to see if we have already scanned this addr.
				// If we have then skip it.
				if _, ok := scannedList[addr]; ok {
					continue
				} else {
					scannedList[addr] = true
				}
				q <- addr
			}
			close(q)
		}
	}
	wg.Wait()

	// We're finished! Stop the clock.
	finish := time.Now()

	// When all of the workers end gracefully
	// go on to report all of the details of
	// the scan as as per the output format.
	r := Report{
		Entries: entryList,
		Scan: ScanInfo{
			Info: p.Report(),

			Finish: finish,
			Start:  start,
		},
	}
	buf, err := json.MarshalIndent(r, " ", "    ")
	if err != nil {
		log.Fatalf("Failed to produce output report: %s", err)
	}
	fmt.Println(string(buf))
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
