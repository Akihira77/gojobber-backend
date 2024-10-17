package types

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Review struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" validate:"required"`
	SellerID  string    `json:"sellerId" gorm:"not null;" validate:"required"`
	BuyerID   string    `json:"buyerId" gorm:"not null;" validate:"required"`
	Rating    uint      `json:"rating" gorm:"not null;" validate:"required"`
	Review    string    `json:"review" gorm:"not null;" validate::"required"`
	CreatedAt time.Time `json:"createdAt" gorm:"not null;" validate:"required"`
}

type UpsertReviewDTO struct {
	SellerID string `json:"sellerId" validate:"required"`
	BuyerID  string `json:"buyerId"`
	Rating   uint   `json:"rating" validate:"required,gte=1,lte=5"`
	Review   string `json:"review" validate:"required"`
}

func ApplyDBSetup(db *gorm.DB) error {
	result := db.Exec(
		`
            ALTER TABLE reviews 
            ADD FOREIGN KEY (seller_id) REFERENCES sellers(id) ON DELETE CASCADE ON UPDATE CASCADE;
        `,
	)
	if result.Error != nil {
		return result.Error
	}

	result = db.Exec(
		`
            ALTER TABLE reviews
            ADD FOREIGN KEY (buyer_id) REFERENCES buyers(id) ON DELETE CASCADE ON UPDATE CASCADE;
        `,
	)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
