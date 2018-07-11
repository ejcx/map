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
	"github.com/ejcx/map/scan/ssh"
	"github.com/spf13/cobra"
)

const (
	defaultSshPort = "22"
)

// sshCmd represents the ssh command
var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "Scan for open ssh instances.",
	Long: `SSH needs no introduction. It's used to
access just about ever linux system that
exists. It can use key or password based
authentication`,
	Run: func(cmd *cobra.Command, args []string) {
		doer := &ssh.SshDoer{
			Username: username,
			Password: password,
		}
		root(cmd, defaultSshPort, doer)
	},
}

func init() {
	addUsername(sshCmd)
	addPassword(sshCmd)
	RootCmd.AddCommand(sshCmd)
}
