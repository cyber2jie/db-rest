package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"os"
)

var (
	genDocCmd = &cobra.Command{
		Use: "gendoc [location]",
		Run: func(cmd *cobra.Command, args []string) {
			location := "."
			if len(args) > 0 {
				location = args[0]
			}
			err := doc.GenMarkdownTree(rootCmd, location)

			if err != nil {
				fmt.Fprint(os.Stderr, err)
				os.Exit(1)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(genDocCmd)
}
