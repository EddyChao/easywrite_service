package model

type OpenaiErrorResponse struct {
	Error OpenaiErrorDetails `json:"error"`
}

type OpenaiErrorDetails struct {
	Message string `json:"message"`

	Type string `json:"type"`

	Param any `json:"param"`

	Code string `json:"code"`
}
