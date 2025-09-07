package api

import (
	"db-rest/service"
	"db-rest/util"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type userForm struct {
	Name string `form:"name"`
	Pass string `form:"pass"`
}

func MountRouters(engine *gin.Engine, sc *service.ServiceContext) {

	authorizeHandler := func(ctx *gin.Context) {

		userForm := userForm{}

		switch ctx.Request.Method {
		case "POST":
			ctx.ShouldBind(&userForm)
		default:
			ctx.BindQuery(&userForm)
		}

		authorized := false

		if userForm.Name != "" && userForm.Pass != "" {
			auth := sc.WorkSpace.Config.BasicConfig.Auth
			if util.StrEquals(userForm.Name, auth.User) && util.StrEquals(userForm.Pass, auth.Pass) {
				authorized = true
			}
		}

		if authorized {
			token, err := service.GenToken(userForm.Name)

			if err != nil {
				ctx.AbortWithStatusJSON(500, NewResult(CodeError, err.Error()))
			} else {
				ctx.JSON(200, NewResultWithData(CodeSuccess, "success", token))
			}
		} else {
			ctx.AbortWithStatusJSON(401, NewResult(CodeError, "incorrect user or pass"))
		}
	}

	engine.GET(pathTokenGet, authorizeHandler)
	engine.POST(pathTokenGet, authorizeHandler)

	engine.GET(pathWorkSpaceList, func(context *gin.Context) {

		result, err := service.ListWorkspace(sc)

		if err != nil {
			context.JSON(500, NewResult(CodeError, err.Error()))
		} else {
			context.JSON(200, NewResultWithData(CodeSuccess, "success", &result))
		}
	})

	if err := mountDbRest(engine, sc); err != nil {
		util.LogError("mountDbRest error,%v", err)
		panic(err)
	}

}

func mountDbRest(engine *gin.Engine, sc *service.ServiceContext) error {

	basicConfig := sc.WorkSpace.Config.BasicConfig

	globalRoutePath := basicConfig.Prefix

	if globalRoutePath == "" {
		globalRoutePath = "/"
	}

	for _, path := range white_list {
		if util.StrEquals(path, globalRoutePath) {
			return errors.New(fmt.Sprintf("path %s can't be used", globalRoutePath))
		}
	}

	util.Log("global route path: %s", globalRoutePath)
	g := engine.Group(globalRoutePath)

	g.GET("/:collection/:table/list", func(context *gin.Context) {
		context.JSON(404, NewResult(CodeError, "not support GET method"))
	})

	g.POST("/:collection/:table/list", func(context *gin.Context) {

		var dbQuery = service.DbQuery{}

		var queryForm = service.DbQueryForm{}

		collection := context.Param("collection")
		table := context.Param("table")

		dbQuery.Collection = collection
		dbQuery.Table = table

		err := context.ShouldBindJSON(&queryForm)

		if err != nil {
			context.JSON(400, NewResult(CodeError, fmt.Sprintf("queryForm error;%v,confirm json body correct", err)))
			return
		}

		dbQuery.Form = &queryForm

		result, err := service.ListDbData(&dbQuery, sc)

		if err != nil {
			context.JSON(500, NewResult(CodeError, err.Error()))
		} else {
			context.JSON(200, NewResultWithData(CodeSuccess, "success", &result))
		}

	})

	//excel
	g.GET("/:collection/:table/excel", func(context *gin.Context) {
		context.JSON(404, NewResult(CodeError, "not support GET method"))
	})

	g.POST("/:collection/:table/excel", func(context *gin.Context) {

		var dbQuery = service.DbQuery{}

		var queryForm = service.DbQueryForm{}

		collection := context.Param("collection")
		table := context.Param("table")

		dbQuery.Collection = collection
		dbQuery.Table = table

		err := context.ShouldBindJSON(&queryForm)

		if err != nil {
			context.JSON(400, NewResult(CodeError, fmt.Sprintf("queryForm error;%v,confirm json body correct", err)))
			return
		}

		dbQuery.Form = &queryForm

		result, err := service.ListDbData(&dbQuery, sc)

		if err != nil {
			context.JSON(500, NewResult(CodeError, err.Error()))
		} else {

			excelFile := excelize.NewFile()
			defer excelFile.Close()

			sheetName := "export"

			idx, err := excelFile.NewSheet(sheetName)
			if err != nil {
				context.JSON(500, NewResult(CodeError, err.Error()))
			}

			resultData := result.Data

			head := []string{}

			for _, column := range resultData.DataStructs {
				head = append(head, column.Name)
			}
			excelFile.SetSheetRow(sheetName, "A1", &head)
			for index, row := range resultData.DataRows {
				rowData := []string{}
				for _, column := range resultData.DataStructs {
					rowData = append(rowData, row.RowData[column.Name])
				}
				excelFile.SetSheetRow(sheetName, fmt.Sprintf("A%d", index+2), &rowData)
			}

			excelFile.SetActiveSheet(idx)

			context.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
			context.Header("Content-Disposition", "attachment; filename=export.xlsx")

			excelFile.Write(context.Writer)

		}

	})

	return nil
}
