package base

import (
	"github.com/basegin/base/log"
	"github.com/basegin/config"
	"github.com/gin-gonic/gin"
)

func Bootstrap(router *gin.Engine) {
	//日志配置 json格式，默认标准输出
	log.InitLog()

	// trace init
	log.InitTracer()
	err := config.LoadConfig("/Users/liyanan/go/src/github.com/basegin/config/config.yaml")
	if err != nil {
		log.ErrorLogger(nil, err)
		return
	}
	// 中间件
	router.Use(log.Tracer())
	router.Use(log.LoggerMiddleware())
	router.Use(gin.Recovery())

	//性能分析工具
	log.Register(router)
}
