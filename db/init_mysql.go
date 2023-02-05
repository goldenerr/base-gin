package db

import (
	"fmt"
	"github.com/basegin/base/jinzhu/gorm"
	"github.com/basegin/base/log"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"time"
)

const SqlTraceContextKey = "traceContext"

type MysqlConf struct {
	User            string
	Password        string
	Addr            string
	Port            int
	DataBase        string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifeTime time.Duration
	LogMode         bool
}

// InitMysqlClient 一定要在bootstrap之后
func InitMysqlClient(conf MysqlConf) (client *gorm.DB, err error) {
	client, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Asia%%2FShanghai",
		conf.User,
		conf.Password,
		conf.Addr,
		conf.DataBase))

	if err != nil {
		return client, err
	}

	client.DB().SetMaxIdleConns(conf.MaxIdleConns)
	client.DB().SetMaxOpenConns(conf.MaxOpenConns)
	client.DB().SetConnMaxLifetime(conf.ConnMaxLifeTime)
	client.LogMode(conf.LogMode)

	// register tracer callback
	setCallback(client, "create")
	setCallback(client, "delete")
	setCallback(client, "update")
	setCallback(client, "query")
	setCallback(client, "row_query")

	return client, nil
}

// InitMysqlClient 一定要在bootstrap之后
func InitMysqlClient1(conf MysqlConf) (client *gorm.DB, err error) {
	client, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Asia%%2FShanghai",
		conf.User,
		conf.Password,
		conf.Addr,
		conf.Port,
		conf.DataBase))

	if err != nil {
		return client, err
	}

	client.DB().SetMaxIdleConns(conf.MaxIdleConns)
	client.DB().SetMaxOpenConns(conf.MaxOpenConns)
	client.DB().SetConnMaxLifetime(conf.ConnMaxLifeTime)
	client.LogMode(conf.LogMode)

	// register tracer callback
	setCallback(client, "create")
	setCallback(client, "delete")
	setCallback(client, "update")
	setCallback(client, "query")
	setCallback(client, "row_query")

	return client, nil
}
func setCallback(client *gorm.DB, callbackName string) {
	beforeName := fmt.Sprintf("tracer:%v_before", callbackName)
	afterName := fmt.Sprintf("tracer:%v_after", callbackName)
	gormCallbackName := fmt.Sprintf("gorm:%v", callbackName)
	switch callbackName {
	case "create":
		client.Callback().Create().Before(gormCallbackName).Register(beforeName, func(scope *gorm.Scope) {
			tracerBefore(scope, callbackName)
		})
		client.Callback().Create().After(gormCallbackName).Register(afterName, func(scope *gorm.Scope) {
			tracerAfter(scope, callbackName)
		})
	case "query":
		client.Callback().Query().Before(gormCallbackName).Register(beforeName, func(scope *gorm.Scope) {
			tracerBefore(scope, callbackName)
		})
		client.Callback().Query().After(gormCallbackName).Register(afterName, func(scope *gorm.Scope) {
			tracerAfter(scope, callbackName)
		})
	case "update":
		client.Callback().Update().Before(gormCallbackName).Register(beforeName, func(scope *gorm.Scope) {
			tracerBefore(scope, callbackName)
		})
		client.Callback().Update().After(gormCallbackName).Register(afterName, func(scope *gorm.Scope) {
			tracerAfter(scope, callbackName)
		})
	case "delete":
		client.Callback().Delete().Before(gormCallbackName).Register(beforeName, func(scope *gorm.Scope) {
			tracerBefore(scope, callbackName)
		})
		client.Callback().Delete().After(gormCallbackName).Register(afterName, func(scope *gorm.Scope) {
			tracerAfter(scope, callbackName)
		})
	case "row_query":
		client.Callback().RowQuery().Before(gormCallbackName).Register(beforeName, func(scope *gorm.Scope) {
			tracerBefore(scope, callbackName)
		})
		client.Callback().RowQuery().After(gormCallbackName).Register(afterName, func(scope *gorm.Scope) {
			tracerAfter(scope, callbackName)
		})
	}
}

func tracerBefore(scope *gorm.Scope, callbackName string) {
	ctx, ok := scope.Search.GetCtx().(*gin.Context)
	if !ok {
		fmt.Println("11111111111111111")
		return
	}
	var parentSpanContext opentracing.SpanContext
	if ctx != nil {
		tempSpan, exist := ctx.Get(log.TraceContextKey)
		if exist {
			parentSpanContext = tempSpan.(opentracing.Span).Context()
		}
	}

	span := log.TracerStartSpan(parentSpanContext, "DB_"+callbackName, map[string]interface{}{
		ext.SpanKindRPCClient.Key: ext.SpanKindRPCClient,
		string(ext.Component):     "gorm",
		string(ext.DBInstance):    scope.InstanceID(),
		string(ext.DBType):        "mysql",
	})
	scope.Set(SqlTraceContextKey, span)
}

func tracerAfter(scope *gorm.Scope, callbackName string) {
	tempSpan, exist := scope.Get(SqlTraceContextKey)
	if !exist {
		return
	}
	span := tempSpan.(opentracing.Span)
	span.SetTag(string(ext.DBStatement), scope.SQL)
	span.Finish()
}
