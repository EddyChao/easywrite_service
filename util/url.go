package util

import (
	"github.com/gorilla/schema"
	"net/url"
)

var encoder = schema.NewEncoder()

func MustEncodeValue(arg any) string {
	var value = url.Values{}
	err := encoder.Encode(arg, value)
	if err != nil {
		panic(err)
	}
	return value.Encode()
}

func EncodeValue(arg any) (string, error) {
	var value url.Values
	err := encoder.Encode(arg, value)
	return value.Encode(), err
}
