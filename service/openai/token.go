package openai

import (
	"fmt"
	"log"
	"strings"

	"github.com/pkoukk/tiktoken-go"
	"github.com/sashabaranov/go-openai"
)

func NumTokensFromMessages(messages []openai.ChatCompletionMessage, model string) (numTokens int) {
	tkm, err := tiktoken.EncodingForModel(model)
	if err != nil {
		err = fmt.Errorf("encoding for model: %v", err)
		log.Println(err)
		return
	}

	tokensPerMessage, tokensPerName := getTokensPerMessageAndTokensPerName(model)

	var nameMap = make(map[string]string)

	for _, message := range messages {
		numTokens += tokensPerMessage
		if message.Content != "" {
			numTokens += len(tkm.Encode(message.Content, nil, nil))
		}
		numTokens += len(tkm.Encode(message.Role, nil, nil))
		if message.FunctionCall != nil {
			numTokens += 1
			numTokens += len(tkm.Encode(message.FunctionCall.Name, nil, nil))
			numTokens += len(tkm.Encode(message.FunctionCall.Arguments, nil, nil))
		}
		if message.Name != "" {
			numTokens += len(tkm.Encode(message.Name, nil, nil))
			if _, ok := nameMap[message.Name]; !ok {
				numTokens += tokensPerName
				nameMap[message.Name] = ""
			}
		}
	}
	numTokens += 3 // every reply is primed with <|start|>assistant<|message|>
	return numTokens
}

func NumTokensFromChatCompletion(response openai.ChatCompletionResponse, model string) (numTokens int) {
	tkm, err := tiktoken.EncodingForModel(model)
	if err != nil {
		err = fmt.Errorf("encoding for model: %v", err)
		log.Println(err)
		return
	}

	if len(response.Choices) <= 0 {
		return 0
	}
	message := response.Choices[0].Message
	numTokens += len(tkm.Encode(message.Content, nil, nil))
	if message.FunctionCall != nil {
		numTokens += 4
		numTokens += len(tkm.Encode(message.FunctionCall.Name, nil, nil))
		numTokens += len(tkm.Encode(message.FunctionCall.Arguments, nil, nil))
	}
	if message.Name != "" {
		numTokens += len(tkm.Encode(message.Name, nil, nil))
	}
	return numTokens
}

func NumTokensFromFunctionsCall(response *openai.FunctionCall, model string) (numTokens int) {
	if response == nil {
		return 0
	}
	tkm, err := tiktoken.EncodingForModel(model)
	if err != nil {
		err = fmt.Errorf("encoding for model: %v", err)
		log.Println(err)
		return
	}

	numTokens += 4
	numTokens += len(tkm.Encode(response.Name, nil, nil))
	numTokens += len(tkm.Encode(response.Arguments, nil, nil))

	return numTokens
}

func getTokensPerMessageAndTokensPerName(model string) (int, int) {
	var tokensPerMessage, tokensPerName int
	switch model {
	case "gpt-3.5-turbo-0613",
		"gpt-3.5-turbo-16k-0613",
		"gpt-4-0314",
		"gpt-4-32k-0314",
		"gpt-4-0613",
		"gpt-4-32k-0613":
		tokensPerMessage = 3
		tokensPerName = 1
	case "gpt-3.5-turbo-0301":
		tokensPerMessage = 4 // every message follows <|start|>{role/name}\n{content}<|end|>\n
		tokensPerName = -1   // if there's a name, the role is omitted
	default:
		if strings.Contains(model, "gpt-3.5-turbo") || strings.Contains(model, "gpt-4") {
			tokensPerMessage = 3
			tokensPerName = 1
		}
	}
	return tokensPerMessage, tokensPerName
}

func NumTokensFromText(message string, model string) int {
	tkm, err := tiktoken.EncodingForModel(model)
	if err != nil {
		err = fmt.Errorf("encoding for model: %v", err)
		log.Println(err)
		return 0
	}

	numTokens := len(tkm.Encode(message, nil, nil))
	return numTokens
}

func NumTokensFromFunctions(functions []openai.FunctionDefinition, model string) int {
	var numTokens int

	tkm, err := tiktoken.EncodingForModel(model)
	if err != nil {
		err = fmt.Errorf("encoding for model: %v", err)
		log.Println(err)
		return 0
	}

	for _, function := range functions {
		functionTokens := len(tkm.Encode(function.Name, nil, nil))
		functionTokens += len(tkm.Encode(function.Description, nil, nil))

		parameters := function.Parameters.(map[string]any)
		if properties, ok := parameters["properties"].(map[string]interface{}); ok {
			for propertiesKey, v := range properties {
				functionTokens += len(tkm.Encode(propertiesKey, nil, nil))
				fields := v.(map[string]interface{})
				for field, fieldValue := range fields {
					switch field {
					case "type":
						functionTokens += 2
						functionTokens += len(tkm.Encode(fieldValue.(string), nil, nil))
					case "description":
						description := fieldValue.(string)
						if description == "" {
							break
						}
						functionTokens += 2
						functionTokens += len(tkm.Encode(description, nil, nil))
					case "enum":
						enums := fieldValue.([]any)
						if len(enums) == 0 {
							break
						}
						functionTokens -= 3
						for _, o := range enums {
							functionTokens += 3
							functionTokens += len(tkm.Encode(o.(string), nil, nil))
						}
					default:
						fmt.Printf("Warning: not supported field %s\n", field)
					}
				}
				functionTokens += 11
			}
		}

		numTokens += functionTokens
	}

	if numTokens > 0 {
		numTokens += 12
	}

	return numTokens
}
