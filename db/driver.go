package db

import (
	"database/sql"
)

import _ "github.com/go-sql-driver/mysql"

import _ "github.com/godror/godror"

import _ "github.com/lib/pq"

import _ "github.com/denisenkom/go-mssqldb"

import _ "github.com/mattn/go-sqlite3"

var dbDrivers = []DbDriver{
	{Name: "sqlite3", DriverName: "sqlite3", Dialect: sqliteDbDialect},
	{Name: "sqlite", DriverName: "sqlite3", Dialect: sqliteDbDialect},
	{Name: "postgres", DriverName: "postgres", Dialect: postgresDbDialect},
	{Name: "oracle", DriverName: "godror", Dialect: oracleDbDialect},
	{Name: "oracle11g", DriverName: "godror", Dialect: oracle11gDbDialect},
	{Name: "mysql", DriverName: "mysql", Dialect: mysqlDbDialect},
	{Name: "mariadb", DriverName: "mysql", Dialect: mysqlDbDialect},
	{Name: "tidb", DriverName: "mysql", Dialect: mysqlDbDialect},
	{Name: "mssql", DriverName: "sqlserver", Dialect: sqlServer2012DbDialect},
	{Name: "sqlserver", DriverName: "sqlserver", Dialect: sqlServerDbDialect},
	{Name: "sqlserver2012", DriverName: "sqlserver", Dialect: sqlServer2012DbDialect},
}

func Drivers() []string {
	return sql.Drivers()
}

type DbDriver struct {
	Name       string
	DriverName string
	Dialect    DbDialect
}

func GetDbDriver(name string) *DbDriver {
	for _, driver := range dbDrivers {
		if driver.Name == name {
			return &driver
		}
	}
	return nil
}
