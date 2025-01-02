package account

import (
	"easywrite-service/constant/code_type"
	"easywrite-service/constant/error_code"
	"easywrite-service/db"
	"easywrite-service/model"
	"easywrite-service/service"
	"easywrite-service/util"
	"github.com/gin-gonic/gin"
)

// RegisterHandler
// @Summary Register
// @Description Handles user registration
// @Tags Account
// @Accept json
// @Produce json
// @Param registerRequest body model.RegistrationParameters true "Registration Request"
// @Success 200 {object} model.JsonResponse[any]
// @Failure 400 {object} model.JsonResponse[any]
// @Router /account/register [post]
func RegisterHandler(c *gin.Context) {
	var p model.RegistrationParameters
	err := c.Bind(&p)
	if err != nil {
		service.HttpParameterError(c)
		return
	}

	if !CheckPasswordIsLegal(c, p.Password) {
		return
	}

	if IsAccountExists(c, p.Username) {
		return
	}

	if IsVerificationCodeCorrect(c, p.VerificationCode, code_type.Register, p.Username) {
		salt := util.GetRandomString(6)
		var user = model.User{
			ID:               0,
			UserInfoReadOnly: model.UserInfoReadOnly{Username: p.Username},
			Password:         util.Sha256Sum(util.Sha256Sum(p.Password) + salt),
			Salt:             salt,
		}

		if util.IsEmail(p.Username) {
			user.Email = p.Username
		} else if util.IsPhone(p.Username) {
			user.Phone = p.Username
		} else {
			c.JSON(200, model.JsonResponse[any]{
				Code: error_code.UsernameError,
				Msg:  "手机号或邮箱错误",
				Data: nil,
			})
			return
		}

		db.Mysql.Create(&user)
		service.HttpOK(c)
	}

}
