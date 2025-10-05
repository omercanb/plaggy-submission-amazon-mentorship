package cmd

import (
	"aiplag-agent/common/config"
	"aiplag-agent/common/db"
	"fmt"

	"github.com/spf13/cobra"
)

// lsCmd lists all watched directories
var lsCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all watched directories",
	Long:  `Lists all directories that are currently being watched by the system.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Call your function to get assignment paths
		editHistory, err := db.NewEditHistoryStore(config.DBPath())
		assignmentPaths, err := editHistory.GetAssignmentFullPaths()
		if err != nil {
			fmt.Println("Error fetching watched directories:", err)
			return
		}

		if len(assignmentPaths) == 0 {
			fmt.Println("No watched directories found.")
			return
		}

		fmt.Println("Watched directories:")
		for _, path := range assignmentPaths {
			fmt.Println(" -", path)
		}
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}
