package types

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type RatingCategory struct {
	Five  uint `json:"five" gorm:"not null; default:0"`
	Four  uint `json:"four" gorm:"not null; default:0"`
	Three uint `json:"three" gorm:"not null; default:0"`
	Two   uint `json:"two" gorm:"not null; default:0"`
	One   uint `json:"one" gorm:"not null; default:0"`
}

type Gig struct {
	ID       uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	SortID   uint      `json:"sortId" gorm:"sort_id;not null;autoIncrement"`
	SellerID string    `json:"sellerId" gorm:"not null"`

	Title       string `json:"title" gorm:"not null"`
	TitleTokens string `json:"titleTokens,omitempty" gorm:"type:tsvector;column:title_tokens;"`

	Description       string `json:"description" gorm:"not null"`
	DescriptionTokens string `json:"descriptionTokens,omitempty" gorm:"type:tsvector;column:description_tokens;"`

	Category       string `json:"category" gorm:"not null"`
	CategoryTokens string `json:"categoryTokens,omitempty" gorm:"type:tsvector;column:category_tokens;"`

	SubCategories       pq.StringArray `json:"subCategories" gorm:"type:text[];not null"`
	SubCategoriesTokens string         `json:"subCategoriesTokens" gorm:"type:tsvector;column:sub_categories_tokens;"`

	Tags       pq.StringArray `json:"tags" gorm:"type:text[];not null"`
	TagsTokens string         `json:"tagsTokens" gorm:"type:tsvector;column:tags_tokens;"`

	Active               bool           `json:"active" gorm:"type:bool;default:true;not null"`
	ExpectedDeliveryDays uint           `json:"expectedDeliveryDays" gorm:"not null;"`
	RatingsCount         uint64         `json:"ratingsCount" gorm:"not null"`
	RatingSum            uint64         `json:"ratingSum" gorm:"not null;"`
	RatingCategories     RatingCategory `json:"ratingCategories" gorm:"type:jsonb;not null;serializer:json;"`
	Price                float64        `json:"price" gorm:"not null;"`
	CoverImage           string         `json:"coverImage" gorm:"not null"`
	CreatedAt            time.Time      `json:"createdAt" gorm:"not null"`
}

func ApplyDBSetup(db *gorm.DB) error {
	result := db.Debug().Exec(`
		CREATE INDEX IF NOT EXISTS idx_rating_sum
		ON gigs USING btree (rating_sum DESC);
		`)
	if result.Error != nil {
		return result.Error
	}

	result = db.Debug().Exec(`
		CREATE INDEX IF NOT EXISTS idx_price_delivery
		ON gigs USING btree (price, expected_delivery_days);
		`)
	if result.Error != nil {
		return result.Error
	}

	result = db.Debug().Exec(`
		ALTER TABLE gigs
		ADD FOREIGN KEY (seller_id) REFERENCES sellers(id) ON DELETE RESTRICT ON UPDATE CASCADE;
		`)
	if result.Error != nil {
		return result.Error
	}

	result = db.Debug().Exec(`
		CREATE INDEX IF NOT EXISTS idx_search_tokens 
		ON gigs USING gin (title_tokens, category_tokens, sub_categories_tokens, tags_tokens, description_tokens);
		`)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

type GigRankDTO struct {
	Gigs      []GigDTO `json:"gigs" gorm:"gigs"`
	RankTotal float64  `json:"rank_total"`
}

type SellerOverview struct {
	SellerID         string         `json:"sellerId"`
	FullName         string         `json:"fullName"`
	RatingsCount     uint64         `json:"ratingsCount"`
	RatingSum        uint64         `json:"ratingSum"`
	RatingCategories RatingCategory `json:"ratingCategories" gorm:"serializer:json"`
}

type GigDTO struct {
	ID                   uuid.UUID      `json:"id"`
	SellerID             string         `json:"sellerId"`
	Title                string         `json:"title"`
	Description          string         `json:"description"`
	Category             string         `json:"category"`
	SubCategories        pq.StringArray `json:"subCategories" gorm:"type:text[]"`
	Tags                 pq.StringArray `json:"tags" gorm:"type:text[]"`
	Active               bool           `json:"active"`
	ExpectedDeliveryDays int            `json:"expectedDeliveryDays"`
	RatingsCount         uint64         `json:"ratingsCount"`
	RatingSum            uint64         `json:"ratingSum"`
	RatingCategories     RatingCategory `json:"ratingCategories" gorm:"serializer:json"`
	Price                float64        `json:"price"`
	CoverImage           string         `json:"coverImage"`
	SortID               uint           `json:"sortId"`
	CreatedAt            time.Time      `json:"createdAt"`
}

type GigSellerDTO struct {
	Seller SellerOverview `json:"seller"`
	Gig    GigDTO         `json:"gig"`
}

type CreateGigDTO struct {
	SellerID             string         `json:"sellerId"`
	Title                string         `json:"title" form:"title" validate:"required"`
	Description          string         `json:"description" form:"description" validate:"required"`
	Category             string         `json:"category" form:"category" validate:"required"`
	SubCategories        pq.StringArray `json:"subCategories" form:"subCategories" validate:"required,min=1"`
	Tags                 pq.StringArray `json:"tags" form:"tags" validate:"required,min=1"`
	Active               bool           `json:"active" form:"active" validate:"required"`
	ExpectedDeliveryDays int            `json:"expectedDeliveryDays" form:"expectedDeliveryDays" validate:"required"`
	Price                float64        `json:"price" form:"price" validate:"required"`
	CoverImage           string         `json:"coverImage"`
	ImageFile            multipart.File
}

type UpdateGigDTO struct {
	SellerID             string         `json:"sellerId" form:"sellerId" validate:"required"`
	Title                string         `json:"title" form:"title" validate:"required"`
	Description          string         `json:"description" form:"description" validate:"required"`
	Category             string         `json:"category" form:"category" validate:"required"`
	SubCategories        pq.StringArray `json:"subCategories" form:"subCategories" validate:"required,min=1"`
	Tags                 pq.StringArray `json:"tags" form:"tags" validate:"required,min=1"`
	ExpectedDeliveryDays int            `json:"expectedDeliveryDays" form:"expectedDeliveryDays" validate:"required"`
	Price                float64        `json:"price" form:"price" validate:"required"`
	CoverImage           string         `json:"coverImage"`
	ImageFile            multipart.File
}

type GigSearchQueryResult struct {
	Total   int64    `json:"total"`
	Matched int64    `json:"matched"`
	Count   int      `json:"count"`
	Gigs    []GigDTO `json:"gigs"`
}

type GigSearchParams struct {
	Page int `json:"page" params:"page"`
	Size int `json:"size" params:"size"`
}

type GigSearchQuery struct {
	Query        string `json:"query" query:"query"`
	DeliveryTime int    `json:"delivery_time" query:"delivery_time"`
	Max          int    `json:"max" query:"max"`
}
