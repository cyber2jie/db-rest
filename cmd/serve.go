package cmd

import (
	"db-rest/api"
	"db-rest/config"
	"db-rest/env"
	"db-rest/service"
	"db-rest/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
)

var (
	workspace string
	serveCmd  = &cobra.Command{
		Use:     "serve <workspace dir>",
		Short:   "Start the server",
		Long:    `Start a server with workspace `,
		Aliases: []string{"s", "serv"},
		Run: func(cmd *cobra.Command, args []string) {

			if workspace == "" {
				if len(args) <= 0 {
					cmd.Help()
					return
				}
				workspace = args[0]
			}

			util.Log("load workspace %s", workspace)

			workspaceConfig, err := config.ParseWorkSpaceConfig(workspace)

			if err != nil {
				panic(err)
			}

			util.LogSuccess("workspace loaded")

			sc, err := service.GetServiceContext(workspaceConfig, context.Background())
			if err != nil {
				panic(err)
			}

			env.SetEnvVar(env.SERVICE_CONTEXT, sc)

			engine, err := api.New(sc)
			if err != nil {
				panic(err)
			}
			if err := engine.Run(); err != nil {
				panic(err)
			}

		},
	}
)

func init() {
	serveCmd.FParseErrWhitelist.UnknownFlags = true
	serveCmd.Flags().StringP("bind", "b", "", "bind address")
	bind := serveCmd.Flags().Lookup("bind")
	viper.BindPFlag(config.VIPER_KEY_ENGINE_BIND, bind)
	rootCmd.AddCommand(serveCmd)
}
