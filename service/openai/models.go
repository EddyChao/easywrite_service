package openai

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// ListModels
// @Summary List models
// @Description Get a list of available models
// @Tags OpenAI
// @Produce json
// @Success 200 {array} Model
// @Router /openai/v1/models [get]
func ListModels(c *gin.Context) {
	resp, err := OpenAiSDKClient.ListModels(context.Background())
	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, "Error")
		return
	}
	c.JSON(200, resp)
}
