package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "video-manager",
	Short: "Idiosyncratic video-manager",
	Long: `My very own system of online video consumption
				  written in Go.
				  Repository located at https://github.com/radoslawg/video-manager/`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
