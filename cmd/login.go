/*
Copyright © 2024 Bridge Digital
*/
package cmd

import (
	"fmt"

	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/processes/login"
	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authorize to service!",
	Long:  `Authorize to service!`,
	Run: func(cmd *cobra.Command, args []string) {
		var result string = login.Execute(cmd)
		fmt.Println(result)
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
