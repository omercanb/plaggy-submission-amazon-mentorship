package cmd

import (
	tcpclient "aiplag-agent/cli/tcp-client"
	"aiplag-agent/common/config"
	"aiplag-agent/common/db"
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// stopWatchingCmd represents the stop-watching command
var stopWatchingCmd = &cobra.Command{
	Use:   "stop-watching",
	Short: "Stop watching a directory",
	Long:  `Stop watching a directory previously added with the watch command. Optionally delete stored edits.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Load watched directories from edit history
		eh, err := db.NewEditHistoryStore(config.DBPath())
		if err != nil {
			fmt.Println("Failed to access edit history:", err)
			return
		}

		assignmentPaths, err := eh.GetAssignmentFullPaths()
		if err != nil {
			fmt.Println("Failed to get watched directories:", err)
			return
		}

		if len(assignmentPaths) == 0 {
			fmt.Println("No directories are currently being watched.")
			return
		}

		// Prompt to select directory to stop watching
		selectDirPrompt := promptui.Select{
			Label:    "Select directory to stop watching",
			Items:    assignmentPaths,
			HideHelp: true,
		}
		_, dirToStop, err := selectDirPrompt.Run()
		fmt.Println(dirToStop)
		if err != nil {
			fmt.Println("Cancelled.")
			return
		}

		fmt.Println("Stopping watch for:", dirToStop)

		// Prompt to optionally delete stored edits
		deletePrompt := promptui.Select{
			Label:    "Also delete stored edits? You can restore this later with 'plaggy watch'",
			Items:    []string{"Yes", "No"},
			HideHelp: true,
		}
		_, deleteChoice, _ := deletePrompt.Run()

		if deleteChoice == "Yes" {
			eh.DeleteEditsByFullPath(dirToStop)
			if err != nil {
				fmt.Println("Failed to delete edits:", err)
			} else {
				fmt.Println("Stored edits deleted for", dirToStop)
			}
		} else {
			fmt.Println("Stored edits retained. You can restore watching later using 'plaggy watch'.")
		}

		// Send command to TCP daemon
		resp, err := tcpclient.SendCommand('X', dirToStop)
		if err != nil {
			// fmt.Println("Error sending command:", err)
			return
		}

		switch resp {
		case 'A':
			// fmt.Println("Stopped watching path!")
		case 'R':
			// fmt.Println("Error while stopping watching path!")
		default:
			fmt.Println("Unknown response from daemon:", resp)
		}
		// fmt.Println("Stop-watching command sent to daemon for:", dirToStop)
	},
}

func init() {
	rootCmd.AddCommand(stopWatchingCmd)
}
