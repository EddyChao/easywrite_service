package tencent

import (
	"easywrite-service/service"
	"github.com/gin-gonic/gin"
	asr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/asr/v20190614"
	tencent "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"log"
)

type TencentCloudConfig struct {
	SecretId  string `json:"secret_id"`
	SecretKey string `json:"secret_key"`
}

var (
	tencentCloudConfig TencentCloudConfig
	asrClient          *asr.Client
)

func InitTencentCloudConfig(config TencentCloudConfig) {
	tencentCloudConfig = config
	credential := tencent.NewCredential(
		tencentCloudConfig.SecretId,
		tencentCloudConfig.SecretKey,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "asr.tencentcloudapi.com"
	asrClient, _ = asr.NewClient(credential, "", cpf)
}

func SentenceRecognition(c *gin.Context) {
	request := asr.NewSentenceRecognitionRequest()
	err := c.Bind(request)
	if err != nil {
		service.HttpParameterError(c)
		return
	}
	response, err := asrClient.SentenceRecognition(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		log.Printf("An API error has returned: %s", err)
	}
	if err != nil {
		service.HttpServerInternalError(c)
		return
	}
	c.JSON(200, response)
}
