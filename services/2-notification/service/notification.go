package service

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/Akihira77/gojobber/services/2-notification/helper"
	"github.com/Akihira77/gojobber/services/common/genproto/notification"
)

type NotificationService struct {
	errCh chan error
}

type NotificationServiceImpl interface {
	UserVerifyingEmail(receiverEmail, htmlTemplateName, verifyLink string) error
	UserForgotPassword(receiverEmail, htmlTemplateName, resetLink, username string) error
	UserSucessResetPassword(receiverEmail, htmlTemplateName, username string) error
	SendEmailChatNotification(receiverEmail, senderEmail, message string) error
	SellerHasCompletedAnOrder(data *notification.SellerCompletedAnOrderRequest) error
	SellerRequestDeadlineExtension(data *notification.SellerDeadlineExtensionRequest) error
	BuyerDeadlineExtensionResponse(data *notification.BuyerDeadlineExtension) error
	BuyerRefundsAnOrder(data *notification.BuyerRefundsOrderRequest) error
	SellerCanceledAnOrder(data *notification.SellerCancelOrderRequest) error
	NotifySellerGotAnOrder(data *notification.NotifySellerGotAnOrderRequest) error
	NotifySellerGotAReview(data *notification.NotifySellerGotAReviewRequest) error
	NotifyBuyerSellerDeliveredOrder(data *notification.NotifyBuyerOrderDeliveredRequest) error
	NotifyBuyerSellerProcessedOrder(data *notification.NotifyBuyerOrderAcknowledgeRequest) error
}

func NewNotificationService() NotificationServiceImpl {
	return &NotificationService{
		errCh: make(chan error, 1),
	}
}

// TODO: IMPLEMENT HTML TEMPLATE
func (ns *NotificationService) NotifyBuyerSellerProcessedOrder(data *notification.NotifyBuyerOrderAcknowledgeRequest) error {
	errCh := make(chan error, 1)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		errCh <- helper.
			SendMail(
				data.ReceiverEmail,
				fmt.Sprint("Seller Has Acknowledge Your Order And Start Working On It"),
				fmt.Sprintf("Check your order %s", data.Url),
			)
	}()

	wg.Wait()
	close(errCh)
	return <-errCh
}

// TODO: IMPLEMENT HTML TEMPLATE
// AND REFACTORE request payload to ADD THE URL TO OUR PLATFORM
func (ns *NotificationService) NotifyBuyerSellerDeliveredOrder(data *notification.NotifyBuyerOrderDeliveredRequest) error {
	errCh := make(chan error, 1)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		errCh <- helper.SendMail(
			data.ReceiverEmail,
			fmt.Sprintf("Seller Has Sent Your Order Progress. Check Out Your Order!"),
			fmt.Sprintf("Check your order %s", data.Url),
		)
	}()

	wg.Wait()
	close(errCh)
	return <-errCh
}

// TODO: IMPLEMENT HTML TEMPLATE
func (ns *NotificationService) NotifySellerGotAReview(data *notification.NotifySellerGotAReviewRequest) error {
	errCh := make(chan error, 1)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		errCh <- helper.SendMail(data.ReceiverEmail, fmt.Sprintf("User Giving You Review"), data.Message)
	}()

	wg.Wait()
	close(errCh)
	return <-errCh
}

// TODO: IMPLEMENT HTML TEMPLATE
func (ns *NotificationService) NotifySellerGotAnOrder(data *notification.NotifySellerGotAnOrderRequest) error {
	b, err := json.Marshal(data.Detail)
	if err != nil {
		return fmt.Errorf("Error parsing payload data %v", err)
	}

	errCh := make(chan error, 1)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		errCh <- helper.SendMail(data.ReceiverEmail, data.Message, string(b))
	}()

	wg.Wait()
	close(errCh)
	return <-errCh
}

func (ns *NotificationService) SellerCanceledAnOrder(data *notification.SellerCancelOrderRequest) error {
	errCh := make(chan error, 1)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		errCh <- helper.SendMail(
			data.ReceiverEmail,
			fmt.Sprintf("Seller Has Canceled Your Order"),
			fmt.Sprintf("Check your order %s", data.Url),
		)
	}()

	wg.Wait()
	close(errCh)
	return <-errCh
}

// TODO: REFACTORE
func (ns *NotificationService) SellerRequestDeadlineExtension(data *notification.SellerDeadlineExtensionRequest) error {
	errCh := make(chan error, 1)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		errCh <- helper.SendMail(
			data.ReceiverEmail,
			fmt.Sprintf("Seller Requested A Deadline Extension"),
			fmt.Sprintf("Check your order %s", data.Url),
		)
	}()

	wg.Wait()
	close(errCh)
	return <-errCh
}

// TODO: REFACTORE
func (ns *NotificationService) BuyerDeadlineExtensionResponse(data *notification.BuyerDeadlineExtension) error {
	errCh := make(chan error, 1)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		errCh <- helper.SendMail(
			data.ReceiverEmail,
			fmt.Sprintf("Buyer Response Your Deadline Extension"),
			fmt.Sprintf("Check your order %s", data.Url),
		)
	}()

	wg.Wait()
	close(errCh)
	return <-errCh
}

// TODO: REFACTORE
func (ns *NotificationService) BuyerRefundsAnOrder(data *notification.BuyerRefundsOrderRequest) error {
	errCh := make(chan error, 1)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		errCh <- helper.SendMail(
			data.ReceiverEmail,
			fmt.Sprintf("Buyer Refunds The Order"),
			fmt.Sprintf("Check your order %s", data.Url),
		)
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

func (ns *NotificationService) SellerHasCompletedAnOrder(data *notification.SellerCompletedAnOrderRequest) error {
	errCh := make(chan error, 1)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		// helper.SellerOrderHasCompleted(errCh, receiverEmail, fmt.Sprintf("Your Order [%s] Has Marked As Complete By Buyer %s", orderID, buyerEmail), buyerEmail, orderID, sellerCurrentBalance)

		errCh <- helper.SendMail(
			data.ReceiverEmail,
			fmt.Sprintf("Buyer [%s] Mark Your Order [%s] As COMPLETED", data.BuyerEmail, data.OrderId),
			fmt.Sprintf("Your Current Balance is: %s.\nCheck Your Order %s", data.SellerCurrentBalance, data.Url),
		)

	}()

	wg.Wait()
	close(errCh)
	return <-errCh
}
