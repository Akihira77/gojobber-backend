package service

import (
	"context"

	"github.com/Akihira77/gojobber/services/4-user/types"
	"gorm.io/gorm"
)

type BuyerService struct {
	db *gorm.DB
}

type BuyerServiceImpl interface {
	FindBuyerByID(ctx context.Context, id string) (*types.Buyer, error)
	FindBuyerByEmailOrUsername(ctx context.Context, str string) (*types.Buyer, error)
	Create(ctx context.Context, b types.Buyer) error
	Delete(ctx context.Context, userId string) error
}

func NewBuyerService(db *gorm.DB) BuyerServiceImpl {
	return &BuyerService{
		db: db,
	}
}

func (bs *BuyerService) FindBuyerByID(ctx context.Context, id string) (*types.Buyer, error) {
	var buyer types.Buyer
	result := bs.db.
		WithContext(ctx).
		Model(&types.Buyer{}).
		First(&buyer, "id = ?", id)

	return &buyer, result.Error
}
func (bs *BuyerService) FindBuyerByEmailOrUsername(ctx context.Context, str string) (*types.Buyer, error) {
	var buyer types.Buyer
	result := bs.db.
		WithContext(ctx).
		Model(&types.Buyer{}).
		First(&buyer, "email = ? OR username = ?", str, str)

	return &buyer, result.Error
}

func (bs *BuyerService) Create(ctx context.Context, b types.Buyer) error {
	return bs.db.
		WithContext(ctx).
		Model(&types.Buyer{}).
		Create(&b).Error
}

func (bs *BuyerService) Delete(ctx context.Context, userId string) error {
	return bs.db.
		WithContext(ctx).
		Model(&types.Buyer{}).
		Where("id = ?", userId).
		Delete(&types.Buyer{}).
		Error
}
