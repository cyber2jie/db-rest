package service

import (
	"context"
	"db-rest/config"
	"db-rest/db"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type ServiceContext struct {
	WorkSpace *config.WorkSpaceParsedConfig
	Context   context.Context
	DB        *gorm.DB
}

func GetServiceContext(workSpace *config.WorkSpaceParsedConfig, c context.Context) (*ServiceContext, error) {
	gdb, err := db.GetWorkSpaceGormDb(workSpace.WorkSpace)
	if err != nil {
		return nil, err
	}
	return &ServiceContext{
		WorkSpace: workSpace,
		Context:   c,
		DB:        gdb,
	}, nil
}

func ListDbData(dbQuery *DbQuery, sc *ServiceContext) (*ListResult, error) {

	form := dbQuery.Form
	if form.Page < 1 {
		form.Page = DEFAULT_PAGE
	}
	if form.PageSize < 1 {
		form.PageSize = DEFAULT_PAGE_SIZE
	}

	if form.PageSize > MAX_PAGE_SIZE {
		form.PageSize = MAX_PAGE_SIZE
	}

	Db := sc.DB

	//find collection
	collection, err := db.GetDbApiConfigCollectionByName(sc.DB, dbQuery.Collection)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("find collection error,%v", err))
	}
	var apiConfig db.DbApiConfig
	err = Db.Where("collection_id = ?", collection.ID).Where("name = ?", dbQuery.Table).First(&apiConfig).Error
	if err != nil {
		return nil, errors.New(fmt.Sprintf("find apiConfig error,%v", err))
	}

	baseQuery := Db.Where("db_api_config_id = ?", apiConfig.ID).Order("row_num")

	dbExecutor, err := buildQueryDb(baseQuery, &apiConfig, form.Query)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("build query db error,%v", err))
	}

	listResult := ListResult{}

	var total int64
	err = dbExecutor.Session(&gorm.Session{}).Model(&db.DbData{}).Group("row_num").Count(&total).Error
	if err != nil {
		return nil, errors.New(fmt.Sprintf("query dbDatas count error,%v", err))
	}
	listResult.Total = total

	subquery := dbExecutor.Session(&gorm.Session{}).Model(&db.DbData{}).Select("row_num").Group("row_num").Limit(form.PageSize).Offset((form.Page - 1) * form.PageSize)
	dbDatas := []*db.DbData{}
	err = dbExecutor.Where(" row_num in (?) ", subquery).Find(&dbDatas).Error
	if err != nil {
		return nil, errors.New(fmt.Sprintf("query dbDatas error,%v", err))
	}

	resultDbData, err := buildResultData(dbDatas, apiConfig)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("build result db data error,%v", err))
	}
	listResult.Data = *resultDbData
	return &listResult, nil
}
