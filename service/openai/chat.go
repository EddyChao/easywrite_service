package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"github.com/shopspring/decimal"
	"log"
	"net/http"
)

// ChatCompletions
// @Summary Chat completions
// @Description Generate completions for a given prompt
// @Tags OpenAI
// @Accept json
// @Produce json
// @Param cookie header string true "Cookie"
// @Param input body ChatCompletionRequest true "Chat completion request object"
// @Success 200 {object} ChatCompletionResponse
// @Router /openai/v1/chat/completions [post]
func ChatCompletions(c *gin.Context) {

	logged, _ := IsLoggedWithOpenaiResponse(c)
	if !logged {
		return
	}

	fmt.Println(c.Request.Header)
	var req openai.ChatCompletionRequest
	err := c.BindJSON(&req)
	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, "Error")
		return
	}

	a := NumTokensFromMessages(req.Messages, req.Model)
	b := NumTokensFromFunctions(req.Functions, req.Model)
	promptToken := a + b

	fmt.Println("promptToken:", promptToken, "a: ", a, "b", b)

	if !req.Stream {

		resp, err := OpenAiSDKClient.CreateChatCompletion(context.Background(), req)
		completionToken := NumTokensFromChatCompletion(resp, req.Model)

		if err == nil {
			price := GetModelPrice(resp.Model)
			expensesInput := price.InputPrice.Mul(decimal.NewFromInt(int64(resp.Usage.PromptTokens))).Div(decimal.NewFromInt32(1000))
			expensesOutput := price.OutputPrice.Mul(decimal.NewFromInt(int64(resp.Usage.CompletionTokens))).Div(decimal.NewFromInt32(1000))
			expenses := expensesInput.Add(expensesOutput)
			fmt.Printf("chat 消费：input: %v, ouput: %v, total: %v\n", expensesInput, expensesOutput, expenses)
		}

		fmt.Println("completionToken:", completionToken)
		usage := resp.Usage
		fmt.Printf("%#v", usage)
		c.JSON(200, resp)

	} else {

		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")

		stream, err := OpenAiSDKClient.CreateChatCompletionStream(context.Background(), req)
		if err != nil {
			log.Println(err)
			c.String(http.StatusInternalServerError, "Error")
			return
		}
		defer stream.Close()

		go func() {
			select {
			case <-c.Request.Context().Done():
				stream.Close()
			}
		}()

		text := ""
		var firstResponse *openai.ChatCompletionStreamResponse = nil
		var functionsCall *openai.FunctionCall = nil
		fmt.Printf("Stream response: ")

		for {
			response, err := stream.Recv()

			if firstResponse == nil {
				firstResponse = &response
			}

			if len(response.Choices) > 0 {
				delta := response.Choices[0].Delta
				text += delta.Content
				if delta.FunctionCall != nil {
					if functionsCall == nil {
						functionsCall = delta.FunctionCall
					} else {
						functionsCall.Arguments += delta.FunctionCall.Arguments
					}
				}
			}

			if err != nil {
				c.SSEvent("", " [DONE]")
				model := firstResponse.Model
				completionToken := NumTokensFromText(text, model) + NumTokensFromFunctionsCall(functionsCall, model)

				price := GetModelPrice(model)
				expensesInput := price.InputPrice.Mul(decimal.NewFromInt(int64(promptToken))).Div(decimal.NewFromInt32(1000))
				expensesOutput := price.OutputPrice.Mul(decimal.NewFromInt(int64(completionToken))).Div(decimal.NewFromInt32(1000))
				expenses := expensesInput.Add(expensesOutput)
				fmt.Printf("chat 消费：input: %v, ouput: %v, total: %v\n", expensesInput, expensesOutput, expenses)

				fmt.Println("completionToken:", completionToken)
				return
			}

			data, _ := json.Marshal(response)
			c.SSEvent("", " "+string(data))
		}
	}
}
