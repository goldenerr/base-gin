package log

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"net/http"
)

type TRenderJson struct {
	ReturnCode    int                    `json:"returnCode"`
	ReturnMsg     string                 `json:"returnMsg"`
	ReturnUserMsg string                 `json:"returnUserMsg"`
	Data          map[string]interface{} `json:"data"`
}

func RenderJsonSucc(ctx *gin.Context, data map[string]interface{}) {
	renderjson := TRenderJson{0, "succ", "成功", data}
	ctx.JSON(http.StatusOK, renderjson)
	ctx.Set("render", renderjson)
	return
}

func RenderJsonFail(ctx *gin.Context, err error) {
	var renderjson TRenderJson

	switch errors.Cause(err).(type) {

	case Error:
		renderjson.ReturnMsg = errors.Cause(err).(Error).ErrorMsg
		renderjson.ReturnUserMsg = errors.Cause(err).(Error).ErrorUserMsg
		renderjson.ReturnCode = errors.Cause(err).(Error).ErrorCode

	case validator.ValidationErrors:
		renderjson.ReturnMsg = "ParamError"
		renderjson.ReturnUserMsg = "参数错误"
		renderjson.ReturnCode = 4000000

		var validateData []interface{}
		for _, err := range errors.Cause(err).(validator.ValidationErrors) {
			validateData = append(validateData, map[string]interface{}{"message": "params check error", "tag": err.Tag, "field": err.Field})
		}

		renderjson.Data = map[string]interface{}{"validateData": validateData}

	default:
		renderjson.ReturnMsg = errors.Cause(err).Error()
		renderjson.ReturnUserMsg = ""
		renderjson.ReturnCode = -1
	}
	ctx.JSON(http.StatusOK, renderjson)
	ctx.Set("render", renderjson)

	//打印错误栈
	StackLogger(ctx, err)
	return
}
