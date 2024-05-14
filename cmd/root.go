/*
Copyright © 2024 Bridge Digital
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "db-manager-client-cli-go",
	Short: "Console Tool",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.db-manager-client-cli-go.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().BoolP("quiet", "q", false, "Do not output any message")
	rootCmd.Flags().BoolP("version", "V", false, "Display this application version")
	//rootCmd.Flags().BoolP("toggle", "--ansi", false, "Force (or disable --no-ansi) ANSI output")
	rootCmd.Flags().BoolP("no-interaction", "n", false, "Do not ask any interactive question")
	rootCmd.Flags().BoolP("verbose", "v", false, "Increase the verbosity of messages: 1 for normal output, 2 for more verbose output and 3 for debug")
}
