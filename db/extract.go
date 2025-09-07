package db

import (
	"db-rest/config"
	"db-rest/env"
	"db-rest/util"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"os"
	"path"
	"sync"
)

type extractJobContext struct {
	id          string
	dbConnector *DbConnector
	db          *gorm.DB
	waitGroup   *sync.WaitGroup
	mutex       *sync.Mutex
	jobs        []*extractJob
}

func (ctx *extractJobContext) Close() {
	if ctx.dbConnector != nil {
		ctx.dbConnector.Close()
	}
}

type extractJob struct {
	dbName     string
	collection *config.ApiConfigCollection
	apiConfig  *config.ApiConfig
}

func newExtractJobContext(id string, dbConnector *DbConnector, db *gorm.DB, waitGroup *sync.WaitGroup, mutex *sync.Mutex, jobs []*extractJob) *extractJobContext {
	return &extractJobContext{id, dbConnector, db, waitGroup, mutex, jobs}
}

func extractJobFn(ctx *extractJobContext) {
	defer ctx.waitGroup.Done()
	util.Log("start extracting job %s", ctx.id)

	usePagination := config.GetEnvValue[bool](config.VIPER_KEY_USE_PAGINATION_EXTRACT)
	workspace := env.GetEnvVar[string](env.WORKSPACE)
	batchSize := config.GetEnvValue[int](config.VIPER_KEY_SAVE_BATCH_SIZE)
	isStrict := config.GetEnvValue[bool](config.VIPER_KEY_MODE_STRICT)

	for _, job := range ctx.jobs {

		if job.collection != nil {

			var collectionId int = -1

			ctx.mutex.Lock()
			c, err := CountDbApiConfigCollectionByName(ctx.db, job.collection.Name)
			if err != nil {
				util.LogWarn("get db [%s] api config collection [%s] occur error: %v", job.dbName, job.collection.Name, err)
			}
			if c <= 0 {
				util.Log("save collection [%s]", job.collection.Name)

				dbCollection := &DbApiConfigCollection{
					Name:        job.collection.Name,
					Description: job.collection.Description,
				}

				err := SaveDbApiConfigCollection(ctx.db, dbCollection)

				if err != nil {
					util.LogWarn("save db [%s] api config collection [%s] occur error: %v", job.dbName, job.collection.Name, err)
				}
				collectionId = int(dbCollection.ID)
			} else {
				dbCollection, err := GetDbApiConfigCollectionByName(ctx.db, job.collection.Name)
				if err != nil {
					util.LogWarn("get db [%s] api config collection [%s] occur error: %v", job.dbName, job.collection.Name, err)
				}
				collectionId = int(dbCollection.ID)
			}
			ctx.mutex.Unlock()

			if collectionId > 0 {

				//save dbApiConfig
				dbApiConfig := &DbApiConfig{
					Name:         job.apiConfig.Name,
					Description:  job.apiConfig.Description,
					Sql:          job.apiConfig.Sql,
					Columns:      util.JoinStr(",", job.apiConfig.Columns...),
					DbName:       job.dbName,
					CollectionId: uint(collectionId),
				}
				err := ctx.db.Create(dbApiConfig).Error

				if err != nil {
					util.LogWarn("save db [%s] api config [%s] occur error: %v", job.dbName, job.apiConfig.Name, err)
				}

				//load data,transform
				var transformFn func(data *DataRow)

				if job.apiConfig.TransformJsPath != "" {

					jsPath := path.Join(workspace, job.apiConfig.TransformJsPath)

					if util.FileExists(jsPath) {

						b, err := os.ReadFile(jsPath)
						if err != nil {
							util.LogWarn("read transform js path [%s] occur error: %v", jsPath, err)
						} else {
							script := string(b)
							_, fn, err := GetVm(script)
							if err != nil {
								util.LogWarn("parse transform js path [%s] occur error: %v", jsPath, err)
							} else {
								transformFn = fn
							}

						}
					} else {
						util.LogWarn("transform js path [%s] not exists", jsPath)
					}
				}

				processFn := func(rows []*DataRow) error {

					if transformFn != nil {
						util.Log("transform datas")
						for _, row := range rows {
							transformFn(row)
						}
					}

					//save data
					dbDatas := []*DbData{}
					for idx, row := range rows {

						rowToDbDatas := dataRowToDbDatas(row)

						for _, dbData := range rowToDbDatas {
							dbData.DbApiConfigId = dbApiConfig.ID
							dbData.Table = job.apiConfig.Name
						}

						dbDatas = append(dbDatas, rowToDbDatas...)

						if idx%batchSize == 0 {
							ctx.mutex.Lock()
							err := ctx.db.CreateInBatches(dbDatas, batchSize).Error
							ctx.mutex.Unlock()
							if err != nil {
								if isStrict {
									return err
								}
							}

							dbDatas = dbDatas[:0]
						}

					}
					if len(dbDatas) > 0 {
						ctx.mutex.Lock()
						err := ctx.db.CreateInBatches(dbDatas, batchSize).Error
						ctx.mutex.Unlock()
						if err != nil {

							if isStrict {
								return err
							}

						}

					}
					return nil
				}

				var loader Loader
				if usePagination {
					loader = NewPaginationDataLoader(ctx.dbConnector, ctx.db, job.apiConfig, processFn)
				} else {
					loader = NewAllDataLoader(ctx.dbConnector, ctx.db, job.apiConfig, processFn)
				}

				if loader != nil {
					extractC, err := loader.Load()
					if err != nil {
						util.LogWarn("load db [%s] api config [%s] occur error: %v", job.dbName, job.apiConfig.Name, err)
					} else {
						util.LogSuccess("load db %s api config [%s]  data rows %d", job.dbName, job.apiConfig.Name, extractC)
					}

				}
			}

		}

	}

}

