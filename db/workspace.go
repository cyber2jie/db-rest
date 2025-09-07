package db

import (
	"db-rest/config"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"path"
	"time"
)

var slowSqlThreshold int64 = config.GetEnvValue[int64](config.VIPER_KEY_SLOW_SQL_THRESHOLD)

var gormLogModeInfo bool = config.GetEnvValue[bool](config.VIPER_KEY_GORM_LOG_INFO_MODE)

func GetWorkSpaceGormDb(p string) (*gorm.DB, error) {
	workSpaceDbPath := path.Join(p, config.WORKSPACE_DB_NAME)

	logLevel := logger.Warn

	if gormLogModeInfo {
		logLevel = logger.Info
	}

	gLogger := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
		SlowThreshold:             time.Duration(slowSqlThreshold) * time.Millisecond,
		LogLevel:                  logLevel,
		IgnoreRecordNotFoundError: false,
		Colorful:                  true,
	})

	_db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?cache=shared&mode=rwc", workSpaceDbPath)), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: "workspace_",
		},
		Logger: gLogger,
	})
	if err != nil {
		return nil, err
	}
	return _db, nil
}
