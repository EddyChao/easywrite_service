package account

import (
	"easywrite-service/constant/error_code"
	"easywrite-service/constant/redis_prefix"
	"easywrite-service/constant/verification_code"
	"easywrite-service/db"
	"easywrite-service/model"
	"easywrite-service/service"
	"easywrite-service/tools"
	"easywrite-service/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

var appName string

func InitAppName(name string) {
	appName = name
}

// VerificationCodeHandler
// @Summary Verification Code
// @Description Handles sending verification code to the user
// @Tags Account
// @Accept json
// @Produce json
// @Param verificationCodeRequest body model.VerificationParameters true "Verification Code Request"
// @Success 200 {object} model.JsonResponse[any]
// @Failure 400 {object} model.JsonResponse[string]
// @Router /account/verify [post]
func VerificationCodeHandler(c *gin.Context) {
	var p model.VerificationParameters
	err := c.Bind(&p)
	if err != nil {
		service.HttpParameterError(c)
		return
	}

	//判断是否超过每日最大发送量
	codeCountKey := util.GetKey(redis_prefix.GetCodeCount, p.PhoneOrEmail)
	countIsExist := db.RedisIsExist(codeCountKey)
	fmt.Println(codeCountKey)
	var count int
	if countIsExist {
		count, _ = db.Redis.Get(db.Context, codeCountKey).Int()
		if count >= verification_code.SingleDayMaximum {
			ttl := db.Redis.TTL(db.Context, codeCountKey).Val()
			c.JSON(200, model.JsonResponse[string]{
				Code: error_code.CodeExceededDailyMax,
				Msg:  fmt.Sprintf("请 %s 后再试", ttl.String()),
				Data: ttl.String(),
			})
			return
		}
	}

	//判断是否超过获取间隔
	codeKey := util.GetKey(redis_prefix.Code, p.CodeType, p.PhoneOrEmail)
	fmt.Println(codeKey)

	if ttl := db.Redis.TTL(db.Context, codeKey).Val(); ttl > 0 {
		if elapsedTime := verification_code.CodeTTL - ttl; elapsedTime < time.Minute*1 {
			countdown := verification_code.Interval - elapsedTime
			c.JSON(200, model.JsonResponse[string]{
				Code: error_code.RequestTooFrequent,
				Msg:  fmt.Sprintf("请在 %s 后获取", countdown.String()),
				Data: countdown.String(),
			})
			return
		}
	}

	//发送
	code := util.GetRandomCode(6)

	if util.IsEmail(p.PhoneOrEmail) {
		err := tools.SendCodeToEmail(p.PhoneOrEmail, fmt.Sprintf("【%s】 您的验证码是：%s。此验证码10分钟后失效，请勿泄露。", appName, code))
		if err != nil {
			c.JSON(200, model.JsonResponse[any]{
				Code: error_code.SendCodeFailed,
				Msg:  "发送验证码失败",
				Data: nil,
			})
			log.Println(err)
			return
		}
	} else if util.IsPhone(p.PhoneOrEmail) {
		sendToPhone(p.PhoneOrEmail)
	} else {
		c.JSON(200, model.JsonResponse[any]{
			Code: error_code.UsernameError,
			Msg:  "手机号或邮箱错误",
			Data: nil,
		})
		return
	}

	//保存验证码
	result, err1 := db.Redis.Set(db.Context, codeKey, code, verification_code.ValidPeriod).Result()

	log.Println(result, err1)

	//更新计数
	if countIsExist {
		db.Redis.Set(db.Context, codeCountKey, count+1, redis.KeepTTL)
	} else {
		db.Redis.Set(db.Context, codeCountKey, 1, time.Hour*24)
	}

	service.HttpOK(c)
}

func sendToPhone(phone string) {

}
