package config

import "embed"

//go:embed internal/*
var SIMPLE embed.FS

const (
	BIND                     = ":8080"
	WORKSPACE_CONFIG         = "workspace.yml"
	WORKSPACE_DB_NAME        = "workspace.db"
	WORKSPACE_API_CONFIG_DIR = "api_configs"
	WORKSPACE_LOCK_FILE      = "workspace.lock"

	Mysql_LinkUrl     = "username:password@protocol(address)/dbname?param=value&param2=value2"
	Oracle_LinkUrl    = "oracle://user:pass@localhost/sid"
	Pgsql_LinkUrl     = "postgres://pqgotest:password@localhost/pqgotest?sslmode=disable"
	Sqlite_LinkUrl    = "file:sqlite.db?cache=shared&mode=rwc"
	Sqlserver_LinkUrl = "sqlserver://username:password@host:port?database=master&param1=value&param2=value"
)
