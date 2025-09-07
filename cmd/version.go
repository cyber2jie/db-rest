package cmd

import (
	"db-rest/version"
	"fmt"
	"github.com/spf13/cobra"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of db-rest",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("db-rest version:", version.Version)
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
}
