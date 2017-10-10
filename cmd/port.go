// Copyright © 2017 Evan Johnson <evan@twiinsen.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/ejcx/map/scan"
	"github.com/ejcx/map/scan/port"
	"github.com/spf13/cobra"
)

// portCmd represents the port command
var (
	portCmd = &cobra.Command{
		Use:   "port",
		Short: "port scan a subnet or host",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: run,
	}
	portCsv  string
	netCsv   string
	portNums []int
	cidrs    []*port.Cidr
)

func parsePortRange(pr string) ([]int, error) {
	s := strings.Split(pr, "-")

	// If it does has a dash, it should be a small number
	// and a larger number. The reverse is invalid.
	if len(s) != 2 {
		return nil, errors.New("Port range must be two ints. Format 1-2")
	}
	a, err := strconv.Atoi(s[0])
	if err != nil {
		return nil, fmt.Errorf("Number must be passed for a port range: %s, %s", a, err)
	}
	b, err := strconv.Atoi(s[1])
	if err != nil {
		return nil, fmt.Errorf("Number must be passed for a port range: %s, %s", b, err)
	}
	if b <= a {
		return nil, fmt.Errorf("Invalid port range. %s is greater than or equal to %s", a, b)
	}
	var (
		nums []int
	)
	for i := a; i < b; i++ {
		if i < 0 || i > 65535 {
			return nil, fmt.Errorf("Invalid port number: %s")
		}
		nums = append(nums, i)
	}
	return nums, nil
}

func run(cmd *cobra.Command, args []string) {

	portList := strings.Split(portCsv, ",")
	netList := strings.Split(netCsv, ",")

	// We don't want to begin scanning and be surprised by an entry
	// that is not a valid HTTP port. Do a quick once over and verify
	// that they are all integers that are greater than 0 and less than
	// 65536
	for _, port := range portList {
		p, err := strconv.Atoi(port)
		if err != nil {
			// If the port list did not
			ports, err := parsePortRange(port)
			if err != nil {
				log.Fatal(err)
			}
			portNums = append(portNums, ports...)
			continue
		} else {
			if p < 0 || p > 65535 {
				log.Fatalf("Invalid port number: %s", port)
			}
		}
		portNums = append(portNums, p)
	}

	// We also want to validate the IP addresses list we got. We can
	// do this by parsing them as a CIDR. This may not work for IPv6
	// depending on what the underlying parse function is doing.
	for _, cidr := range netList {
		i, s, err := net.ParseCIDR(cidr)
		if err != nil {
			log.Fatalf("Invalid cidr passed as arg: %s", err)
		}
		cidrs = append(cidrs, &port.Cidr{IP: i, IPNet: s})
	}

	// Run the scan.
	s := &port.PortScan{
		CIDR:  cidrs,
		Ports: portNums,
	}
	scan.Do(s)
}

func init() {
	portCmd.Flags().StringVarP(&portCsv, "ports", "p", "80,443", "The set of ports to scan")
	portCmd.Flags().StringVarP(&netCsv, "nets", "n", "10.0.0.0/24", "The cidrs to scan")
	RootCmd.AddCommand(portCmd)
}
