package httpServe

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
)

type Response struct {
	Code    int         `json:"err_code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Json(ctx *gin.Context, code int, msg string, data interface{}) {
	result := new(Response)
	result.Code = code
	result.Message = msg
	if data != nil {
		result.Data = data
	}
	if code != CodeOK {
		ctx.Error(errors.New(msg))
	}
	ctx.JSON(http.StatusOK, result)
}
