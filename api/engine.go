package api

import (
	"db-rest/config"
	"db-rest/service"
	"db-rest/util"
	"github.com/gin-gonic/gin"
	"strings"
)

type RestEngine struct {
	engine *gin.Engine
	Sc     *service.ServiceContext
}

func New(sc *service.ServiceContext) (RestEngine, error) {

	isDebugMode := config.GetEnvValue[bool](config.VIPER_KEY_ENGINE_DEBUG_MODE)
	if isDebugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.New(func(eg *gin.Engine) {

		UseMiddlewares(eg, sc)
		MountRouters(eg, sc)

	})
	return RestEngine{
		engine: engine,
		Sc:     sc,
	}, nil
}
func (re *RestEngine) Run() error {

	bind := config.GetEnvValue[string](config.VIPER_KEY_ENGINE_BIND)

	if strings.TrimSpace(bind) == "" {
		bind = re.Sc.WorkSpace.Config.BasicConfig.Bind
	}

	if strings.TrimSpace(bind) == "" {
		bind = config.BIND
	}

	bind = strings.TrimSpace(bind)
	util.Log("starting rest engine on %s", bind)
	return re.engine.Run(bind)
}
