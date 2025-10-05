/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	tcpclient "aiplag-agent/cli/tcp-client"

	"github.com/spf13/cobra"
)

// watchCmd represents the watch command
var watchCmd = &cobra.Command{
	Use:   "watch [path]",
	Short: "Starts watching files in the specified directory",
	Long:  `.`,
	Args:  cobra.MaximumNArgs(1), // allow at most one argument
	Run: func(cmd *cobra.Command, args []string) {
		// Get current working directory
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Failed to get current directory:", err)
			return
		}

		// Resolve path
		var pathToWatch string
		if len(args) == 0 || args[0] == "" {
			pathToWatch = cwd
		} else {
			pathToWatch = filepath.Join(cwd, args[0])
		}
		// Check if path exists
		info, err := os.Stat(pathToWatch)
		if os.IsNotExist(err) {
			fmt.Println("Path does not exist:", pathToWatch)
			return
		}
		if err != nil {
			fmt.Println("Error accessing path:", err)
			return
		}
		if !info.IsDir() {
			fmt.Println("Specified path is not a directory:", pathToWatch)
			return
		}

		fmt.Println("Watching path:", pathToWatch)

		// Send command to TCP daemon
		resp, err := tcpclient.SendCommand('W', pathToWatch)
		if err != nil {
			fmt.Println("Error sending command:", err)
			return
		}

		switch resp {
		case 'A':
			fmt.Println("Started watching path!")
		case 'R':
			fmt.Println("Error while watching path!")
		default:
			fmt.Println("Unknown response from daemon:", resp)
		}
	},
}

func init() {
	rootCmd.AddCommand(watchCmd)
}
