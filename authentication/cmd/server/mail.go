package main

import (
	"crypto/tls"
	"fmt"
	mail "gopkg.in/mail.v2"
	"telython/authentication/pkg/cfg"
)

var mailDialer *mail.Dialer

func initMailClient() {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		// "true" is only needed when SSL/TLS certificate is not valid on server.
		// In production this should be set to false.
		ServerName: cfg.GetString("smtpHost"),
	}

	//cfg.smtpHost

	mailDialer = mail.NewDialer(cfg.GetString("smtpHost"), cfg.Value.GetInt("smtpPort"), cfg.GetString("smtpSender"), cfg.GetString("smtpPassword"))
	mailDialer.TLSConfig = tlsConfig
}

func TestMail() {
	err := sendEmail("maxim2006722@gmail.com", "Is it in junk?", "<h1>Test html email from telython-auth service.</h1>")
	if err != nil {
		panicIfError(err)
	}
}

func sendEmail(receiver string, subject string, html string) error {
	m := mail.NewMessage()

	m.SetHeader("From", cfg.GetString("smtpSender"))
	m.SetHeader("To", receiver)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", html)

	err := mailDialer.DialAndSend(m)
	if err != nil {
		return err
	}
	return nil
}

func sendMultiEmail(receivers []string, subject string, html string) {
	for i := 0; i < len(receivers); i++ {
		go func() {
			err := sendEmail(receivers[i], subject, html)
			if err != nil {
				fmt.Println(err.Error())
			}
		}()
	}
}
