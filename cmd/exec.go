/*
Copyright © 2026 alexortegaleonato@gmail.com
SPDX-License-Identifier: MIT
*/
package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	customErr "github.com/CALEXCO/goautomate/pkg/errors"
	"github.com/CALEXCO/goautomate/pkg/server"
	"github.com/spf13/cobra"
)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "exec command to target server",
	Long:  `execute a command to a target group of servers`,
	Run: func(cmd *cobra.Command, args []string) {
		f, err := cmd.Flags().GetString("file")
		if err != nil {
			panic(err)
		}
		c, err := cmd.Flags().GetString("command")
		if err != nil {
			panic(err)
		}
		targetServer, err := server.BuildServersFromYaml(f)
		if err != nil {
			if gae, ok := errors.AsType[*customErr.GoAutomateErrors](err); ok {
				os.Exit(gae.Code)
			}
		}

		var wg sync.WaitGroup
		for _, ts := range targetServer.Servers {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			ts.SshConfig.Ctx = ctx
			wg.Add(1)
			go func(*context.Context) {
				defer wg.Done()
				fmt.Println("Executing:", c, "in", ts.Ip)
				output, err := ts.ConnectAndExecute(c)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					if gae, ok := errors.AsType[*customErr.GoAutomateErrors](err); ok {
						os.Exit(gae.Code)
					}
				}
				fmt.Println(output)
			}(&ctx)
		}
		wg.Wait()
	},
}

func init() {
	rootCmd.AddCommand(execCmd)

	rootCmd.PersistentFlags().String("command", "", "Command to execute to the server/s. Ex: \"ls -l .* | grep bash\"")
	rootCmd.PersistentFlags().String("file", "", `yaml file with the hosts you want to target with the command.
For example:
target_servers:
target_servers:
- ip: 1.1.1.1
  user: vagrant
  password: vagrant
- ip: 1.1.1.1
  user: vagrant
  password: vagrant`)
}
