package cmd

import (
	"db-rest/config"
	"db-rest/db"
	"db-rest/env"
	"db-rest/util"
	"github.com/gofrs/flock"
	"github.com/spf13/cobra"
	"path"
)

var (
	extractCmd = &cobra.Command{
		Use:   "extract <workspace dir> ",
		Short: "extract data to  workspace",
		Long:  `use config file,extract data to workspace `,
		Run: func(cmd *cobra.Command, args []string) {
			workspace := ""

			if len(args) > 0 {
				workspace = args[0]
			}
			if workspace == "" {
				workspace = "."
			}

			env.SetEnvVar(env.WORKSPACE, workspace)

			if !config.CheckWorkSpaceExists(workspace) {
				util.LogFatal("%s is not a  workspace", workspace)
			}

			lock := flock.New(path.Join(workspace, config.WORKSPACE_LOCK_FILE))

			locked, err := lock.TryLock()

			if err != nil {
				util.LogFatal("lock workspace dir error: %s", err)
			}

			defer func() {
				if locked {
					lock.Unlock()
				}
				if err := recover(); err != nil {
					util.LogError("workspace %s extract error: %s", workspace, err)
					panic(err)
				}
			}()

			config, err := config.ParseWorkSpaceConfig(workspace)

			if locked {

				util.Log("workspace %s begin extract", workspace)

				err = db.Extract(config)

				if err != nil {
					util.LogFatal("workspace %s extract error: %s", workspace, err)
					return
				}
				util.LogSuccess("workspace %s extract success", workspace)
			}

		},
	}
)

func init() {
	rootCmd.AddCommand(extractCmd)
}
