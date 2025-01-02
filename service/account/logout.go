package account

import (
	"easywrite-service/common"
	"easywrite-service/model"
	"easywrite-service/service"
	"github.com/gin-gonic/gin"
	"log"
)

// LogoutHandler
// @Summary Logout
// @Description Handles user logout
// @Tags Account
// @Produce json
// @Param cookie header string true "Cookie"
// @Success 200 {object} model.JsonResponse[string]
// @Failure 400 {object} model.JsonResponse[any]
// @Router /account/logout [delete]
func LogoutHandler(c *gin.Context) {
	logged, _ := IsLoggedWithResponse(c)
	if !logged {
		return
	}
	session, err := common.Sessions.Get(c.Request, "session-key")
	session.Options.MaxAge = -1
	if err = session.Save(c.Request, c.Writer); err != nil {
		service.HttpServerInternalError(c)
		log.Println("failed deleting session: ", err)
		return
	}
	c.JSON(200, model.JsonResponse[string]{
		Code: 200,
		Msg:  "ok",
		Data: "退出登录成功",
	})
}
