package helper

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"gopkg.in/gomail.v2"
)

func SendMail(to string, subject string, body string) error {
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

	return dialer.DialAndSend(mailer)
}

func VerifyAccountURLMail(errCh chan<- error, to, subject, verifyLink string) {
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

	errCh <- SendMail(to, subject, body.String())
	return
}

func ResetPasswordSuccessMail(errCh chan<- error, to, subject, username string) {
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
	if err = tmpl.Execute(&body, data); err != nil {
		errCh <- err
		return
	}

	errCh <- SendMail(to, subject, body.String())
	return
}

func ForgotPasswordMail(errCh chan<- error, to, subject, resetLink, username string) {
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
		Username  string
		ResetLink string
	}{
		AppLink:   os.Getenv("CLIENT_URL"),
		AppIcon:   "https://i.ibb.co/Kyp2m0t/cover.png",
		Username:  username,
		ResetLink: resetLink,
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		errCh <- err
		return
	}

	errCh <- SendMail(to, subject, body.String())
	return
}

func SellerOrderHasCompleted(errCh chan<- error, to, subject, buyerEmail, orderID, currentBalance string) {
	dir, err := os.Getwd()
	if err != nil {
		errCh <- err
		return
	}

	//TODO: IMPLEMENT EMAIL HTML TEMPLATE FOR SELLER ORDER COMPLETED
	tmpl, err := template.ParseFiles(fmt.Sprintf("%s/emails/sellerOrderCompleted.html", dir))
	if err != nil {
		errCh <- err
		return
	}

	data := &struct {
		AppLink              string
		AppIcon              string
		OrderID              string
		BuyerEmail           string
		SellerCurrentBalance string
	}{
		AppLink:              os.Getenv("CLIENT_URL"),
		AppIcon:              "https://i.ibb.co/Kyp2m0t/cover.png",
		OrderID:              orderID,
		BuyerEmail:           buyerEmail,
		SellerCurrentBalance: currentBalance,
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		errCh <- err
		return
	}

	errCh <- SendMail(to, subject, body.String())
	return
}
