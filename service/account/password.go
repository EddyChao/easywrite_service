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

// ChangePasswordHandler
// @Summary Change Password
// @Description Handles password change for a user
// @Tags Account
// @Accept json
// @Produce json
// @Param changePasswordRequest body model.ChangePasswordParameters true "Change Password Request"
// @Success 200 {object} model.JsonResponse[any]
// @Failure 400 {object} model.JsonResponse[any]
// @Router /account/password/change [put]
func ChangePasswordHandler(c *gin.Context) {
	var p model.ChangePasswordParameters
	err := c.Bind(&p)
	if err != nil {
		service.HttpParameterError(c)
		return
	}

	if !CheckPasswordIsLegal(c, p.NewPassword) {
		return
	}

	if IsPasswordCorrect(c, p.Username, p.OldPassword) {
		UpdatePassword(p.Username, p.NewPassword)
		service.HttpOK(c)
	}
}

// ResetPasswordHandler
// @Summary Reset Password
// @Description Handles password reset for a user
// @Tags Account
// @Accept json
// @Produce json
// @Param resetPasswordRequest body model.ResetPasswordParameters true "Reset Password Request"
// @Success 200 {object} model.JsonResponse[any]
// @Failure 400 {object} model.JsonResponse[any]
// @Router /account/password/reset [put]
func ResetPasswordHandler(c *gin.Context) {
	var p model.ResetPasswordParameters
	if c.Bind(&p) != nil {
		service.HttpParameterError(c)
		return
	}
	if !CheckPasswordIsLegal(c, p.NewPassword) {
		return
	}
	if IsVerificationCodeCorrect(c, p.VerificationCode, code_type.ResetPassword, p.Username) {
		UpdatePassword(p.Username, p.NewPassword)
		service.HttpOK(c)
	}
}

func UpdatePassword(username string, newPassword string) {
	var user = model.User{}
	salt := util.GetRandomString(6)
	user.Password = util.Sha256Sum(util.Sha256Sum(newPassword) + salt)
	user.Salt = salt
	db.Mysql.Model(&user).Select("password", "salt").Where("username=?", username).Updates(&user)
}

func CheckPasswordIsLegal(c *gin.Context, password string) bool {
	if !util.IsPasswordLegal(password) {
		c.JSON(200, model.JsonResponse[any]{
			Code: error_code.IllegalPassword,
			Msg:  "密码至少包含 数字和英文，长度8-20",
			Data: nil,
		})
		return false
	}
	return true
}
