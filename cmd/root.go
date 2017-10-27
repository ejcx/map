// Copyright © 2017 Evan Johnson <e@ejj.io>
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
	"github.com/spf13/cobra"
)

var (
	cfgFile   string
	scanTypes []string

	PasswordAdded = false
)

// RootCmd represents the base command when called without any subcommands
var (
	RootCmd = &cobra.Command{
		Use:   "map",
		Short: "map is a simple and powerful network scanning tool",
		Long: `You can use map to scan a network for open ports
or available services on the network.`,
		//Run: run,
	}
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
	portCsv  string
	netCsv   string
	workers  int
	portNums []int
	cidrs    []*scan.Cidr

	password string
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func addPassword() {
	if !PasswordAdded {
		addFlag("password", "s", "", "The password attempt to use in the scan", &password)
	}
	PasswordAdded = true
}

// addFlag is a small abstraction around adding flags.
// We want to only add a flag one time.
func addFlag(longOpt, opt, def, desc string, v *string) {
	RootCmd.PersistentFlags().StringVarP(v, longOpt, opt, def, desc)
}

func init() {
	//RootCmd.PersistentFlags().StringVarP(&portCsv, "ports", "p", "80,443", "The set of ports to scan")
	RootCmd.PersistentFlags().StringVarP(&netCsv, "net", "n", "", "The cidrs to scan")
	RootCmd.PersistentFlags().StringVarP(&portCsv, "ports", "p", "0", "The set of ports to scan")
	RootCmd.PersistentFlags().IntVarP(&workers, "workers", "w", 50, "The number of goroutinues to use for parallel scanning")
	//cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	// RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.map.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
// func initConfig() {
// 	if cfgFile != "" {
// 		// Use config file from the flag.
// 		viper.SetConfigFile(cfgFile)
// 	} else {
// 		// Find home directory.
// 		home, err := homedir.Dir()
// 		if err != nil {
// 			log.Fatal(err)
// 		}
//
// 		// Search config in home directory with name ".map" (without extension).
// 		viper.AddConfigPath(home)
// 		viper.SetConfigName(".map")
// 	}
//
// 	viper.AutomaticEnv() // read in environment variables that match
//
// 	// If a config file is found, read it in.
// 	//if err := viper.ReadInConfig(); err == nil {
// 	//	fmt.Println("Using config file:", viper.ConfigFileUsed())
// 	//}
// }

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

func root(cmd *cobra.Command, port string, f scan.Connector) {

	if portCsv == "0" {
		portCsv = port
	}
	portList := strings.Split(portCsv, ",")
	netList := strings.Split(netCsv, ",")

	// We don't want to begin scanning and be surprised by an entry
	// that is not a valid HTTP port. Do a quick once over and verify
	// that they are all integers that are greater than 0 and less than
	// 65536
	for _, xport := range portList {
		p, err := strconv.Atoi(xport)
		if err != nil {
			// If the port list did not
			ports, err := parsePortRange(xport)
			if err != nil {
				log.Fatal(err)
			}
			portNums = append(portNums, ports...)
			continue
		} else {
			if p < 0 || p > 65535 {
				log.Fatalf("Invalid port number: %s", xport)
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
		cidrs = append(cidrs, &scan.Cidr{IP: i, IPNet: s, Original: cidr})
	}

	scans := []string{}

	// Run the scan.
	s := &scan.NetScan{
		CIDR:    cidrs,
		Ports:   portNums,
		Scans:   scans,
		Workers: workers,
	}
	scan.Do(s, []scan.Connector{f})
}
