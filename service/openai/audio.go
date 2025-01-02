package openai

import (
	"easywrite-service/common"
	"easywrite-service/tools"
	"easywrite-service/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/weili71/go-filex"
	"io"
	"log"
	"net/http"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

var myServiceHttpClient *resty.Client

func isFFmpegAvailable() bool {
	cmd := exec.Command("ffmpeg", "-version")
	err := cmd.Run()
	if err != nil {
		return false
	}
	return true
}

func init() {
	if !isFFmpegAvailable() {
		log.Println("ffmpeg command not available!")
		return
	}
}

func InitMyServiceAddress(port int) {
	host := "localhost"
	if port == 0 {
		log.Panicln("port is 0")
	}
	myServiceAddress := fmt.Sprintf("http://%s:%d", host, port)
	myServiceHttpClient = resty.New().SetBaseURL(myServiceAddress)
}

// AudioTranscriptions
// @Summary Audio transcriptions
// @Description Transcribe audio to text
// @Tags OpenAI
// @Produce json
// @Param cookie header string true "Cookie"
// @Param file formData file true "Audio file"
// @Success 200
// @Router /openai/v1/audio/transcriptions [post]
func AudioTranscriptions(c *gin.Context) {
	audio(c, "audio/transcriptions")
}

// AudioTranslations
// @Summary Audio translations
// @Description Translate audio from one language to another
// @Tags OpenAI
// @Produce json
// @Param cookie header string true "Cookie"
// @Param file formData file true "Audio file"
// @Success 200
// @Router /openai/v1/audio/translations [post]
func AudioTranslations(c *gin.Context) {
	audio(c, "audio/translations")
}

func audio(c *gin.Context, openaiPath string) {

	if !isFFmpegAvailable() {
		log.Println("ffmpeg command not available!")
		c.JSON(http.StatusBadRequest, gin.H{"error": "service not available! Contact the administrator please!"})
		return
	}
	logged, _ := IsLoggedWithOpenaiResponse(c)
	if !logged {
		return
	}

	var duration time.Duration = 0
	var resp *resty.Response = nil

	var wg sync.WaitGroup
	wg.Add(2)

	reader1, reader2 := util.CopyReader(c.Request.Body)

	contentType := c.GetHeader("content-type")
	go func() {
		defer wg.Done()
		var err error
		resp, err = OpenAiHttpClient.R().
			SetHeader("content-type", contentType).
			SetDoNotParseResponse(true).
			SetBody(reader1).
			Post(openaiPath)
		if err != nil {
			fmt.Println(err)
		}
	}()

	go func() {
		defer wg.Done()
		var err error
		uri := OpenaiToolsServiceKeyPerfix + "/openai/tools/audio/duration"
		resp2, err := myServiceHttpClient.R().
			SetHeader("content-type", contentType).
			SetBody(reader2).
			Post(uri)

		if err != nil {
			fmt.Println(err)
		}
		duration, err = time.ParseDuration(resp2.String())
		if err != nil {
			fmt.Println(err)
		}
	}()
	// 等待两个goroutine完成
	wg.Wait()

	c.Writer.WriteHeader(resp.StatusCode())
	util.CopyResponseHeader(c, resp.RawResponse)
	io.Copy(c.Writer, resp.RawResponse.Body)

	if resp.StatusCode() == 200 {
		expenses := PriceWhisper.
			Mul(decimal.NewFromInt(duration.Milliseconds())).
			Div(decimal.NewFromInt32(60 * 1000))

		fmt.Println("audio: 消费：", expenses, duration.Milliseconds())
	}

}

func GetAudioDuration(c *gin.Context) {
	multipartForm, _ := c.MultipartForm()
	fmt.Println(multipartForm.Value["model"])
	files := multipartForm.File["file"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is null"})
		return
	}
	file := files[0]
	relativePath := filepath.Join("/temp", uuid.NewString()+filepath.Ext(file.Filename))

	f := filex.NewFile(filepath.Join(common.UploadSavePath, relativePath))
	dir := f.ParentFile()
	if !dir.IsExist() {
		err := dir.MkdirAll(0666)
		if err != nil {
			log.Panicln(err)
			return
		}
	}
	destPath := f.Pathname
	defer f.Delete()

	if err := c.SaveUploadedFile(file, destPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	duration, err := tools.GetAudioDuration(destPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.String(200, duration.String())
}
