package handler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	svc "github.com/Akihira77/gojobber/services/4-user/service"
	"github.com/Akihira77/gojobber/services/4-user/types"
	"github.com/Akihira77/gojobber/services/4-user/util"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v80"
	"github.com/stripe/stripe-go/v80/account"
	"gorm.io/gorm"
)

type SellerHandler struct {
	sellerSvc svc.SellerServiceImpl
	buyerSvc  svc.BuyerServiceImpl
	validate  *validator.Validate
	stripeKey string
}

func NewSellerHandler(buyerSvc svc.BuyerServiceImpl, sellerSvc svc.SellerServiceImpl) *SellerHandler {
	return &SellerHandler{
		sellerSvc: sellerSvc,
		buyerSvc:  buyerSvc,
		validate:  validator.New(validator.WithRequiredStructEnabled()),
		stripeKey: os.Getenv("STRIPE_SECRET_KEY"),
	}
}

func (sh *SellerHandler) GetMySellerInfo(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 1*time.Second)
	defer cancel()

	userInfo, ok := c.UserContext().Value("current_user").(*types.JWTClaims)
	if !ok {
		log.Println("invalid userInfo", userInfo)
		return fiber.NewError(http.StatusUnauthorized, "sign-in first")
	}

	seller, err := sh.sellerSvc.FindSellerByUsername(ctx, userInfo.Username)
	if err != nil {
		log.Println("get my seller info", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(http.StatusNotFound, "seller data is not found")
		}
		return fiber.NewError(http.StatusInternalServerError, "Error while finding your data")
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"seller": seller,
	})
}

func (sh *SellerHandler) FindSellerByID(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 1*time.Second)
	defer cancel()

	seller, err := sh.sellerSvc.FindSellerByID(ctx, c.Params("id"))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(http.StatusNotFound, "data is not found")
		}

		return fiber.NewError(http.StatusInternalServerError, "error while searching data.")
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"seller": seller,
	})
}

func (sh *SellerHandler) FindSellerByUsername(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 1*time.Second)
	defer cancel()

	seller, err := sh.sellerSvc.FindSellerByUsername(ctx, c.Params("username"))
	if err != nil {
		fmt.Printf("findSellerByUsername error:\n%+v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(http.StatusNotFound, "data is not found")
		}

		return fiber.NewError(http.StatusInternalServerError, "error while searching data.")
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"seller": seller,
	})
}

func (sh *SellerHandler) GetRandomSellers(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 1*time.Second)
	defer cancel()

	value, err := strconv.Atoi(c.Params("count"))
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid sample count number")
	}

	sellers, err := sh.sellerSvc.GetRandomSellers(ctx, value)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"total":   len(sellers),
		"sellers": sellers,
	})
}

func (sh *SellerHandler) Create(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	userInfo, ok := c.UserContext().Value("current_user").(*types.JWTClaims)
	if !ok {
		return fiber.NewError(http.StatusBadRequest, "invalid data. Please re-signin")
	}

	if !userInfo.VerifiedUser {
		return fiber.NewError(http.StatusForbidden, "Verify Your Email First")
	}

	data := new(types.CreateSellerDTO)
	if err := c.BodyParser(data); err != nil {
		fmt.Printf("%+v", err)
		return fiber.NewError(http.StatusBadRequest, "invalid data. Please check again your data")
	}

	err := sh.validate.Struct(data)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors": util.CustomValidationErrors(err),
		})
	}

	sellerDataInBuyerDB, err := sh.buyerSvc.FindBuyerByID(ctx, userInfo.UserID)
	if err != nil {
		fmt.Printf("%+v", err)
		return fiber.NewError(http.StatusNotFound, "data does not found in buyer database")
	}

	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	acctParams := &stripe.AccountParams{
		Country:      stripe.String(util.FindMyCountryCode(sellerDataInBuyerDB.Country)),
		Email:        &sellerDataInBuyerDB.Email,
		BusinessType: stripe.String(string(stripe.AccountBusinessTypeIndividual)),
		Controller: &stripe.AccountControllerParams{
			StripeDashboard: &stripe.AccountControllerStripeDashboardParams{
				Type: stripe.String("none"),
			},
			Fees: &stripe.AccountControllerFeesParams{
				Payer: stripe.String("application"),
			},
			Losses: &stripe.AccountControllerLossesParams{
				Payments: stripe.String("application"),
			},
			RequirementCollection: stripe.String("application"),
		}, Capabilities: &stripe.AccountCapabilitiesParams{
			CardPayments: &stripe.AccountCapabilitiesCardPaymentsParams{
				Requested: stripe.Bool(true),
			},
			Transfers: &stripe.AccountCapabilitiesTransfersParams{
				Requested: stripe.Bool(true),
			},
		},
	}
	acc, err := account.New(acctParams)
	if err != nil {
		fmt.Printf("%+v", err)
		return fiber.NewError(http.StatusInternalServerError, "Error while creating seller account")
	}

	data.StripeAccountID = acc.ID
	seller, err := sh.sellerSvc.Create(ctx, sellerDataInBuyerDB, data)
	if err != nil {
		fmt.Printf("%+v", err)
		return fiber.NewError(http.StatusBadRequest, "error saving your data")
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"seller": seller,
	})
}

func (sh *SellerHandler) Update(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	userInfo, ok := c.UserContext().Value("current_user").(*types.JWTClaims)
	if !ok {
		return fiber.NewError(http.StatusBadRequest, "invalid data. Please re-signin")
	}

	if !userInfo.VerifiedUser {
		return fiber.NewError(http.StatusForbidden, "Verify Your Email First")
	}

	data := new(types.UpdateSellerDTO)
	if err := c.BodyParser(data); err != nil {
		fmt.Printf("%+v", err)
		return fiber.NewError(http.StatusBadRequest, "invalid data. Please check again your data")
	}

	err := sh.validate.Struct(data)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors": util.CustomValidationErrors(err),
		})
	}

	seller, err := sh.sellerSvc.FindSellerByBuyerID(ctx, userInfo.UserID)
	if err != nil {
		fmt.Printf("%+v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(http.StatusNotFound, "your data is not found")
		}

		return fiber.NewError(http.StatusInternalServerError, "error while searching data")
	}

	err = sh.sellerSvc.Update(ctx, seller, data)
	if err != nil {
		fmt.Printf("%+v", err)
		return fiber.NewError(http.StatusInternalServerError, "update data error. Try again")
	}

	return c.Status(http.StatusOK).SendString("update success")
}
