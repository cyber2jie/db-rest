package db

import (
	"db-rest/config"
	"github.com/jmoiron/sqlx"
	"gorm.io/gorm"
)

type DbModel struct {
	gorm.Model
}

// entity
type DbApiConfigCollection struct {
	DbModel
	Name        string
	Description string
}

type DbApiConfig struct {
	DbModel
	CollectionId uint
	Name         string
	Sql          string
	Description  string
	DbName       string
	Columns      string `gorm:"columns"`
}

const (
	DB_DATA_TYPE_STRING = "string"
	DB_DATA_TYPE_NUMBER = "number"
)

type DbData struct {
	DbModel
	DbApiConfigId uint `gorm:"db_api_config_id"`
	Table         string
	RowNum        int
	Column        string `gorm:"column:column;index:idx_column_value"`
	Value         string `gorm:"column:value;index:idx_column_value"`
	LargeValue    string `gorm:"large_value;type:text;"` //大值
	ValueType     string `gorm:"value_type"`
}

// Db
type DbConnector struct {
	Id     string
	Config *config.DbConfig
	Db     *sqlx.DB
	Driver *DbDriver
}

func (dbconn *DbConnector) Close() error {
	return dbconn.Db.Close()
}

type DataRow struct {
	RowNum  int           `json:"row_num"`
	Columns []*DataColumn `json:"columns"`
}

type DataColumn struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
