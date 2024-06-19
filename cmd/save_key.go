/*
Copyright © 2024 Bridge Digital
*/
package cmd

import (
	"fmt"

	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/processes/savekey"
	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services/predefined"
	"github.com/spf13/cobra"
)

// saveKeyCmd represents the saveKey command
var saveKeyCmd = &cobra.Command{
	Use:   "save-key",
	Short: "Add public key!",
	Long:  `Creating/editing a PEM public key.`,
	Run: func(cmd *cobra.Command, args []string) {
		var result string = savekey.Execute(false, "")

		if result != "" {
			fmt.Println(predefined.BuildOk("The public key has been saved successfully"))
		}
	},
}

func init() {
	rootCmd.AddCommand(saveKeyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// saveKeyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// saveKeyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