func newExtractJob(collection *config.ApiConfigCollection, apiCofig *config.ApiConfig) *extractJob {
	return &extractJob{
		dbName:     apiCofig.DbName,
		collection: collection,
		apiConfig:  apiCofig,
	}
}
func Extract(cfg *config.WorkSpaceParsedConfig) error {
	initDbDriver()
	if cfg.Config == nil {
		return errors.New("config is nil")
	}

	workspaceDir := cfg.WorkSpace

	gdb, err := GetWorkSpaceGormDb(workspaceDir)

	if err != nil {
		return err
	}

	sqlDb, err := gdb.DB()

	if err != nil {
		return err
	}

	defer sqlDb.Close()

	maxConnection := config.GetEnvValue[int](config.VIPER_KEY_MAX_CONNECTION)

	if maxConnection <= 0 {
		maxConnection = config.MAX_CONNECTION
	}
	sqlDb.SetMaxIdleConns(maxConnection / 3 * 2)
	sqlDb.SetMaxOpenConns(maxConnection)

	err = CleanAllWorkSpaceData(gdb)
	if err != nil {
		return err
	}

	maxDatasource := config.GetEnvValue[uint](config.VIPER_KEY_MAX_DATASOURCE)

	if maxDatasource > 0 && len(cfg.Config.BasicConfig.DbList) > int(maxDatasource) {
		return errors.New(fmt.Sprintf("max datasource exceeded,max %d ,now is %d ", maxDatasource, len(cfg.Config.BasicConfig.DbList)))
	}

	extractJobContexts := []*extractJobContext{}
	extractJobs := []*extractJob{}
	mutex := sync.Mutex{}
	waitGroup := sync.WaitGroup{}

	for _, collection := range cfg.ApiCollection {
		for _, apiConfig := range collection.ApiList {
			extractJob := newExtractJob(&collection, &apiConfig)
			extractJobs = append(extractJobs, extractJob)

		}
	}

	for idx, dbConfig := range cfg.Config.BasicConfig.DbList {
		dbConnector, err := getDbConnector(&dbConfig)
		if err != nil {
			util.LogWarn("get db [%s] connector occur error: %v", dbConfig.Name, err)
			continue
		}
		if dbConnector != nil {
			jobs := []*extractJob{}

			for _, job := range extractJobs {
				if util.StrEquals(job.dbName, dbConfig.Name) {
					jobs = append(jobs, job)
				}
			}

			jobId := util.FormatStr("extract_job_%d", idx+1)
			extractJobContexts = append(extractJobContexts, newExtractJobContext(jobId, dbConnector, gdb, &waitGroup, &mutex, jobs))
		}

	}

	for _, jobCtx := range extractJobContexts {
		if jobCtx.dbConnector != nil && len(jobCtx.jobs) > 0 {
			waitGroup.Add(1)
			go extractJobFn(jobCtx)
		}
	}

	waitGroup.Wait()
	for _, jobCtx := range extractJobContexts {
		jobCtx.Close()
	}

	return nil
}

func dataRowToDbDatas(row *DataRow) []*DbData {
	dbDataRows := []*DbData{}
	rowNum := row.RowNum
	for _, column := range row.Columns {
		dbData := &DbData{
			Column:    column.Name,
			Value:     column.Value,
			RowNum:    rowNum,
			ValueType: DB_DATA_TYPE_STRING,
		}
		dbDataRows = append(dbDataRows, dbData)
	}
	return dbDataRows
}
