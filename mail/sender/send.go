package sender

import (
	"fmt"
	conf "github.com/glossd/pokergloss/goconf"
	"net/smtp"
)

func SendEmail(toEmail, subject, body string) error {
	smtpAuth := smtp.PlainAuth("", conf.Props.Mail.Username, conf.Props.Mail.Password, conf.Props.Mail.Host)
	address := fmt.Sprintf("%s:%d", conf.Props.Mail.Host, conf.Props.Mail.Port)
	msg := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n%s", conf.Props.From, toEmail, subject, body)
	return smtp.SendMail(address, smtpAuth, conf.Props.From, []string{toEmail}, []byte(msg))
}
