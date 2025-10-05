/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"aiplag-agent/cli/models"
	"aiplag-agent/common/api"
	"aiplag-agent/common/config"
	"aiplag-agent/common/db"
	"errors"
	"fmt"
	"slices"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// submitCmd represents the submit command
var submitCmd = &cobra.Command{
	Use:   "submit",
	Short: "Submit an assignment to the backend",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// Load session info
		email := viper.GetString("session.email")
		token := viper.GetString("session.token")
		if email == "" || token == "" {
			fmt.Println("No session found. Please login first.")
			return
		}

		eh, err := db.NewEditHistoryStore(config.DBPath())
		assignmentPaths, err := eh.GetAssignmentFullPaths()

		selectDirectoryToSubmitPrompt := promptui.Select{
			Label:    "Select Directory To Submit",
			Items:    assignmentPaths,
			HideHelp: true,
		}

		_, dirToSubmit, err := selectDirectoryToSubmitPrompt.Run()
		if err != nil {
			return
		}

		assignments, err := api.FetchAssignments(email, token)
		if err != nil {
			if errors.Is(err, api.ServerError) {
				fmt.Println("Server unavailable, please try again later")
			} else {
				fmt.Println("Unknown error occured, please try again later")
			}
			return
		}
		assignmentTitles := []string{}
		for _, assignment := range assignments {
			assignmentTitles = append(assignmentTitles, assignment.Title)
		}

		prompt := promptui.Select{
			Label:    "Select The Assignment To Submit",
			Items:    assignmentTitles,
			HideHelp: true,
		}
		_, selectedAssignmentTitle, err := prompt.Run()
		selectedAssignmentIdx := slices.IndexFunc(assignments, func(a models.Assignment) bool {
			return a.Title == selectedAssignmentTitle
		})
		selectedAssignment := assignments[selectedAssignmentIdx]
		err = api.SubmitEdits(selectedAssignment.ID, eh, dirToSubmit, token)
		if err == nil {
			fmt.Println("Assignment submitted!")
		} else {
			if errors.Is(err, api.ServerError) {
				fmt.Println("Server unavailable, please try again later")
			} else {
				fmt.Println("Submission failed! Please try again")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(submitCmd)
}
