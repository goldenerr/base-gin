package log

import (
	"encoding/json"
	"fmt"
	"github.com/rifflock/lfshook"
	"math"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/file-rotatelogs"
	mylog "github.com/sirupsen/logrus"
)

//var logger *mylog.Logger

func InitLog() {
	//logger = mylog.New()
	// 创建一个文件夹，用来存储分隔后的日志
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		os.Mkdir("logs", os.ModePerm)
	}

	logFile := "./logs/app-%Y%m%d.log"
	logWriter, err := rotatelogs.New(
		logFile,
		//rotatelogs.WithLinkName("./logs/app.log"),
		rotatelogs.WithMaxAge(time.Duration(30*24)*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(24)*time.Hour),
	)
	if err != nil {
		mylog.Fatalf("failed to create log writer: %v", err)
	}
	lfHook := lfshook.NewHook(lfshook.WriterMap{
		mylog.DebugLevel: logWriter,
		mylog.InfoLevel:  logWriter,
		mylog.WarnLevel:  logWriter,
		mylog.ErrorLevel: logWriter,
		mylog.FatalLevel: logWriter,
		mylog.PanicLevel: logWriter,
	}, &mylog.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	mylog.AddHook(lfHook)

	// 设置日志输出到rotatelogs
	mylog.SetOutput(os.Stdout)
	//logger.SetOutput(logWriter)
	//logger.SetOutput(io.MultiWriter(os.Stdout, logWriter))

	// todo 设置日志级别
	mylog.SetLevel(mylog.InfoLevel)

	// 将函数名和行数放在日志里
	//logger.SetReportCaller(true)
}

// 记录日志
func LoggerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

		ctx.Next()

		if GetNoLogFlag(ctx) == true {
			return
		}

		stop := time.Since(start)
		latency := int(math.Ceil(float64(stop.Nanoseconds()) / 1000000.0))

		dataLength := ctx.Writer.Size()
		if dataLength < 0 {
			dataLength = 0
		}

		_ = ctx.Request.ParseForm()
		requestBody := ctx.Request.PostForm.Encode()
		if len(requestBody) > 10240 {
			requestBody = "too big not show"
		}

		render, exist := ctx.Get("render")
		var retCode int
		if exist == true {
			retCode = render.(TRenderJson).ReturnCode
		} else {
			retCode = 0
		}

		response, _ := json.Marshal(render)
		if len(string(response)) > 10240 {
			response = []byte("too big not show")
		}

		mylog.NewEntry(mylog.StandardLogger()).WithFields(mylog.Fields{
			"requestId": GetRequestId(ctx),
			"handle":    ctx.HandlerName(),
			"type":      "http",
			"retCode":   retCode,
			"latency":   latency, // time to process

			"requestBody": requestBody,
			"requestPath": ctx.Request.URL.RequestURI(),
			"httpCode":    ctx.Writer.Status(),
			"clientIp":    ctx.ClientIP(),
			"method":      ctx.Request.Method,
			"referer":     ctx.Request.Referer(),
			"dataLength":  dataLength,
			"userAgent":   ctx.Request.UserAgent(),
		}).Info(string(response))
	}
}

const NOLOGFlAG string = "NOLOGFlAG"

//func getLogLevel(loglevel string) (mylog.Level, bool) {
//	switch loglevel {
//	case "debug":
//		return mylog.DebugLevel, true
//	case "info":
//		return mylog.InfoLevel, true
//	case "warn":
//		return mylog.WarnLevel, true
//	case "error":
//		return mylog.ErrorLevel, true
//	case "fatal":
//		return mylog.FatalLevel, true
//	case "panic":
//		return mylog.PanicLevel, true
//	default:
//		return 0, false
//	}
//}

func Logger(ctx *gin.Context) *mylog.Entry {
	entry := mylog.NewEntry(mylog.StandardLogger()).WithFields(mylog.Fields{
		"requestId": GetRequestId(ctx),
	})

	return entry
}

