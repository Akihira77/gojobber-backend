package handler

import (
	"context"
	"errors"
	"net/http"
	"time"

	svc "github.com/Akihira77/gojobber/services/4-user/service"
	"github.com/Akihira77/gojobber/services/4-user/types"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type BuyerHandler struct {
	buyerSvc svc.BuyerServiceImpl
}

func NewBuyerHandler(buyerSvc svc.BuyerServiceImpl) *BuyerHandler {
	return &BuyerHandler{
		buyerSvc: buyerSvc,
	}
}

func (bh *BuyerHandler) GetMyBuyerInfo(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 1*time.Second)
	defer cancel()

	currUser, ok := c.UserContext().Value("current_user").(*types.JWTClaims)
	if !ok {
		return fiber.NewError(http.StatusUnauthorized, "invalid data. Please re-signin")
	}

	myInfo, err := bh.buyerSvc.FindBuyerByEmailOrUsername(ctx, currUser.Email)
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return fiber.NewError(http.StatusNotFound, "your buyer information is not found. Please re-signin")
		}
		return fiber.NewError(http.StatusInternalServerError, "error while finding your data")
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"buyer": myInfo,
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
