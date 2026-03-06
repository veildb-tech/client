/*
Copyright © 2024 Bridge Digital
*/
package cmd

import (
	"github.com/dbvisor-pro/client/processes/download"
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
		latestDump, _ := cmd.Flags().GetBool("latest-dump")

		download.Execute(dbUid, dumpUid, latestDump)
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().String("dump-uid", "", "Enter DB Dump UID")
	downloadCmd.Flags().String("db-uid", "", "Enter DB UID")
	downloadCmd.Flags().Bool("latest-dump", false, "Download the latest dump without selection")
}
