package route

import (
	"easywrite-service/common"
	_ "easywrite-service/docs"
	"easywrite-service/service/account"
	"easywrite-service/service/appversion"
	"easywrite-service/service/bill"
	"easywrite-service/service/feedback"
	"easywrite-service/service/openai"
	"easywrite-service/service/proxy"
	"easywrite-service/service/tencent"
	"easywrite-service/service/textin"
	"easywrite-service/service/welcome"
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Run(host string, port int) {
	r := gin.Default()
	gin.SetMode(gin.DebugMode)
	r.Use(func(c *gin.Context) {
		/*	log.Println(c.Request.Header)
			c.Next()
			log.Println(c.Writer.Header())*/
	})
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Static(common.UploadDir, common.UploadSavePath)
	accountRouter := r.Group("/account")
	accountRouter.POST("/login", account.LoginHandler)
	accountRouter.POST("/login_with_code", account.LoginWithCodeHandler)
	accountRouter.POST("/register", account.RegisterHandler)
	accountRouter.POST("/verify", account.VerificationCodeHandler)
	accountRouter.PUT("/password/reset", account.ResetPasswordHandler)
	accountRouter.PUT("/password/change", account.ChangePasswordHandler)
	accountRouter.PUT("/info", account.SetUserInfoHandler)
	accountRouter.GET("/info", account.GetUserInfoHandler)
	accountRouter.DELETE("/logout", account.LogoutHandler)

	r.GET("/welcome", welcome.WelcomeHandler)

	billRouter := r.Group("/bill")
	billRouter.GET("", bill.GetBillHandler)
	billRouter.POST("", bill.AddBillHandler)
	billRouter.PUT("", bill.UpdateBillHandler)
	billRouter.DELETE("", bill.DeleteBillHandler)
	billRouter.POST("/list", bill.AddBillListHandler)

	r.POST("/feedback", feedback.PostFeedbackHandler)

	r.GET("/app_version", appversion.GetAppVersion)

	r.Any("/proxy", proxy.HttpProxyHandler)

	openAiV1 := r.Group("/openai/v1")
	openAiV1.GET("/models", openai.ListModels)
	openAiV1.GET("/prices", openai.GetPriceHandler)
	openAiV1.POST("/chat/completions", openai.ChatCompletions)
	openAiV1.POST("/audio/transcriptions", openai.AudioTranscriptions)
	openAiV1.POST("/audio/translations", openai.AudioTranslations)
	openAiV1.POST("/images/generations", openai.ImagesGenerations)
	openAiV1.POST("/images/edits", openai.ImagesEdits)
	openAiV1.POST("/images/variations", openai.ImagesVariations)

	openAiTools := r.Group(openai.OpenaiToolsServiceKeyPerfix + "/openai/tools/")
	openAiTools.POST("/audio/duration", openai.GetAudioDuration)
	openAiTools.POST("/image/size", openai.GetImageSize)

	textinApi := r.Group("/textin")
	textinApi.POST("/ai/service/v1/dewarp", textin.Dewarp)
	textinApi.POST("/ai/service/v1/crop_enhance_image", textin.CropEnhanceImage)
	textinApi.POST("/robot/v1.0/api/bills_crop", textin.BillsCrop)

	tencentApi := r.Group("/tencent")
	tencentApi.POST("/sentence_recognition", tencent.SentenceRecognition)

	err := r.Run(fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		panic(err)
	}
}
