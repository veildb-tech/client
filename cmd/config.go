/*
Copyright © 2024 Bridge Digital
*/
package cmd

import (
	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/processes/config"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Set new path for saving dumps",
	Long:  `Sets a new default path for saving database dumps`,
	Run: func(cmd *cobra.Command, args []string) {
		dumpPath, _ := cmd.Flags().GetString("dump-path")

		config.Execute(dumpPath)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.Flags().StringP("dump-path", "d", "", "Set new path for saving dumps")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
