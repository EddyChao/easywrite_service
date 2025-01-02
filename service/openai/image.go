package openai

import (
	"easywrite-service/util"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/shopspring/decimal"
	"io"
	"sync"
)

// ImagesGenerations
// @Summary Audio translations
// @Description Translate audio from one language to another
// @Tags OpenAI
// @Produce json
// @Param cookie header string true "Cookie"
// @Param file formData file true "Audio file"
// @Success 200
// @Router /openai/v1/audio/translations [post]
func ImagesGenerations(c *gin.Context) {
	image(c, "images/generations")
}

// ImagesEdits
// @Summary Images edits
// @Description Edit images based on a given prompt
// @Tags OpenAI
// @Accept json
// @Produce json
// @Param cookie header string true "Cookie"
// @Success 200
// @Router /openai/v1/images/edits [post]
func ImagesEdits(c *gin.Context) {
	image(c, "images/edits")
}

// ImagesVariations
// @Summary Images variations
// @Description Generate variations of images based on a given prompt
// @Tags OpenAI
// @Accept json
// @Produce json
// @Param cookie header string true "Cookie"
// @Success 200
// @Router /openai/v1/images/variations [post]
func ImagesVariations(c *gin.Context) {
	image(c, "images/variations")
}

func image(c *gin.Context, openaiPath string) {

	logged, _ := IsLoggedWithOpenaiResponse(c)
	if !logged {
		return
	}

	var imageRequestBody *ImageRequestBody = nil
	var resp *resty.Response = nil

	var wg sync.WaitGroup
	wg.Add(2)

	reader1, reader2 := util.CopyReader(c.Request.Body)

	contentType := c.GetHeader("content-type")
	go func() {
		defer wg.Done()
		resp, _ = OpenAiHttpClient.R().
			SetHeader("content-type", contentType).
			SetDoNotParseResponse(true).
			SetBody(reader1).
			Post(openaiPath)
	}()

	go func() {
		defer wg.Done()
		resp2, _ := myServiceHttpClient.R().
			SetHeader("content-type", contentType).
			SetBody(reader2).
			Post(OpenaiToolsServiceKeyPerfix + "/openai/tools/image/size")

		err := json.Unmarshal(resp2.Body(), &imageRequestBody)
		if err != nil {
			return
		}
	}()
	// 等待两个goroutine完成
	wg.Wait()

	c.Writer.WriteHeader(resp.StatusCode())
	util.CopyResponseHeader(c, resp.RawResponse)
	io.Copy(c.Writer, resp.RawResponse.Body)

	if resp.StatusCode() == 200 {
		expenses := GetImagePrices(imageRequestBody.Size).
			Mul(decimal.NewFromInt(int64(imageRequestBody.N)))
		fmt.Println("image 消费：", expenses)
	}

}

type ImageRequestBody struct {
	Size string `json:"size" form:"size"`
	N    int    `json:"n" form:"n"`
}

func GetImageSize(c *gin.Context) {

	var imageRequestBody ImageRequestBody
	err := c.Bind(&imageRequestBody)
	if err != nil {
		return
	}

	fmt.Println(imageRequestBody)
	c.JSON(200, imageRequestBody)

}
