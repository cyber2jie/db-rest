package db

import (
	"db-rest/config"
	"db-rest/util"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"gorm.io/gorm"
	"strings"
)

func initDbDriver() {
	util.Log("可用数据库驱动:%v", strings.Join(Drivers(), ","))
}

func getDbConnector(dbconfig *config.DbConfig) (*DbConnector, error) {
	driver := GetDbDriver(dbconfig.DbType)
	if driver == nil {
		return nil, errors.New(fmt.Sprintf("dbtype %s not support", dbconfig.DbType))
	}

	sqlxDb, err := sqlx.Open(driver.DriverName, dbconfig.LinkUrl)

	if err != nil {
		return nil, err
	}

	//ping datasource
	err = sqlxDb.Ping()
	if err != nil {
		return nil, err
	}

	return &DbConnector{Id: dbconfig.Name,
		Db:     sqlxDb,
		Config: dbconfig,
		Driver: driver,
	}, nil
}

func CleanAllWorkSpaceData(gdb *gorm.DB) error {
	session := gdb.Session(&gorm.Session{AllowGlobalUpdate: true})
	err := session.Unscoped().Delete(&DbApiConfigCollection{}).Error
	if err != nil {
		return err
	}
	err = session.Unscoped().Delete(&DbApiConfig{}).Error
	if err != nil {
		return err
	}
	return session.Unscoped().Delete(&DbData{}).Error
}

func CountDbApiConfigCollectionByName(gdb *gorm.DB, name string) (int64, error) {
	var count int64
	res := gdb.Model(&DbApiConfigCollection{}).Where("name = ?", name).Count(&count)
	if res.Error != nil {
		return 0, res.Error
	}
	return count, nil
}

func GetDbApiConfigCollectionByName(gdb *gorm.DB, name string) (*DbApiConfigCollection, error) {
	var dbApiConfigCollection DbApiConfigCollection
	res := gdb.Where("name = ?", name).First(&dbApiConfigCollection)
	if res.Error != nil {
		return nil, res.Error
	}
	return &dbApiConfigCollection, nil
}

func SaveDbApiConfigCollection(gdb *gorm.DB, dbApiConfigCollection *DbApiConfigCollection) error {
	res := gdb.Create(dbApiConfigCollection)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
