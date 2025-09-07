package cmd

import (
	"db-rest/config"
	"db-rest/db"
	"db-rest/util"
	"github.com/gofrs/flock"
	"github.com/spf13/cobra"
	"os"
	"path"
)

var (
	initCmd = &cobra.Command{
		Use:   "init <init dir> ",
		Short: "init workspace",
		Long:  `init db-rest workspace,which includes simple config file,cache file,etc `,
		Run: func(cmd *cobra.Command, args []string) {
			workspace := ""

			if len(args) > 0 {
				workspace = args[0]
			}
			if workspace == "" {
				workspace = "."
			}

			if config.CheckWorkSpaceExists(workspace) {
				util.LogFatal("%s is exist workspace", workspace)
			}

			if !util.FileExists(workspace) {
				os.MkdirAll(workspace, os.ModePerm)
			}

			lock := flock.New(path.Join(workspace, config.WORKSPACE_LOCK_FILE))

			locked, err := lock.TryLock()

			if err != nil {
				util.LogFatal("lock workspace dir error:%v", err)
			}

			if locked {
				if err := config.InitWorkSpace(workspace); err != nil {
					util.LogFatal("init workspace error:%v", err)
				}

				if err := db.InitWorkSpaceDb(workspace); err != nil {
					util.LogFatal("init workspace db error:%v", err)
				}
				util.LogSuccess("init workspace success")

				lock.Unlock()
			}

		},
	}
)

func init() {
	rootCmd.AddCommand(initCmd)
}
