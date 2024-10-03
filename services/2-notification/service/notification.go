package service

import (
	"fmt"
	"sync"

	"github.com/Akihira77/gojobber/services/2-notification/helper"
)

type NotificationService struct {
	errCh chan error
}

type NotificationServiceImpl interface {
	UserVerifyingEmail(receiverEmail, htmlTemplateName, verifyLink string) error
	UserForgotPassword(receiverEmail, htmlTemplateName, resetLink, username string) error
	UserSucessResetPassword(receiverEmail, htmlTemplateName, username string) error
	SendEmailChatNotification(receiverEmail, senderEmail, message string) error
}

func NewNotificationService() NotificationServiceImpl {
	return &NotificationService{
		errCh: make(chan error, 1),
	}
}

func (ns *NotificationService) SendEmailChatNotification(receiverEmail string, senderEmail string, message string) error {
	errCh := make(chan error, 1)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		errCh <- helper.SendMail(receiverEmail, fmt.Sprintf("You receive message from user: %s", senderEmail), message)
	}()

	wg.Wait()
	close(errCh)
	return <-errCh
}
func (ns *NotificationService) UserForgotPassword(receiverEmail string, htmlTemplateName string, resetLink string, username string) error {
	errCh := make(chan error, 1)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		helper.ForgotPasswordMail(errCh, receiverEmail, "Reset Password URL", resetLink, username)
	}()

	wg.Wait()
	close(errCh)
	return <-errCh
}
func (ns *NotificationService) UserSucessResetPassword(receiverEmail string, htmlTemplateName string, username string) error {
	errCh := make(chan error, 1)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		helper.ResetPasswordSuccessMail(errCh, receiverEmail, "Success Reseting Your Password", username)
	}()

	wg.Wait()
	close(errCh)
	return <-errCh
}

func (ns *NotificationService) UserVerifyingEmail(receiverEmail string, htmlTemplateName string, verifyLink string) error {
	errCh := make(chan error, 1)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		helper.VerifyAccountURLMail(errCh, receiverEmail, "Verify Account URL", verifyLink)
	}()

	wg.Wait()
	close(errCh)
	return <-errCh
}
