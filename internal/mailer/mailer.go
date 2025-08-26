package mailer

import (
	"bytes"
	"html/template"
	"log"

	"github.com/ipincamp/edsa/internal/config"
	"gopkg.in/gomail.v2"
)

type Mailer struct {
	dialer *gomail.Dialer
	sender string
}

func New(env *config.Env) Mailer {
	dialer := gomail.NewDialer(env.MailHost, env.MailPort, env.MailUsername, env.MailPassword)
	return Mailer{
		dialer: dialer,
		sender: env.MailFrom,
	}
}

func (m Mailer) Send(recipient, templateFile string, data interface{}) error {
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return err
	}

	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}

	msg := gomail.NewMessage()
	msg.SetHeader("To", recipient)
	msg.SetHeader("From", m.sender)
	msg.SetHeader("Subject", subject.String())
	msg.SetBody("text/plain", plainBody.String())
	msg.AddAlternative("text/html", htmlBody.String())

	err = m.dialer.DialAndSend(msg)
	if err != nil {
		log.Printf("Failed to send email to %s: %v", recipient, err)
		return err
	}

	log.Printf("Sent email to %s", recipient)
	return nil
}
