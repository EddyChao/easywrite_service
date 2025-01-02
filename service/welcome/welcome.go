package welcome

import (
	"easywrite-service/model"
	"easywrite-service/service/account"
	"github.com/gin-gonic/gin"
)

func WelcomeHandler(c *gin.Context) {
	logged, _ := account.IsLoggedWithResponse(c)
	if !logged {
		return
	}
	c.JSON(200, model.JsonResponse[any]{
		Code: 200,
		Msg:  "hello,you are logged in!",
		Data: nil,
	})
}
