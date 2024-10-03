package util

import (
	"bytes"
	"fmt"
	"html/template"
	"os"

	"gopkg.in/gomail.v2"
)

func SendMail(errCh chan<- error, to string, subject string, body string) {
	defer close(errCh)
	var (
		CONFIG_SMTP_HOST     = "smtp.gmail.com"
		CONFIG_SMTP_PORT     = 587
		CONFIG_SENDER_NAME   = os.Getenv("EMAIL_SENDER_NAME")
		CONFIG_AUTH_EMAIL    = os.Getenv("EMAIL_SENDER")
		CONFIG_AUTH_PASSWORD = os.Getenv("PASSWORD_SENDER")
	)

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", CONFIG_AUTH_EMAIL)
	mailer.SetHeader("To", to)
	mailer.SetAddressHeader("Cc", CONFIG_AUTH_EMAIL, CONFIG_SENDER_NAME)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", body)

	dialer := gomail.NewDialer(
		CONFIG_SMTP_HOST,
		CONFIG_SMTP_PORT,
		CONFIG_AUTH_EMAIL,
		CONFIG_AUTH_PASSWORD,
	)

	errCh <- dialer.DialAndSend(mailer)
}

func VerifyAccountURLMail(errCh chan<- error, to, subject, verifyLink string) {
	defer close(errCh)
	dir, err := os.Getwd()
	if err != nil {
		errCh <- err
		return
	}

	tmpl, err := template.ParseFiles(fmt.Sprintf("%s/emails/verifyEmail.html", dir))
	if err != nil {
		errCh <- err
		return
	}

	data := &struct {
		AppLink    string
		AppIcon    string
		VerifyLink string
	}{
		AppLink:    os.Getenv("CLIENT_URL"),
		AppIcon:    "https://i.ibb.co/Kyp2m0t/cover.png",
		VerifyLink: verifyLink,
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		errCh <- err
		return
	}

	go SendMail(errCh, to, subject, body.String())

}

func ResetPasswordSuccessMail(errCh chan<- error, to, subject, username string) {
	defer close(errCh)
	dir, err := os.Getwd()
	if err != nil {
		errCh <- err
		return
	}

	tmpl, err := template.ParseFiles(fmt.Sprintf("%s/emails/resetPasswordSuccess.html", dir))
	if err != nil {
		errCh <- err
		return
	}

	data := &struct {
		AppLink  string
		AppIcon  string
		Username string
	}{
		AppLink:  os.Getenv("CLIENT_URL"),
		AppIcon:  "https://i.ibb.co/Kyp2m0t/cover.png",
		Username: username,
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		errCh <- err
		return
	}

	go SendMail(errCh, to, subject, body.String())
}

func ResetPasswordMail(errCh chan<- error, to, subject, resetLink string) {
	defer close(errCh)
	dir, err := os.Getwd()
	if err != nil {
		errCh <- err
		return
	}

	tmpl, err := template.ParseFiles(fmt.Sprintf("%s/emails/resetPassword.html", dir))
	if err != nil {
		errCh <- err
		return
	}

	data := &struct {
		AppLink   string
		AppIcon   string
		ResetLink string
	}{
		AppLink:   os.Getenv("CLIENT_URL"),
		AppIcon:   "https://i.ibb.co/Kyp2m0t/cover.png",
		ResetLink: resetLink,
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		errCh <- err
		return
	}

	go SendMail(errCh, to, subject, body.String())
}
