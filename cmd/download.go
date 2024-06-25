/*
Copyright © 2024 Bridge Digital
*/
package cmd

import (
	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/processes/download"
	"github.com/spf13/cobra"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download dump",
	Long:  `Downloading a dump of the database.`,
	Run: func(cmd *cobra.Command, args []string) {
		dumpUid, _ := cmd.Flags().GetString("dump-uid")
		dbUid, _ := cmd.Flags().GetString("db-uid")

		download.Execute(dbUid, dumpUid)
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().String("dump-uid", "", "Enter DB Dump UID")
	downloadCmd.Flags().String("db-uid", "", "Enter DB UID")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downloadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downloadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