func DebugLogger(ctx *gin.Context, value ...interface{}) {
	if GetNoLogFlag(ctx) == true {
		return
	}
	if ctx != nil {
		Logger(ctx).Debug(value...)
	} else {
		mylog.Debug(value...)
	}
}

func DebugfLogger(ctx *gin.Context, format string, value ...interface{}) {
	if GetNoLogFlag(ctx) == true {
		return
	}
	if ctx != nil {
		Logger(ctx).Debugf(format, value...)
	} else {
		mylog.Debugf(format, value...)
	}
}

func InfoLogger(ctx *gin.Context, value ...interface{}) {
	if GetNoLogFlag(ctx) == true {
		return
	}
	if ctx != nil {
		Logger(ctx).Info(value...)
	} else {
		mylog.Info(value...)
	}
}

func InfofLogger(ctx *gin.Context, format string, value ...interface{}) {
	if GetNoLogFlag(ctx) == true {
		return
	}
	if ctx != nil {
		Logger(ctx).Infof(format, value...)
	} else {
		mylog.Infof(format, value...)
	}
}

func WarnLogger(ctx *gin.Context, value ...interface{}) {
	if GetNoLogFlag(ctx) == true {
		return
	}
	if ctx != nil {
		Logger(ctx).Warn(value...)
	} else {
		mylog.Warn(value...)
	}
}

func WarnfLogger(ctx *gin.Context, format string, value ...interface{}) {
	if GetNoLogFlag(ctx) == true {
		return
	}
	if ctx != nil {
		Logger(ctx).Warnf(format, value...)
	} else {
		mylog.Warnf(format, value...)
	}
}

func ErrorLogger(ctx *gin.Context, value ...interface{}) {
	if ctx != nil {
		Logger(ctx).Error(value...)
	} else {
		mylog.Error(value...)
	}
}

func ErrorfLogger(ctx *gin.Context, format string, value ...interface{}) {
	if ctx != nil {
		Logger(ctx).Errorf(format, value...)
	} else {
		mylog.Errorf(format, value...)
	}
}

func PanicLogger(ctx *gin.Context, value ...interface{}) {
	if ctx != nil {
		Logger(ctx).Panic(value...)
	} else {
		mylog.Panic(value...)
	}
}

func PanicfLogger(ctx *gin.Context, format string, value ...interface{}) {
	if ctx != nil {
		Logger(ctx).Panicf(format, value...)
	} else {
		mylog.Panicf(format, value...)
	}
}

func SetNoLogFlag(ctx *gin.Context) {
	ctx.Set(NOLOGFlAG, true)
}

func GetNoLogFlag(ctx *gin.Context) bool {
	if ctx == nil {
		return false
	}
	flag, ok := ctx.Get(NOLOGFlAG)
	if ok == true && flag == true {
		return true
	} else {
		return false
	}
}

func StackLogger(ctx *gin.Context, err error) {
	if !strings.Contains(fmt.Sprintf("%+v", err), "\n") {
		return
	}

	var info []byte
	if ctx != nil {
		info, _ = json.Marshal(map[string]interface{}{"time": time.Now().Format("2006-01-02 15:04:05"), "level": "error", "module": "errorstack", "requestId": GetRequestId(ctx)})
	} else {
		info, _ = json.Marshal(map[string]interface{}{"time": time.Now().Format("2006-01-02 15:04:05"), "level": "error", "module": "errorstack"})
	}

	fmt.Printf("%s\n-------------------stack-start-------------------\n%+v\n-------------------stack-end-------------------\n", string(info), err)
}

func GetRequestId(ctx *gin.Context) string {
	if ctx == nil {
		return ""
	}
	requestId, exist := ctx.Get("requestId")
	if exist {
		return requestId.(string)
	}
	return ""
}

//func GetLogger() *mylog.Logger {
//	return logger
//}
