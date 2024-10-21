package handler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Akihira77/gojobber/services/6-chat/service"
	"github.com/Akihira77/gojobber/services/6-chat/types"
	"github.com/Akihira77/gojobber/services/6-chat/util"
	"github.com/Akihira77/gojobber/services/common/genproto/notification"
	"github.com/Akihira77/gojobber/services/common/genproto/user"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ChatHandler struct {
	cs         service.ChatServiceImpl
	validate   *validator.Validate
	cld        *util.Cloudinary
	grpcClient *GRPCClients
}

func NewChatHandler(cld *util.Cloudinary, cs service.ChatServiceImpl, grpcServices *GRPCClients) *ChatHandler {
	return &ChatHandler{
		cs:         cs,
		validate:   validator.New(validator.WithRequiredStructEnabled()),
		cld:        cld,
		grpcClient: grpcServices,
	}
}

func (ch *ChatHandler) GetAllMyConversations(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 1*time.Second)
	defer cancel()

	userInfo, ok := c.UserContext().Value("current_user").(*types.JWTClaims)
	if !ok {
		return fiber.NewError(http.StatusUnauthorized, "signin first")
	}

	conversations, err := ch.cs.GetAllMyConversations(ctx, userInfo.UserID)
	if err != nil {
		fmt.Printf("Get All My Conversations Error:\n+%v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(http.StatusOK).JSON(fiber.Map{
				"conversations": conversations,
			})
		}
		return fiber.NewError(http.StatusInternalServerError, "Error while searching conversations data")
	}

	cc, err := ch.grpcClient.GetClient("USER_SERVICE")
	if err != nil {
		fmt.Printf("Get All My Conversations Error:\n+%v", err)
		return fiber.NewError(http.StatusInternalServerError, "Error while searching conversations data")
	}

	var wg sync.WaitGroup
	errCh := make(chan error, len(conversations))
	userGrpcClient := user.NewUserServiceClient(cc)

	for i := range conversations {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			newCtx, canc := context.WithTimeout(ctx, 200*time.Millisecond)
			defer canc()

			userId := conversations[idx].UserOneID
			if userId == userInfo.UserID {
				userId = conversations[idx].UserTwoID
			}

			b, err := userGrpcClient.FindBuyer(newCtx, &user.FindBuyerRequest{
				BuyerId: userId,
			})
			if err != nil {
				errCh <- err
				return
			}

			conversations[idx].SenderName = b.Username
			conversations[idx].SenderEmail = b.Email
			conversations[idx].SenderProfilePicture = b.ProfilePicture

			errCh <- nil
		}(i)

	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		if err != nil {
			fmt.Printf("Get All My Conversations Error:\n+%v", err)
			return fiber.NewError(http.StatusInternalServerError, "Error while searching conversations data")
		}
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"conversations": conversations,
		"total":         len(conversations),
	})
}

func (ch *ChatHandler) GetMessagesInsideConversation(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 1*time.Second)
	defer cancel()

	chat, err := ch.cs.GetMessages(ctx, c.Params("conversationId"))
	if err != nil {
		fmt.Printf("Get All My Conversations Error:\n+%v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(http.StatusNotFound, "Oops, your chat data does not found")
		}
		return fiber.NewError(http.StatusInternalServerError, "Error while searching conversations data")
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"conversation": chat,
		"total":        len(chat),
	})
}

