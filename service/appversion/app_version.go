package appversion

import (
	"easywrite-service/db"
	"easywrite-service/model"
	"easywrite-service/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetAppVersion 获取应用版本信息
// @Summary 获取应用版本信息
// @Description 根据应用ID获取最新的应用版本信息
// @Tags AppVersion
// @Accept json
// @Produce json
// @Param appid query string true "应用ID"
// @Success 200 {object} model.JsonResponse[model.AppVersion] "成功"
// @Failure 500 {object} model.JsonResponse[any] "服务器内部错误"
// @Router /app_version [get]
func GetAppVersion(c *gin.Context) {
	// 获取参数
	appid := c.Query("appid")
	// 从数据库中查询最新版本信息
	latestVersion := new(model.AppVersion)
	err := db.Mysql.Where("appid = ?", appid).Order("version_code desc").First(latestVersion).Error
	if err == gorm.ErrRecordNotFound {
		latestVersion = nil
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		service.HttpServerInternalError(c)
		return
	}
	// 返回结果
	c.JSON(200, model.JsonResponse[model.AppVersion]{
		Code: 200,
		Msg:  "ok",
		Data: *latestVersion,
	})
}
