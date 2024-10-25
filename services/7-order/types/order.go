package types

import (
	"database/sql/driver"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeadlineExtensionStatus string

const (
	ACCEPTED DeadlineExtensionStatus = "ACCEPTED"
	REJECTED DeadlineExtensionStatus = "REJECTED"
)

type OrderStatus string

const (
	AWAITING_PAYMENT OrderStatus = "AWAITING PAYMENT" //NOTE: BUYER HAS NOT PAID
	PENDING          OrderStatus = "PENDING"          // BUYER HAS PAID, BUT SELLER HAS NOT SEE THE ORDER
	PROCESS          OrderStatus = "PROCESS"          // BUYER HAS PAID AND SELLER DECIDE TO PROCESS THE ORDER
	CANCELED         OrderStatus = "CANCELED"         // BUYER HAS PAID BUT SELLER CANCEL THE ORDER
	REFUNDED         OrderStatus = "REFUNDED"         // BUYER HAS PAID BUT DECIDE TO REFUND THE ORDER
	COMPLETED        OrderStatus = "COMPLETED"        // BUYER HAS PAID AND CONFIRM THAT THE ORDER IS COMPLETE
)

func (p *OrderStatus) Scan(value interface{}) error {
	*p = OrderStatus(value.(string))
	// log.Println("scan p", p)
	return nil
}

func (p OrderStatus) Value() (driver.Value, error) {
	// log.Println("value p", p)
	return string(p), nil
}

type OrderEvent struct {
	ID        uint64    `json:"id" gorm:"primaryKey;autoIncrement;"`
	Event     string    `json:"event" gorm:"not null;"`
	CreatedAt time.Time `json:"createdAt" gorm:"not null;"`
	OrderID   string    `json:"orderId"`
}

type DeliveredHistory struct {
	ID            uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4();"`
	OrderID       string    `json:"orderId" validate:"required"`
	ProgressNote  string    `json:"progressNote" validate:"required" gorm:"not null;"`
	ResultURL     string    `json:"resultUrl" validate:"url" gorm:""`
	BuyerNote     string    `json:"buyerNote" gorm:""`
	DeliveredDate time.Time `json:"deliveredDate" gorm:"not null;"`
}

type BuyerResponseOrderDelivered struct {
	ID            uuid.UUID `json:"id" validate:"required,uuid"`
	OrderID       string    `json:"orderId" validate:"required"`
	ProgressNote  string    `json:"progressNote" validate:"required"`
	ResultURL     string    `json:"resultUrl"`
	BuyerNote     string    `json:"buyerNote" validate:"required"`
	DeliveredDate time.Time `json:"deliveredDate" validate:"required"`
}

type Order struct {
	ID                 string             `json:"id" gorm:"primaryKey; not null"`
	SellerID           string             `json:"sellerId" gorm:"not null;"`
	BuyerID            string             `json:"buyerId" gorm:"not null;"`
	GigTitle           string             `json:"gigTitle" gorm:"not null;"`
	GigDescription     string             `json:"gigDescription" gorm:"not null;"`
	Status             OrderStatus        `json:"status" gorm:"type:order_status;not null;"`
	Price              uint64             `json:"price" gorm:"not null;"`
	ServiceFee         uint               `json:"serviceFee" gorm:"not null; default:0;"`
	PaymentIntentID    string             `json:"paymentIntentId" gorm:"unique; not null;"`
	StripeClientSecret string             `json:"stripeClientSecret" gorm:"unique; not null;"`
	DeliveredHistories []DeliveredHistory `json:"deliveredHistories,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	OrderEvents        []OrderEvent       `json:"orderEvents,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	InvoiceID          string             `json:"invoiceId,omitempty"`
	StartDate          time.Time          `json:"startDate" gorm:"not null;"`
	Deadline           time.Time          `json:"deadline" gorm:"not null;"`
	// Unread             bool               `json:"unread" gorm:"default:true;not null;"`
}

type OrderNotificationDTO struct {
	ID        string      `json:"id"`
	GigTitle  string      `json:"gigTitle"`
	Status    OrderStatus `json:"status"`
	Price     uint64      `json:"price"`
	StartDate time.Time   `json:"startDate"`
	Deadline  time.Time   `json:"deadline"`
}

type CreateOrderDTO struct {
	SellerID           string `json:"sellerId" validate:"required"`
	BuyerID            string `json:"buyerId"`
	GigTitle           string `json:"gigTitle" validate:"required"`
	GigDescription     string `json:"gigDescription"`
	Price              uint64 `json:"price" validate:"required"`
	ServiceFee         uint   `json:"serviceFee"`
	PaymentIntentID    string `json:"paymentIntentId"`
	StripeClientSecret string `json:"stripeClientSecret"`
	Deadline           int    `json:"deadline" validate:"required"`
	MessageID          string `json:"messageId,omitempty"`
}

type DeadlineExtensionRequest struct {
	NumberOfDays  int                     `json:"numberOfDays" validate:"required,gte=1,lte=365"`
	Reason        string                  `json:"reason" validate:"required"`
	BuyerResponse DeadlineExtensionStatus `json:"buyerResponse" validate:"oneof=ACCEPTED REJECTED"`
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
