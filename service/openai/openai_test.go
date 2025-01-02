package openai

import (
	"fmt"
	"github.com/shopspring/decimal"
	"testing"
	"time"
)

func TestExpenses(t *testing.T) {
	d, _ := time.ParseDuration("60m")
	expenses := PriceWhisper.Mul(decimal.NewFromInt(d.Milliseconds())).Div(decimal.NewFromInt32(60 * 1000))
	fmt.Println(expenses)
	fmt.Println(GetModelPrice("gpt-4-0314"))
	fmt.Println(GetModelPrice("gpt-4-0314"))
}

func TestKey(t *testing.T) {
	fmt.Println(GenerateKey(48))
}
