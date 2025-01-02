package service

import (
	"easywrite-service/model"
	"github.com/gin-gonic/gin"
)

func HttpParameterError(c *gin.Context) {
	c.JSON(200, model.JsonResponse[any]{
		Code: 403,
		Msg:  "参数错误",
		Data: nil,
	})
}

func HttpServerInternalError(c *gin.Context) {
	c.JSON(200, model.JsonResponse[any]{
		Code: 500,
		Msg:  "Http服务器内部错误",
		Data: nil,
	})
}

func HttpOK(c *gin.Context) {
	c.JSON(200, model.JsonResponse[any]{
		Code: 200,
		Msg:  "ok",
		Data: nil,
	})
}
