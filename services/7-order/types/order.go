package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type orderStatus string

const (
	PENDING   orderStatus = "PENDING"
	PROCESS   orderStatus = "PROCESS"
	CANCELED  orderStatus = "CANCELED"
	REFUNDED  orderStatus = "REFUNDED"
	COMPLETED orderStatus = "COMPLETED"
)

func (p *orderStatus) Scan(value interface{}) error {
	*p = orderStatus(value.([]byte))
	return nil
}

func (p orderStatus) Value() (driver.Value, error) {
	return string(p), nil
}

type OrderEvent struct {
	Event     string    `json:"event" gorm:"not null;"`
	CreatedAt time.Time `json:"createdAt" gorm:"not null;"`
}
type OrderEvents []OrderEvent

func (o OrderEvents) Value() (driver.Value, error) {
	return json.Marshal(o)
}
func (o *OrderEvents) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value")
	}

	return json.Unmarshal(bytes, &o)
}

type DeliveredHistory struct {
	ID            uuid.UUID `json:"id" gorm:"primaryKey"`
	OrderID       string    `json:"orderId"`
	ProgressNote  string    `json:"progressNote" gorm:"not null;"`
	ResultURL     string    `json:"resultUrl" gorm:""`
	BuyerNote     string    `json:"buyerNote" gorm:""`
	DeliveredDate time.Time `json:"deliveredDate" gorm:"not null;"`
}

type Order struct {
	ID                 string             `json:"id" gorm:"primaryKey"`
	SellerID           string             `json:"sellerId" gorm:"not null;"`
	BuyerID            string             `json:"buyerId" gorm:"not null;"`
	GigTitle           string             `json:"gigTitle" gorm:"not null;"`
	GigDescription     string             `json:"gigDescription" gorm:"not null;"`
	Status             orderStatus        `json:"status" gorm:"type:order_status;not null;"`
	Price              uint64             `json:"price" gorm:"not null;"`
	ServiceFee         uint               `json:"serviceFee" gorm:"not null; default:0;"`
	PaymentIntent      string             `json:"paymentIntent" gorm:"not null;"`
	DeliveredHistories []DeliveredHistory `json:"deliveredHistories" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	OrderEvents        OrderEvents        `json:"orderEvents" gorm:"type:jsonb; serializer:json;"`
	InvoiceID          string             `json:"invoiceId,omitempty"`
	StartDate          time.Time          `json:"startDate" gorm:"not null;"`
	Deadline           time.Time          `json:"deadline" gorm:"not null;"`
}

func ApplyDBSetup(db *gorm.DB) error {
	err := db.Exec(
		`
		ALTER TABLE orders
		ADD FOREIGN KEY (seller_id) REFERENCES sellers(id) ON DELETE RESTRICT ON UPDATE CASCADE;
		`,
	).Error
	if err != nil {
		return err
	}

	err = db.Exec(
		`
		ALTER TABLE orders
		ADD FOREIGN KEY (buyer_id) REFERENCES buyers(id) ON DELETE RESTRICT ON UPDATE CASCADE;
		`,
	).Error
	if err != nil {
		return err
	}

	err = db.Exec(
		`
		CREATE INDEX IF NOT EXISTS idx_seller_id
		ON orders USING btree (seller_id);
		`,
	).Error
	if err != nil {
		return err
	}

	err = db.Exec(
		`
		CREATE INDEX IF NOT EXISTS idx_buyer_id
		ON orders USING btree (buyer_id);
		`,
	).Error
	if err != nil {
		return err
	}
	return nil
}
