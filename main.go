package main

import (
	"github.com/basegin/base"
	"github.com/basegin/db"
	"github.com/basegin/router"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化 gin 框架
	r := gin.Default()

	//启动配置
	base.Bootstrap(r)
	db.InitMysql()
	// 添加路由
	router.Http(r)
	// 健康检查curl -I http://host:port/health
	r.HEAD("/health", func(c *gin.Context) {
		c.Status(200)
	})

	endless.ListenAndServe(":8080", r)
}
