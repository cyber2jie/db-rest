package service

import "db-rest/config"

const (
	DEFAULT_PAGE      = 1
	DEFAULT_PAGE_SIZE = 10
)

var (
	MAX_PAGE_SIZE = config.GetEnvValue[int](config.VIPER_KEY_API_MAX_PAGE_SIZE)
)
