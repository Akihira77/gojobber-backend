package handler

import (
	"context"
	"errors"
	"net/http"
	"time"

	svc "github.com/Akihira77/gojobber/services/4-user/service"
	"github.com/Akihira77/gojobber/services/4-user/types"
	"github.com/Akihira77/gojobber/services/4-user/util"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type BuyerHandler struct {
	buyerSvc svc.BuyerServiceImpl
	validate *validator.Validate
}

func NewBuyerHandler(buyerSvc svc.BuyerServiceImpl) *BuyerHandler {
	return &BuyerHandler{
		buyerSvc: buyerSvc,
		validate: validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (bh *BuyerHandler) GetMyBuyerInfo(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 1*time.Second)
	defer cancel()

	userInfo, ok := c.UserContext().Value("current_user").(*types.JWTClaims)
	if !ok {
		return fiber.NewError(http.StatusUnauthorized, "invalid data. Please re-signin")
	}

	u, err := bh.buyerSvc.FindBuyerByID(ctx, userInfo.UserID)
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return fiber.NewError(http.StatusNotFound, "your buyer information is not found. Please re-signin")
		}
		return fiber.NewError(http.StatusInternalServerError, "error while finding your data")
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"buyer": u,
	})
}

func (bh *BuyerHandler) FindBuyerByUsername(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 1*time.Second)
	defer cancel()

	buyer, err := bh.buyerSvc.FindBuyerByEmailOrUsername(ctx, c.Params("username"))
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return fiber.NewError(http.StatusNotFound, "your buyer information is not found. Please re-signin")
		}
		return fiber.NewError(http.StatusInternalServerError, "error while finding your data")
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"buyer": buyer,
	})
}

func (bh *BuyerHandler) Update(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 1*time.Second)
	defer cancel()

	userInfo, ok := c.UserContext().Value("current_user").(*types.JWTClaims)
	if !ok {
		return fiber.NewError(http.StatusUnauthorized, "invalid data. Please re-signin")
	}

	if !userInfo.VerifiedUser {
		return fiber.NewError(http.StatusForbidden, "Verify Your Email First")
	}

	u, err := bh.buyerSvc.FindBuyerByID(ctx, userInfo.UserID)
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return fiber.NewError(http.StatusNotFound, "your buyer information is not found. Please re-signin")
		}
		return fiber.NewError(http.StatusInternalServerError, "error while finding your data")
	}

	data := new(types.EditBuyerDTO)
	err = c.BodyParser(data)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "Invalid Provided Data")
	}

	err = bh.validate.Struct(data)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors": util.CustomValidationErrors(err),
		})
	}

	u, err = bh.buyerSvc.Update(ctx, *u, data)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Erro while saving your data")
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"buyer": u,
	})

}
