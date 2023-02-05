package controller

import (
	"github.com/basegin/base/log"
	models "github.com/basegin/model"
	"github.com/gin-gonic/gin"
)

type ConfigParam struct {
	Key string `form:"key" binding:"required"`
}

func GetConfig(c *gin.Context) {
	var input ConfigParam
	if err := c.ShouldBind(&input); err != nil {
		log.ErrorLogger(c, err.Error())
		return
	}

	v, err := models.GetConfigValue(c, input.Key)
	if err != nil {
		log.RenderJsonFail(c, err)
	}

	log.RenderJsonSucc(c, gin.H{"value": v})
}
