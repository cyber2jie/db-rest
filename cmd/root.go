package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	rootCmd = &cobra.Command{
		Use:   "db-rest",
		Short: " make your database accessible via REST API ",
		Long:  " db-rest is a tool which can make your database accessible via REST API, includes table schema and exposes table data。",
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

}
