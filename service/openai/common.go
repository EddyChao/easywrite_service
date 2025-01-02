package openai

import (
	"easywrite-service/model"
	"easywrite-service/service/account"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

var OpenaiToolsServiceKeyPerfix = "/" + uuid.NewString()

func IsLoggedWithOpenaiResponse(c *gin.Context) (bool, string) {
	logged, username, err := account.IsLogged(c)
	fmt.Println("username", username)
	if !logged || err != nil {
		c.JSON(http.StatusUnauthorized, model.OpenaiErrorResponse{Error: model.OpenaiErrorDetails{
			Message: "请登录",
			Type:    "invalid_request_error",
			Param:   nil,
			Code:    "invalid_cookie",
		}})
	}
	return logged, username
}
