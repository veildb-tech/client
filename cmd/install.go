/*
Copyright © 2024 Bridge Digital
*/
package cmd

import (
	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/processes/install"
	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install application",
	Long:  `Installing a console application.`,
	Run: func(cmd *cobra.Command, args []string) {
		install.Execute()
	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// installCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
