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
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/ejcx/map/port"
	"github.com/ejcx/map/scan"
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
			log.Fatalf("Invalid port value passed: %s", err)
		}
		if p < 0 || p > 65535 {
			log.Fatal("Invalid port number: ", port)
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
