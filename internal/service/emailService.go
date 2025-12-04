package service

import (
	"crypto/tls"
	"fmt"
	"gopkg.in/gomail.v2"
	"log"
	"vibe-music-server/internal/config"
	"vibe-music-server/internal/pkg/util"
)

type EmailService struct {
	from   string // 显示发件人
	dialer *gomail.Dialer
}

func NewEmailService() *EmailService {
	mailConf := config.Get().Mail
	d := gomail.NewDialer(mailConf.Host, mailConf.Port, mailConf.User, mailConf.Password)
	if mailConf.SSL {
		d.SSL = true
		d.TLSConfig = &tls.Config{ServerName: mailConf.Host}
	} else {
		d.SSL = false
	}
	return &EmailService{
		from:   fmt.Sprintf("%s <%s>", mailConf.SenderName, mailConf.User),
		dialer: d,
	}
}

func (e EmailService) SendEmail(to string, subject string, content string) bool {
	m := gomail.NewMessage()
	m.SetHeader("From", e.from) // 任意发件人
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", content)

	if err := e.dialer.DialAndSend(m); err != nil {
		log.Printf("EmailService.SendEmail err: %v\n", err)
		return false
	}
	return true
}

func (e EmailService) SendVerificationCodeEmail(email string) string {
	code := util.GenRandomDigitalCode(6)
	subject := config.Get().App.Name + " - Verification Code"
	content := fmt.Sprintf("<h1>Your verification code is: %s</h1><p>Please use this code to complete your action. The code is valid for 5 minutes.</p>", code)
	if e.SendEmail(email, subject, content) {
		return code
	} else {
		return ""
	}
}
