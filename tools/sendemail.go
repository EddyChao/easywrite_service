package tools

import (
	"github.com/wneessen/go-mail"
)

type EmailConfig struct {
	SMTPUsername string `json:"username"`
	SMTPPassword string `json:"password"`
	SMTPHost     string `json:"host"`
	SMTPPort     int    `json:"smtp_port"`
	FromAddress  string `json:"from_address"`
	DisplayName  string `json:"display_name"`
	UseSSL       bool   `json:"use_ssl"`
}

var (
	emailConfig EmailConfig
)

func InitEmail(config EmailConfig) {
	emailConfig = config
}

func SendEmail(m *mail.Msg) error {
	if err := m.FromFormat(emailConfig.DisplayName, emailConfig.FromAddress); err != nil {
		return err
	}

	options := []mail.Option{
		mail.WithPort(emailConfig.SMTPPort), mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(emailConfig.SMTPUsername), mail.WithPassword(emailConfig.SMTPPassword),
	}
	if emailConfig.UseSSL {
		options = append(options, mail.WithSSLPort(true))
	}
	c, err := mail.NewClient(emailConfig.SMTPHost, options...)
	if err != nil {
		return err
	}
	return c.DialAndSend(m)
}

func SendCodeToEmail(to string, body string) error {
	m := mail.NewMsg()
	if err := m.To(to); err != nil {
		return err
	}
	m.Subject("验证码")
	m.SetBodyString(mail.TypeTextPlain, body)
	return SendEmail(m)
}
