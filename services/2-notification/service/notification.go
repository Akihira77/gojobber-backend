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
	SellerHasCompletedAnOrder(receiverEmail, buyerEmail, orderID, sellerCurrentBalance string) error
	BuyerDeadlineExtensionResponse(receiverEmail, message string) error
	BuyerRefundsAnOrder(receiverEmail, reason string) error
	SellerCanceledAnOrder(receiverEmail, reason string) error
}

func NewNotificationService() NotificationServiceImpl {
	return &NotificationService{
		errCh: make(chan error, 1),
	}
}

func (ns *NotificationService) SellerCanceledAnOrder(receiverEmail string, reason string) error {
	errCh := make(chan error, 1)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		errCh <- helper.SendMail(receiverEmail, fmt.Sprintf("Seller Has Canceled Your Order"), reason)
	}()

	wg.Wait()
	close(errCh)
	return <-errCh
}

// TODO: REFACTORE
func (ns *NotificationService) BuyerDeadlineExtensionResponse(receiverEmail string, message string) error {
	errCh := make(chan error, 1)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		errCh <- helper.SendMail(receiverEmail, fmt.Sprintf("Buyer Response Your Deadline Extension"), message)
	}()

	wg.Wait()
	close(errCh)
	return <-errCh
}

// TODO: REFACTORE
func (ns *NotificationService) BuyerRefundsAnOrder(receiverEmail string, reason string) error {
	errCh := make(chan error, 1)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		errCh <- helper.SendMail(receiverEmail, fmt.Sprintf("Buyer Refunds The Order"), reason)
	}()

	wg.Wait()
	close(errCh)
	return <-errCh
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

func (ns *NotificationService) SellerHasCompletedAnOrder(receiverEmail, buyerEmail, orderID, sellerCurrentBalance string) error {
	errCh := make(chan error, 1)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		helper.SellerOrderHasCompleted(errCh, receiverEmail, fmt.Sprintf("Your Order [%s] Has Marked As Complete By Buyer %s", orderID, buyerEmail), buyerEmail, orderID, sellerCurrentBalance)
	}()

	wg.Wait()
	close(errCh)
	return <-errCh
}
