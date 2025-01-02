package openai

import (
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"net/http"
)

var (
	// 单位：1k token
	ExchangeRate  = NewDecimalFromString("7")
	Magnification = NewDecimalFromString("1.5")

	// 单位：张
	PriceImage256  = NewDecimalFromString("0.016").Mul(ExchangeRate).Mul(Magnification)
	PriceImage512  = NewDecimalFromString("0.018").Mul(ExchangeRate).Mul(Magnification)
	PriceImage1024 = NewDecimalFromString("0.02").Mul(ExchangeRate).Mul(Magnification)
	// 单位：分钟
	PriceWhisper = NewDecimalFromString("0.006").Mul(ExchangeRate).Mul(Magnification)
)

type ChatModelPrice struct {
	InputPrice  decimal.Decimal `json:"input_price"`
	OutputPrice decimal.Decimal `json:"output_price"`
}

var (
	PriceGpt4        = ChatModelPrice{NewDecimalFromString("0.03"), NewDecimalFromString("0.06")}
	PriceGpt432k     = ChatModelPrice{NewDecimalFromString("0.06"), NewDecimalFromString("0.12")}
	PriceGpt3Dot5    = ChatModelPrice{NewDecimalFromString("0.0015"), NewDecimalFromString("0.002")}
	PriceGpt3Dot516k = ChatModelPrice{NewDecimalFromString("0.003"), NewDecimalFromString("0.004")}
)

var modelPrices = map[string]ChatModelPrice{
	"gpt-4":                  PriceGpt4,
	"gpt-4-0314":             PriceGpt4,
	"gpt-4-0613":             PriceGpt4,
	"gpt-4-32k":              PriceGpt432k,
	"gpt-4-32k-0314":         PriceGpt432k,
	"gpt-4-32k-0613":         PriceGpt432k,
	"gpt-3.5-turbo":          PriceGpt3Dot5,
	"gpt-3.5-turbo-0301":     PriceGpt3Dot5,
	"gpt-3.5-turbo-0613":     PriceGpt3Dot5,
	"gpt-3.5-turbo-16k":      PriceGpt3Dot516k,
	"gpt-3.5-turbo-16k-0613": PriceGpt3Dot516k,
}

var imagePrices = map[string]decimal.Decimal{
	"1024x1024": PriceImage1024,
	"512x512":   PriceImage512,
	"256x256":   PriceImage256,
}

func GetImagePrices(size string) decimal.Decimal {
	return imagePrices[size]
}

func GetModelPrice(model string) ChatModelPrice {
	price, ok := modelPrices[model]
	if !ok {
		return ChatModelPrice{decimal.Zero, decimal.Zero}
	}

	price.InputPrice = price.InputPrice.Mul(ExchangeRate).Mul(Magnification)
	price.OutputPrice = price.OutputPrice.Mul(ExchangeRate).Mul(Magnification)
	return price
}

// GetPriceHandler
// @Summary Get prices
// @Description Get the pricing details for OpenAI services
// @Tags OpenAI
// @Produce json
// @Success 200 {object} Price
// @Router /openai/v1/prices [get]
func GetPriceHandler(c *gin.Context) {
	priceTable := getPriceTable()
	c.JSON(http.StatusOK, priceTable)
}

type Price struct {
	Unit   string `json:"unit"`
	Values any    `json:"values"`
}

type PriceTable struct {
	Chat    Price `json:"chat"`
	Image   Price `json:"image"`
	Whisper Price `json:"whisper"`
}

func getPriceTable() PriceTable {
	chatList := calculatePrices(modelPrices, ExchangeRate, Magnification)
	imageMap := getImageMap()
	whisperPrice := getWhisperPrice()

	priceTable := PriceTable{
		Chat: Price{
			Unit:   "￥/1k token",
			Values: chatList,
		},
		Image: Price{
			Unit:   "￥/1 picture",
			Values: imageMap,
		},
		Whisper: Price{
			Unit:   "￥/1 minute",
			Values: whisperPrice,
		},
	}

	return priceTable
}

func getImageMap() map[string]decimal.Decimal {
	imageMap := map[string]decimal.Decimal{
		"image256":  PriceImage256,
		"image512":  PriceImage512,
		"image1024": PriceImage1024,
	}

	return imageMap
}

func getWhisperPrice() decimal.Decimal {
	return PriceWhisper
}

func calculatePrices(modelPrices map[string]ChatModelPrice, exchangeRate, magnification decimal.Decimal) map[string]ChatModelPrice {
	newModelPrices := make(map[string]ChatModelPrice)
	for key, price := range modelPrices {
		price.InputPrice = price.InputPrice.Mul(exchangeRate).Mul(magnification)
		price.OutputPrice = price.OutputPrice.Mul(exchangeRate).Mul(magnification)
		newModelPrices[key] = price
	}
	return newModelPrices
}

func NewDecimalFromString(str string) decimal.Decimal {
	num, err := decimal.NewFromString(str)
	if err != nil {
		panic(err)
	}
	return num
}
