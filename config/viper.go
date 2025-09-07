package config

import (
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"os"
	"path"
)

const (
	//extract
	VIPER_KEY_MAX_DATASOURCE = "MAX_DATASOURCE"
	MAX_DATASOURCE           = 16

	VIPER_KEY_MAX_CONNECTION = "MAX_CONNECTION"
	MAX_CONNECTION           = 16

	VIPER_KEY_USE_PAGINATION_EXTRACT = "USE_PAGINATION_EXTRACT"
	USE_PAGINATION_EXTRACT           = true

	VIPER_KEY_CHECK_SQL = "CHECK_SQL"
	CHECK_SQL           = true

	VIPER_KEY_USE_TIDB_PARSER = "USE_TIDB_PARSER"
	USE_TIDB_PARSER           = true

	VIPER_KEY_SAVE_BATCH_SIZE = "BATCH_SIZE"
	BATCH_SIZE                = 1000

	VIPER_KEY_MODE_STRICT = "MODE_STRICT"
	MODE_STRICT           = false

	VIPER_KEY_SLOW_SQL_THRESHOLD = "SLOW_SQL_THRESHOLD"
	SLOW_SQL_THRESHOLD           = 3000

	VIPER_KEY_PAGINATION_SQL_LIMIT = "PAGINATION_SQL_LIMIT"
	PAGINATION_SQL_LIMIT           = 10000

	VIPER_KEY_GORM_LOG_INFO_MODE = "GORM_LOG_INFO_MODE"
	GORM_LOG_INFO_MODE           = false

	//engine
	VIPER_KEY_ENGINE_BIND = "ENGINE_BIND"
	ENGINE_BIND           = ""

	VIPER_KEY_ENGINE_DEBUG_MODE = "ENGINE_DEBUG_MODE"
	ENGINE_DEBUG_MODE           = false

	VIPER_KEY_JWT_SECRET = "JWT_SECRET"

	VIPER_KEY_API_MAX_PAGE_SIZE = "API_MAX_PAGE_SIZE"
	API_MAX_PAGE_SIZE           = 200

	VIPER_KEY_MAX_QUERY_CONDITION = "MAX_QUERY_CONDITION"
	MAX_QUERY_CONDITION           = 10
)

func init() {
	viper.SetEnvPrefix("DB-REST")

	viper.SetDefault(VIPER_KEY_MAX_DATASOURCE, MAX_DATASOURCE)
	viper.SetDefault(VIPER_KEY_MAX_CONNECTION, MAX_CONNECTION)
	viper.SetDefault(VIPER_KEY_USE_PAGINATION_EXTRACT, USE_PAGINATION_EXTRACT)
	viper.SetDefault(VIPER_KEY_CHECK_SQL, CHECK_SQL)
	viper.SetDefault(VIPER_KEY_USE_TIDB_PARSER, USE_TIDB_PARSER)
	viper.SetDefault(VIPER_KEY_MODE_STRICT, MODE_STRICT)
	viper.SetDefault(VIPER_KEY_SAVE_BATCH_SIZE, BATCH_SIZE)
	viper.SetDefault(VIPER_KEY_SLOW_SQL_THRESHOLD, SLOW_SQL_THRESHOLD)
	viper.SetDefault(VIPER_KEY_PAGINATION_SQL_LIMIT, PAGINATION_SQL_LIMIT)
	viper.SetDefault(VIPER_KEY_GORM_LOG_INFO_MODE, GORM_LOG_INFO_MODE)

	viper.SetDefault(VIPER_KEY_ENGINE_BIND, ENGINE_BIND)
	viper.SetDefault(VIPER_KEY_ENGINE_DEBUG_MODE, ENGINE_DEBUG_MODE)
	viper.SetDefault(VIPER_KEY_JWT_SECRET, "")
	viper.SetDefault(VIPER_KEY_API_MAX_PAGE_SIZE, API_MAX_PAGE_SIZE)
	viper.SetDefault(VIPER_KEY_MAX_QUERY_CONDITION, MAX_QUERY_CONDITION)

	viper.SetConfigName("db-rest")
	viper.SetConfigType("yaml")
	home, err := os.UserHomeDir()
	if err == nil {
		p := path.Join(home, ".db_rest")
		os.MkdirAll(p, os.ModePerm)
		viper.AddConfigPath(p)
		viper.SafeWriteConfig()
		viper.ReadInConfig()
	}
}

func GetEnvValue[T any](key string) T {
	v := viper.Get(key)
	var zero T
	if v == nil {
		return zero
	}
	return convert(v, any(zero)).(T)
}
func convert(v any, t any) any {
	switch t.(type) {
	case string:
		return cast.ToString(v)
	case bool:
		return cast.ToBool(v)
	case int:
		return cast.ToInt(v)
	case int8:
		return cast.ToInt8(v)
	case int16:
		return cast.ToInt16(v)
	case int32:
		return cast.ToInt32(v)
	case int64:
		return cast.ToInt64(v)
	case uint:
		return cast.ToUint(v)
	case uint32:
		return cast.ToUint32(v)
	case float32:
		return cast.ToFloat32(v)
	case float64:
		return cast.ToInt64(v)
	}
	return v
}
