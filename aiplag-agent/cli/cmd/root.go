package cmd

import (
	"aiplag-agent/common/config"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "plaggy help",
	Short: "plaggy cli",
	Long:  `This cli is the student interface for plaggy which manages file tracking and submissions`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to the plaggy cli")
		fmt.Println("Run plaggy help to see a list of useful commands")
	},
}

func init() {
	viper.SetConfigFile(config.ConfigPath())
	viper.SetConfigType("yaml")

	// Read config (if exists)
	if err := viper.ReadInConfig(); err == nil {
		// Optional: feedback to user
		// fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
