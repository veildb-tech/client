/*
Copyright © 2024 Bridge Digital
*/
package cmd

import (
	"github.com/dbvisor-pro/client/processes/config"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure application settings",
	Long:  `Configure application settings including service URL and default path for saving database dumps`,
	Run: func(cmd *cobra.Command, args []string) {
		dumpPath, _ := cmd.Flags().GetString("dump-path")
		serviceUrl, _ := cmd.Flags().GetString("url")

		config.Execute(dumpPath, serviceUrl)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.Flags().StringP("dump-path", "d", "", "Set new path for saving dumps")
	configCmd.Flags().StringP("url", "u", "", "Set service URL (e.g., https://app.veildb.com)")
}
