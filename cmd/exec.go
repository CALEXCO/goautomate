/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "exec command to target server",
	Long:  `execute a command to a target group of servers`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Executing")
	},
}

// TODO: execute command in a server using ssh

func init() {
	rootCmd.AddCommand(execCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// execCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// execCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().String("command", "", "Command to execute to the server/s. Ex: \"ls -l .* | grep bash\"")
	rootCmd.PersistentFlags().String("file", "", `yaml file with the hosts you want to target with the command.
For example:
target_servers:
	- ip: 1.1.1.1
	user: vagrant
	password: vagrant
	- ip: 1.1.1.1
	user: vagrant
	password: vagrant`)
}
