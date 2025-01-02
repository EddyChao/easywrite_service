package tools

import (
	"gopkg.in/gomail.v2"
)

type EmailConfig struct {
	SMTPUsername string `json:"username"`
	SMTPPassword string `json:"password"`
	SMTPHost     string `json:"host"`
	SMTPPort     int    `json:"smtp_port"`
	FromAddress  string `json:"from_address"`
}

var (
	emailConfig EmailConfig
)

func InitEmail(config EmailConfig) {
	emailConfig = config
}

func SendEmail(m *gomail.Message) error {
	m.SetHeader("From", emailConfig.FromAddress)

	d := gomail.NewDialer(emailConfig.SMTPHost, emailConfig.SMTPPort, emailConfig.SMTPUsername, emailConfig.SMTPPassword)
	return d.DialAndSend(m)
}

func SendCodeToEmail(to string, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("To", to)
	m.SetHeader("Subject", "验证码")
	m.SetBody("text/plain", body)
	return SendEmail(m)
}
