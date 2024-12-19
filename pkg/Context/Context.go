package Context

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"layout/global"
	"layout/pkg/utils"
)

const (
	injectGinContextKey = "custom-inject-Ctx-info"
)

type Context struct {
	*zap.SugaredLogger
	ctx context.Context
}

func GetContextFromGin(ginCtx *gin.Context) *Context {
	val, exist := ginCtx.Get(injectGinContextKey)
	if exist {
		if ref, ok := val.(*Context); ok {
			return ref
		}
	}
	ctx := Context{
		SugaredLogger: global.GetLogger().With("req_id", utils.GetUUid()),
		ctx:           context.Background(),
	}
	ginCtx.Set(injectGinContextKey, &ctx)
	return &ctx
}
