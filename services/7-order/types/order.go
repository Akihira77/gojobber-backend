package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderStatus string

const (
	REQUIRE_PAYMENT OrderStatus = "REQUIRE PAYMENT" //NOTE: BUYER HAS NOT PAID
	PENDING         OrderStatus = "PENDING"         // BUYER HAS PAID, BUT SELLER HAS NOT SEE THE ORDER
	PROCESS         OrderStatus = "PROCESS"         // BUYER HAS PAID AND SELLER DECIDE TO PROCESS THE ORDER
	CANCELED        OrderStatus = "CANCELED"        // BUYER HAS PAID BUT SELLER CANCEL THE ORDER
	REFUNDED        OrderStatus = "REFUNDED"        // BUYER HAS PAID BUT DECIDE TO REFUND THE ORDER
	COMPLETED       OrderStatus = "COMPLETED"       // BUYER HAS PAID AND CONFIRM THAT THE ORDER IS COMPLETE
)

var OrderStatuses = map[string]OrderStatus{
	string(PENDING):   PENDING,
	string(PROCESS):   PROCESS,
	string(CANCELED):  CANCELED,
	string(REFUNDED):  REFUNDED,
	string(COMPLETED): COMPLETED,
}

func (p *OrderStatus) Scan(value interface{}) error {
	*p = OrderStatus(value.([]byte))
	return nil
}

func (p OrderStatus) Value() (driver.Value, error) {
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
	Status             OrderStatus        `json:"status" gorm:"type:order_status;not null;"`
	Price              uint64             `json:"price" gorm:"not null;"`
	ServiceFee         uint               `json:"serviceFee" gorm:"not null; default:0;"`
	PaymentIntentID    string             `json:"paymentIntentId" gorm:"not null;"`
	DeliveredHistories []DeliveredHistory `json:"deliveredHistories" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	OrderEvents        OrderEvents        `json:"orderEvents" gorm:"type:jsonb; serializer:json;"`
	InvoiceID          string             `json:"invoiceId,omitempty"`
	StartDate          time.Time          `json:"startDate" gorm:"not null;"`
	Deadline           time.Time          `json:"deadline" gorm:"not null;"`
}

type CreateOrderDTO struct {
	ClientSecret    string `json:"clientSecret" validate:"required"`
	SellerID        string `json:"sellerId" validate:"required"`
	BuyerID         string `json:"buyerId"`
	GigTitle        string `json:"gigTitle" validate:"required"`
	GigDescription  string `json:"gigDescription" validate:"required"`
	Price           uint64 `json:"price" validate:"required"`
	ServiceFee      uint   `json:"serviceFee"`
	PaymentIntentID string `json:"paymentIntentId"`
	Deadline        int    `json:"deadline" validate:"required"`
}

type CreatePaymentIntentDTO struct {
	Amount int64 `json:"amount"`
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
