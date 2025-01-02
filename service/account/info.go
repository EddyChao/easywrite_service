package account

import (
	"easywrite-service/db"
	"easywrite-service/model"
	"easywrite-service/service"
	"github.com/gin-gonic/gin"
)

// SetUserInfoHandler 更新用户信息
// @Summary 更新用户信息
// @Description 更新用户的个人信息
// @Tags Account
// @Accept json
// @Produce json
// @Param cookie header string true "Cookie"
// @Param info body model.UserInfo true "用户信息"
// @Success 200 {object} model.JsonResponse[any] "成功"
// @Failure 400 {object} model.JsonResponse[any] "请求参数错误"
// @Failure 500 {object} model.JsonResponse[any] "服务器内部错误"
// @Router /account/info [put]
func SetUserInfoHandler(c *gin.Context) {
	var p model.UserInfo
	err := c.Bind(&p)
	if err != nil {
		service.HttpParameterError(c)
		return
	}

	logged, username := IsLoggedWithResponse(c)
	if !logged {
		return
	}

	err = db.Mysql.Where(map[string]any{"username": username}).Updates(&model.User{
		UserInfoReadOnly: model.UserInfoReadOnly{UserInfo: p},
	}).Error
	if err != nil {
		service.HttpServerInternalError(c)
		return
	}
	c.JSON(200, model.JsonResponse[any]{
		Code: 200,
		Msg:  "ok",
		Data: nil,
	})
}

// GetUserInfoHandler 获取用户信息
// @Summary 获取用户信息
// @Description 获取用户的个人信息
// @Tags Account
// @Accept json
// @Produce json
// @Param cookie header string true "Cookie"
// @Success 200 {object} model.JsonResponse[model.UserInfoReadOnly] "成功"
// @Failure 500 {object} model.JsonResponse[any] "服务器内部错误"
// @Router /account/info [get]
func GetUserInfoHandler(c *gin.Context) {
	logged, username := IsLoggedWithResponse(c)
	if !logged {
		return
	}

	var info model.UserInfoReadOnly
	err := db.Mysql.Model(&model.User{}).Where(map[string]any{"username": username}).Take(&info).Error
	if err != nil {
		service.HttpServerInternalError(c)
		return
	}
	c.JSON(200, model.JsonResponse[model.UserInfoReadOnly]{
		Code: 200,
		Msg:  "ok",
		Data: info,
	})
}
