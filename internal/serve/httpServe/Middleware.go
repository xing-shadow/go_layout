package httpServe

import (
	"github.com/gin-gonic/gin"
	"layout/pkg/Context"
	"net/http"
	"reflect"
	"time"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := Context.GetContextFromGin(c)
		now := time.Now()
		req := c.Request
		path := req.URL.Path
		var quota float64
		if deadline, ok := c.Request.Context().Deadline(); ok {
			quota = time.Until(deadline).Seconds()
		}
		c.Next()
		dt := time.Since(now)
		err := c.Errors.Last()
		errMsg := ""
		isSlow := dt >= (time.Second)
		isWarn, isError := false, false
		if !reflect.ValueOf(&err).Elem().IsNil() {
			errMsg = err.Error()
			isError = true
		} else {
			if isSlow {
				isWarn = true
			}
		}
		if isError {
			ctx.Errorf("[HttpServer] method=%s | path=%s | ts=%fs | timeout_quota=%fs | err=%s ",
				req.Method, path, dt.Seconds(), quota, errMsg)
		} else if isWarn {
			ctx.Warnf("[HttpServer] method=%s | path=%s | ts=%fs | timeout_quota=%fs | err=%s ",
				req.Method, path, dt.Seconds(), quota, errMsg)
		} else {
			ctx.Warnf("[HttpServer] method=%s | path=%s | ts=%fs | timeout_quota=%fs | err=%s ",
				req.Method, path, dt.Seconds(), quota, errMsg)
		}
	}
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		//Origin
		origin := c.Request.Header.Get("Origin")
		if origin == "" {
			origin = "*"
		}
		c.Header("Access-Control-Allow-Origin", origin)

		c.Header("Access-Control-Allow-Headers", "Origin,Authorization,Access-Control-Allow-Origin,Access-Control-Allow-Headers,Content-Type,Signe,Sign")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length,Access-Control-Allow-Origin,Access-Control-Allow-Headers,Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		// 放行所有OPTIONS方法，因为有的模板是要请求两次的
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}

		// 处理请求
		c.Next()
	}
}