func (ch *ChatHandler) InsertMessage(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	userInfo, ok := c.UserContext().Value("current_user").(*types.JWTClaims)
	if !ok {
		log.Println(userInfo)
		return fiber.NewError(http.StatusUnauthorized, "Sender data is invalid")
	}

	data := new(types.CreateMessageDTO)
	err := c.BodyParser(data)
	if err != nil {
		fmt.Printf("InsertMessage Error:\n+%v", err)
		return fiber.NewError(http.StatusBadRequest, "Error reading request body")
	}

	cc, err := ch.grpcClient.GetClient("USER_SERVICE")
	userGrpcClient := user.NewUserServiceClient(cc)
	if err != nil {
		fmt.Printf("InsertMessage Error:\n%+v", err)
		return fiber.NewError(http.StatusBadRequest, "Error saving message")
	}

	if data.Offer != nil {
		_, err = userGrpcClient.FindSeller(ctx, &user.FindSellerRequest{
			SellerId: userInfo.UserID,
		})
		if err != nil {
			fmt.Printf("InsertMessage Error:\n+%v", err)
			return fiber.NewError(http.StatusBadRequest, "Non Seller User Can Make An Offer To Another User")
		}
	}

	receiverUser, err := userGrpcClient.FindBuyer(ctx, &user.FindBuyerRequest{
		BuyerId: data.ReceiverID,
	})
	if err != nil {
		fmt.Printf("InsertMessage Error:\n+%v", err)
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	data.ReceiverEmail = receiverUser.Email
	if err := ch.validate.Struct(data); err != nil {
		fmt.Printf("InsertMessage Error:\n+%v", err)
		return fiber.NewError(http.StatusBadRequest, "Error validating request body")
	}

	formHeader, err := c.FormFile("file")
	if err != nil {
		fmt.Printf("InsertMessage Error:\n+%v", err)
		if data.Body == "" {
			return fiber.NewError(http.StatusBadRequest, "Can't send chat with empty body and file")
		}
	} else {
		data.File, err = formHeader.Open()
		if err != nil {
			fmt.Printf("InsertMessage Error:\n+%v", err)
			return fiber.NewError(http.StatusBadRequest, "Error reading file")
		}

		filePath := util.RandomStr(64)
		//FIX: IDK IF IT IS CORRECT HEADER
		uploadResult, err := ch.cld.UploadFile(ctx, formHeader, data.File, filePath, formHeader.Header.Get("file-type"))
		if err != nil {
			fmt.Printf("InsertMessage Error:\n+%v", err)
			return fiber.NewError(http.StatusInternalServerError, "Error processing file")
		}

		data.FileURL = uploadResult.SecureURL
	}

	chat, err := ch.cs.InsertMessage(ctx, userInfo.UserID, data)
	if err != nil {
		fmt.Printf("InsertMessage Error:\n+%v", err)
		return fiber.NewError(http.StatusInternalServerError, "Error saving your chat")
	}

	go func() {
		if chat.Offer == nil {
			return
		}

		cc, err = ch.grpcClient.GetClient("NOTIFICATION_SERVICE")
		if err != nil {
			fmt.Printf("InsertMessage Error:\n+%v", err)
			return
		}

		notificationGrpcClient := notification.NewNotificationServiceClient(cc)
		_, err = notificationGrpcClient.SendEmailChatNotification(context.TODO(), &notification.EmailChatNotificationRequest{
			ReceiverEmail: data.ReceiverEmail,
			SenderEmail:   userInfo.Email,
			Message:       fmt.Sprintf("You receive a Gig Offer from seller: %s", userInfo.Email),
		})
		if err != nil {
			fmt.Printf("InsertMessage Error:\n+%v", err)
			return
		}
	}()

	unreadMessages := ch.cs.CalculateUnreadMessages(ctx, chat.ConversationID, userInfo.UserID)

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"senderId":       userInfo.UserID,
		"receiver":       receiverUser,
		"unreadMessages": unreadMessages,
		"chat":           chat,
	})
}

func (ch *ChatHandler) SellerCancelOffer(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	userInfo, ok := c.UserContext().Value("current_user").(*types.JWTClaims)
	if !ok {
		log.Println(userInfo)
		return fiber.NewError(http.StatusUnauthorized, "Sender data is invalid")
	}

	m, err := ch.cs.FindMessageByID(ctx, c.Params("messageId"))
	if err != nil {
		fmt.Printf("SellerCancelOffer Error:\n+%v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(http.StatusNotFound, "Message is not found")
		}
		return fiber.NewError(http.StatusInternalServerError, "Unexpected error happened. Please try again.")
	}

	err = ch.cs.ChangeOfferStatus(ctx, m, types.CANCELED)
	if err != nil {
		fmt.Printf("SellerCancelOffer Error:\n+%v", err)
		return fiber.NewError(http.StatusInternalServerError, "Unexpected error happened. Please try again.")
	}

	return c.SendStatus(http.StatusOK)
}
