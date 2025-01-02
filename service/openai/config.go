package openai

import (
	"github.com/go-resty/resty/v2"
	"github.com/sashabaranov/go-openai"
	"log"
)

type OpenAiConfig struct {
	Key   string `json:"key"`
	Proxy string `json:"proxy"`
}

var (
	openAiConfig     OpenAiConfig
	OpenAiHttpClient *resty.Client
	OpenAiSDKClient  *openai.Client
)

const (
	openAiBaseUrl = "https://api.openai.com/v1/"
)

func InitOpenAi(config OpenAiConfig) {
	openAiConfig = config

	OpenAiHttpClient = resty.New()
	OpenAiHttpClient.SetBaseURL(openAiBaseUrl)
	OpenAiHttpClient.OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
		// 添加 token 到请求头
		req.SetHeader("Authorization", "Bearer "+openAiConfig.Key)
		return nil
	})
	if openAiConfig.Proxy != "" {
		OpenAiHttpClient.SetProxy(openAiConfig.Proxy)
	}
	OpenAiHttpClient.OnError(func(request *resty.Request, err error) {
		log.Println(err)
	})
	//OpenAiHttpClient.Debug = true
	config1 := openai.DefaultConfig(config.Key)
	config1.HTTPClient = OpenAiHttpClient.GetClient()
	config1.BaseURL = openAiBaseUrl
	OpenAiSDKClient = openai.NewClientWithConfig(config1)
}
