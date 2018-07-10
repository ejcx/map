// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
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
	m "github.com/ejcx/map/scan/mysql"
	"github.com/spf13/cobra"
)

const (
	defaultMysqlPort = "3306"
)

var (
	username string
	protocol string
	dbname   string
)

// mysqlCmd represents the mysql command
var mysqlCmd = &cobra.Command{
	Use:   "mysql",
	Short: "Scan for mysql servers",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		doer := &m.MysqlDoer{
			Username: username,
			Password: password,
			Protocol: protocol,
			DBName:   dbname,
		}
		root(cmd, defaultMysqlPort, doer)
	},
}

func init() {
	mysqlCmd.Flags().StringVarP(&username, "username", "u", "root", "The username used to connect to the db")
	addPassword(mysqlCmd)
	mysqlCmd.Flags().StringVarP(&protocol, "protocol", "q", "tcp", "The protocol used to connect to the db")
	mysqlCmd.Flags().StringVarP(&dbname, "dbname", "d", "mysql", "The database name used to connect to the db")
	RootCmd.AddCommand(mysqlCmd)
}
