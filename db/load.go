package db

import (
	"db-rest/config"
	"db-rest/util"
	"fmt"
	"github.com/spf13/cast"
	"gorm.io/gorm"
)

type Loader interface {
	Load() (uint32, error)
}

type ProcessFn = func(row []*DataRow) error

type BaseLoader struct {
	dbConnector *DbConnector
	db          *gorm.DB
	config      *config.ApiConfig
	processFn   ProcessFn
}
type AllDataLoader struct {
	BaseLoader
}

var isStrict = config.GetEnvValue[bool](config.VIPER_KEY_MODE_STRICT)
var pageSize = config.GetEnvValue[int](config.VIPER_KEY_PAGINATION_SQL_LIMIT)

func (loader *AllDataLoader) Load() (uint32, error) {
	util.Log("load data by all strategy")
	sqlExecutor := loader.dbConnector.Db
	sql := loader.config.Sql
	columns := loader.config.Columns

	rows, err := sqlExecutor.Queryx(sql)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	dataRows := []*DataRow{}

	row := make(map[string]interface{})

	row_num := 1
	for rows.Next() {
		err = rows.MapScan(row)
		if err != nil {
			if isStrict {
				return 0, err
			}
			util.LogWarn("mapscan row error: %v", err)
			continue
		}
		dataRow, err := mapToDataRow(columns, row)
		if err != nil {
			if isStrict {
				return 0, err
			}
			util.LogWarn("map row to data row error: %v", err)
			continue
		}
		dataRow.RowNum = row_num
		dataRows = append(dataRows, dataRow)
		row_num++

	}

	if loader.processFn != nil {
		err := loader.processFn(dataRows)
		if err != nil {
			return 0, err
		}
	}

	return uint32(len(dataRows)), nil
}

func mapToDataRow(columns []string, row map[string]interface{}) (*DataRow, error) {
	dataColumns := []*DataColumn{}
	for _, c := range columns {
		v := row[c]
		dataColumns = append(dataColumns, &DataColumn{
			Name:  c,
			Value: columnToString(v),
		})
	}
	return &DataRow{
		Columns: dataColumns,
	}, nil
}

// TODO 按数据类型优化存储
func columnToString(v any) string {
	return cast.ToString(v)
}

func NewAllDataLoader(dbConnector *DbConnector, db *gorm.DB, config *config.ApiConfig, processFn ProcessFn) *AllDataLoader {
	return &AllDataLoader{
		BaseLoader: BaseLoader{
			dbConnector: dbConnector,
			db:          db,
			config:      config,
			processFn:   processFn,
		},
	}
}

type PaginationDataLoader struct {
	BaseLoader
}

func (loader *PaginationDataLoader) _load(sql string, row_num int) (uint32, error) {
	sqlExecutor := loader.dbConnector.Db
	columns := loader.config.Columns
	rows, err := sqlExecutor.Queryx(sql)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	dataRows := []*DataRow{}

	row := make(map[string]interface{})

	for rows.Next() {
		err = rows.MapScan(row)
		if err != nil {
			if isStrict {
				return 0, err
			}
			util.LogWarn("mapscan row error: %v", err)
			continue
		}
		dataRow, err := mapToDataRow(columns, row)
		if err != nil {
			if isStrict {
				return 0, err
			}
			util.LogWarn("map row to data row error: %v", err)
			continue
		}
		dataRow.RowNum = row_num
		dataRows = append(dataRows, dataRow)
		row_num++

	}

	if loader.processFn != nil {
		err := loader.processFn(dataRows)
		if err != nil {
			return 0, err
		}
	}

	return uint32(len(dataRows)), nil
}

func (loader *PaginationDataLoader) Load() (uint32, error) {
	util.Log("load data by pagination strategy")
	sqlExecutor := loader.dbConnector.Db
	sql := loader.config.Sql
	columns := loader.config.Columns
	driver := loader.dbConnector.Driver
	dialect := driver.Dialect

	logTag := fmt.Sprintf("db[%s],api[%s]", loader.config.DbName, loader.config.Name)

	var total uint32

	if dialect != nil {
		countSql := dialect.getCountSql(sql)
		util.Log("%s count sql: %s", logTag, countSql)
		crow, err := sqlExecutor.Queryx(countSql)
		if err != nil {
			return 0, err
		}
		defer crow.Close()

		var count int

		if crow.Next() {
			err := crow.Scan(&count)
			if err != nil {
				return 0, err
			}
		}

		if count > 0 {
			util.LogSuccess("%s total count: %d", logTag, count)
			totalPages := (count + pageSize - 1) / pageSize
			util.LogSuccess("%s total pages: %d", logTag, totalPages)
			for page := 1; page <= totalPages; page++ {
				util.Log("%s load page %d", logTag, page)
				pageSql := dialect.getPaginationSql(sql, columns, page, pageSize)
				util.Log("%s page sql: %s", logTag, pageSql)
				row_num := (page-1)*pageSize + 1
				_total, err := loader._load(pageSql, row_num)
				if err != nil {
					if isStrict {
						return total, err
					}
					util.LogWarn("%s load page %d occur error: %v", logTag, page, err)
				}
				total += _total
			}

		}

	} else {
		util.Log("not support pagination dialect for driver [%s]", driver)
		return loader._load(sql, 1)
	}
	return total, nil
}
func NewPaginationDataLoader(dbConnector *DbConnector, db *gorm.DB, config *config.ApiConfig, processFn ProcessFn) *PaginationDataLoader {
	return &PaginationDataLoader{
		BaseLoader: BaseLoader{
			dbConnector: dbConnector,
			db:          db,
			config:      config,
			processFn:   processFn,
		},
	}
}
