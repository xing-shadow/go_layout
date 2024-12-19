package httpServe

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"layout/global"
	"net/http"
	"time"
)

var srv = new(HttpServer)

type Option struct {
	Port      int
	ReadTime  int //单位秒
	WriteTime int //单位秒
}

type HttpServer struct {
	opt   *Option
	Route *gin.Engine
	*http.Server
}

func Init(o *Option) error {
	//
	srv.opt = o
	//
	srv.Route = gin.New()
	srv.SetRoute()
	//
	srv.Server = &http.Server{
		Addr:         fmt.Sprintf(":%d", srv.opt.Port),
		Handler:      srv.Route,
		ReadTimeout:  time.Duration(srv.opt.ReadTime) * time.Second,
		WriteTimeout: time.Duration(srv.opt.WriteTime) * time.Second,
	}
	return nil
}

func Start() error {
	global.GetLogger().Infof("开启http服务 %s:%d", "", srv.opt.Port)
	return srv.ListenAndServe()
}

func Stop() error {
	return srv.Shutdown(context.Background())
}
