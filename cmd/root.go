/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "goautomate",
	Short: "Execute commands to targets machines",
	Long: `A simple CLI to execute commands to targets machines through ssh.
			
This CLI is only for daly-basics ro commands in non-production environments.
For example:
	- ls -l /path/of/target/server
	- ls -l /path/of/target/server | grep -R "regex"
	- grep "regex" /path/to/grep
	.
	.
	.
	- ro commands`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.goautomate.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
}
