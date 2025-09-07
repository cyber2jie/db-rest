package api

import (
	"db-rest/service"
	"db-rest/util"
	"github.com/gin-gonic/gin"
	"strings"
)

func authenrize(ctx *gin.Context) {
	needAuth := true

	uri := ctx.Request.URL.Path

	for _, path := range white_list {
		if util.StrEquals(path, uri) {
			needAuth = false
			break
		}
	}
	bear := ctx.GetHeader("Authorization")
	if needAuth && (bear == "" || !service.IsValidToken(strings.Replace(bear, "Bearer ", "", -1))) {
		ctx.AbortWithStatusJSON(401, NewResult(CodeError, "unauthorize"))
	} else {
		ctx.Next()
	}
}

func UseMiddlewares(engine *gin.Engine, sc *service.ServiceContext) {

	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	if sc.WorkSpace.Config.BasicConfig.Auth.Enable {
		engine.Use(authenrize)
	}

}
