package util

import (
	gomail "gopkg.in/gomail.v1"
)

// SendEmail 发送
func SendEmail(from, code string, to []string, title string, file []string, server string, port int) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", from)
	msg.SetHeader("To", to...)
	msg.SetHeader("Subject", title)

	if file != nil {
		for _, f := range file {
			gf, err := gomail.OpenFile(f)
			if err != nil {
				return err
			}
			msg.Attach(gf)
		}
	}
	mailer := gomail.NewMailer(server, from, code, port)
	if err := mailer.Send(msg); err != nil {
		return err
	}
	return nil
}
