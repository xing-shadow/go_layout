package httpServe

import (
	"github.com/gin-gonic/gin"
	"layout/pkg/Context"
	"net/http"
)

type HandlerFunc func(ctx *Context.Context)

func (pThis *HttpServer) SetRoute() {
	//中间件添加
	pThis.Route.Use(gin.Recovery(), Cors())
	//
	srv.Route.NoRoute(func(ctx *gin.Context) {
		ctx.Data(http.StatusNotFound, "text/plain", []byte(http.StatusText(http.StatusNotFound)))
		ctx.Abort()
	})
	srv.Route.NoMethod(func(ctx *gin.Context) {
		ctx.Data(http.StatusMethodNotAllowed, "text/plain", []byte(http.StatusText(http.StatusMethodNotAllowed)))
		ctx.Abort()
	})

	//业务handle
	apiV1 := pThis.Route.Group("/api/v1")
	apiV1.Use(LoggerMiddleware())

}
