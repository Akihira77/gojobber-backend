package service

import (
	"context"
	"time"

	"github.com/Akihira77/gojobber/services/8-review/types"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ReviewService struct {
	db *gorm.DB
}

type ReviewServiceImpl interface {
	FindSellerReviews(ctx context.Context, sellerID string) ([]types.Review, error)
	FindReviewByID(ctx context.Context, id uuid.UUID) (*types.Review, error)
	FindReviewByIDAndUserID(ctx context.Context, id uuid.UUID, userId string) (*types.Review, error)
	Add(ctx context.Context, data types.UpsertReviewDTO) (*types.Review, error)
	Update(ctx context.Context, data types.Review) (*types.Review, error)
	Remove(ctx context.Context, reviewID string) error
}

func NewReviewService(db *gorm.DB) ReviewServiceImpl {
	return &ReviewService{
		db: db,
	}
}

func (rs *ReviewService) FindSellerReviews(ctx context.Context, sellerID string) ([]types.Review, error) {
	var rvs []types.Review

	result := rs.db.
		Debug().
		WithContext(ctx).
		Model(&types.Review{}).
		Where("seller_id = ?", sellerID).
		Find(&rvs)

	return rvs, result.Error
}

func (rs *ReviewService) FindReviewByIDAndUserID(ctx context.Context, id uuid.UUID, userId string) (*types.Review, error) {
	var r types.Review

	result := rs.db.
		Debug().
		WithContext(ctx).
		Model(&types.Review{}).
		Where("id = ? AND buyer_id = ?", id, userId).
		First(&r)

	return &r, result.Error
}

func (rs *ReviewService) FindReviewByID(ctx context.Context, id uuid.UUID) (*types.Review, error) {
	var r types.Review

	result := rs.db.
		Debug().
		WithContext(ctx).
		Model(&types.Review{}).
		Where("id = ?", id).
		First(&r)

	return &r, result.Error
}

func (rs *ReviewService) Add(ctx context.Context, data types.UpsertReviewDTO) (*types.Review, error) {
	r := types.Review{
		Rating:    data.Rating,
		Review:    data.Review,
		SellerID:  data.SellerID,
		BuyerID:   data.BuyerID,
		CreatedAt: time.Now(),
	}
	result := rs.db.
		Debug().
		WithContext(ctx).
		Model(&types.Review{}).
		Create(&r)

	return &r, result.Error
}

func (rs *ReviewService) Remove(ctx context.Context, reviewID string) error {
	result := rs.db.
		Debug().
		WithContext(ctx).
		Model(&types.Review{}).
		Where("id = ?", reviewID).
		Delete(&types.Review{})

	return result.Error
}

func (rs *ReviewService) Update(ctx context.Context, data types.Review) (*types.Review, error) {
	result := rs.db.
		Debug().
		WithContext(ctx).
		Model(&data).
		Clauses(clause.Returning{}).
		Updates(types.Review{Rating: data.Rating, Review: data.Review})

	return &data, result.Error
}
