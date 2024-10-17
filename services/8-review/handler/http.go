package handler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Akihira77/gojobber/services/8-review/service"
	"github.com/Akihira77/gojobber/services/8-review/types"
	"github.com/Akihira77/gojobber/services/common/genproto/notification"
	"github.com/Akihira77/gojobber/services/common/genproto/user"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ReviewHandler struct {
	reviewSvc    service.ReviewServiceImpl
	validate     *validator.Validate
	grpcServices *GRPCClients
}

func NewReviewHandler(reviewSvc service.ReviewServiceImpl) *ReviewHandler {
	return &ReviewHandler{
		reviewSvc:    reviewSvc,
		validate:     validator.New(validator.WithRequiredStructEnabled()),
		grpcServices: NewGRPCClients(),
	}
}

func (rh *ReviewHandler) FindSellerReviews(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 500*time.Millisecond)
	defer cancel()

	reviews, err := rh.reviewSvc.FindSellerReviews(ctx, c.Params("sellerId"))
	if err != nil {
		log.Printf("FindSellerReviews Error:\n+%v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"reviews": []types.Review{},
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"reviews": reviews,
	})
}

func (rh *ReviewHandler) Add(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 500*time.Millisecond)
	defer cancel()

	userInfo, ok := c.UserContext().Value("current_user").(*types.JWTClaims)
	if !ok {
		log.Println(userInfo)
		return fiber.NewError(http.StatusUnauthorized, "Sign-in first")
	}

	data := new(types.UpsertReviewDTO)
	err := c.BodyParser(data)
	if err != nil {
		log.Printf("Add Review Error:\n+%v", err)
		return fiber.NewError(http.StatusBadRequest, "Invalid Data")
	}

	err = rh.validate.Struct(data)
	if err != nil {
		log.Printf("Add Review Error:\n+%v", err)
		return fiber.NewError(http.StatusBadRequest, "Invalid Data")
	}

	cc, err := rh.grpcServices.GetClient("USER_SERVICE")
	if err != nil {
		log.Printf("Add Review Error:\n+%v", err)
		return fiber.NewError(http.StatusNotFound, "Error while validating seller")
	}

	userGrpcClient := user.NewUserServiceClient(cc)
	s, err := userGrpcClient.FindSeller(ctx, &user.FindSellerRequest{
		SellerId: data.SellerID,
	})
	if err != nil {
		log.Printf("Add Review Error:\n+%v", err)
		return fiber.NewError(http.StatusNotFound, "Error while validating seller")
	}

	r, err := rh.reviewSvc.Add(ctx, *data)
	if err != nil {
		log.Printf("Add Review Error:\n+%v", err)
		return fiber.NewError(http.StatusBadRequest, "Error while saving review")
	}

	go func() {
		cc, err := rh.grpcServices.GetClient("NOTIFICATION_SERVICE")
		if err != nil {
			log.Printf("Add Review Error:\n+%v", err)
			return
		}

		notificationGrpcClient := notification.NewNotificationServiceClient(cc)
		_, err = notificationGrpcClient.NotifySellerGotAReview(context.TODO(), &notification.NotifySellerGotAReviewRequest{
			ReceiverEmail: s.Email,
			Message:       fmt.Sprintf("Buyer [%s] Giving You A Rating [%v] And Review:\n%s", data.BuyerID, data.Rating, data.Review),
		})
		if err != nil {
			log.Printf("Add Review Error:\n+%v", err)
			return
		}
	}()

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"review": r,
	})
}

func (rh *ReviewHandler) Update(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 500*time.Millisecond)
	defer cancel()

	userInfo, ok := c.UserContext().Value("current_user").(*types.JWTClaims)
	if !ok {
		log.Println(userInfo)
		return fiber.NewError(http.StatusUnauthorized, "Sign-in first")
	}

	data := new(types.Review)
	err := c.BodyParser(data)
	if err != nil {
		log.Printf("Update Review Error:\n+%v", err)
		return fiber.NewError(http.StatusBadRequest, "Invalid Data")
	}

	err = rh.validate.Struct(data)
	if err != nil {
		log.Printf("Update Review Error:\n+%v", err)
		return fiber.NewError(http.StatusBadRequest, "Invalid Data")
	}

	r, err := rh.reviewSvc.FindReviewByID(ctx, c.Params("reviewId"))
	if err != nil {
		log.Printf("Update Review Error:\n+%v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(http.StatusNotFound, "Review is not found")
		}
		return fiber.NewError(http.StatusBadRequest, "Error while saving review")
	}

	r, err = rh.reviewSvc.Update(ctx, *data)
	if err != nil {
		log.Printf("Update Review Error:\n+%v", err)
		return fiber.NewError(http.StatusBadRequest, "Error while updating review")
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"review": r,
	})
}

func (rh *ReviewHandler) Remove(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 500*time.Millisecond)
	defer cancel()

	userInfo, ok := c.UserContext().Value("current_user").(*types.JWTClaims)
	if !ok {
		log.Println(userInfo)
		return fiber.NewError(http.StatusUnauthorized, "Sign-in first")
	}

	_, err := rh.reviewSvc.FindReviewByID(ctx, c.Params("reviewId"))
	if err != nil {
		log.Printf("Update Review Error:\n+%v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(http.StatusNotFound, "Review is not found")
		}
		return fiber.NewError(http.StatusBadRequest, "Error while searching review")
	}

	err = rh.reviewSvc.Remove(ctx, c.Params("reviewId"))
	if err != nil {
		log.Printf("Update Review Error:\n+%v", err)
		return fiber.NewError(http.StatusBadRequest, "Error while removing review")
	}

	return c.SendStatus(http.StatusOK)
}
