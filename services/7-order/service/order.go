package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/Akihira77/gojobber/services/7-order/types"
	"github.com/Akihira77/gojobber/services/7-order/util"
	"gorm.io/gorm"
)

type OrderService struct {
	db *gorm.DB
}

type OrderServiceImpl interface {
	FindOrderByID(ctx context.Context, id string) (*types.Order, error)
	FindOrderByIDAndSellerID(ctx context.Context, id string, sellerId string) (*types.Order, error)
	FindOrderByPaymentIntentID(ctx context.Context, id string) (*types.Order, error)
	FindOrdersByBuyerID(ctx context.Context, id string) ([]types.Order, error)
	FindOrdersBySellerID(ctx context.Context, id string) ([]types.Order, error)
	SaveOrder(ctx context.Context, data *types.CreateOrderDTO) (*types.Order, error)
	ChangeOrderStatus(ctx context.Context, o types.Order, newStatus types.OrderStatus, msg string) (*types.Order, error)
	RequestDeadlineExtension(ctx context.Context, o types.Order, data *types.DeadlineExtensionRequest) error
	DeadlineExtensionResponse(ctx context.Context, o types.Order, status types.DeadlineExtensionStatus, data *types.DeadlineExtensionRequest) (string, error)
	DeliveringOrder(ctx context.Context, o types.Order, dh types.DeliveredHistory) (*types.Order, error)
	OrderDeliveredResponse(ctx context.Context, o types.Order, r *types.BuyerResponseOrderDelivered) (*types.Order, error)
	FindMyOrderNotifications(ctx context.Context, userID string) ([]types.OrderNotificationDTO, error)
	MarkReadsMyOrderNotifications(ctx context.Context, userID string) error
}

func NewOrderService(db *gorm.DB) OrderServiceImpl {
	return &OrderService{
		db: db,
	}
}

func (os *OrderService) ChangeOrderStatus(ctx context.Context, o types.Order, newStatus types.OrderStatus, msg string) (*types.Order, error) {
	tx := os.db.
		Debug().
		WithContext(ctx).
		Begin()

	o.Status = newStatus
	result := tx.
		Model(&types.Order{}).
		Where("id = ?", o.ID).
		Update("status", newStatus)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}

	result = tx.
		Model(&types.OrderEvent{}).
		Create(&types.OrderEvent{
			OrderID:   o.ID,
			Event:     msg,
			CreatedAt: time.Now(),
		})
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}

	result = tx.Commit()

	return &o, result.Error
}

func (os *OrderService) SaveOrder(ctx context.Context, data *types.CreateOrderDTO) (*types.Order, error) {
	tx := os.db.
		Debug().
		WithContext(ctx).
		Begin()

	startDate := time.Now()
	newOrder := types.Order{
		ID:                 fmt.Sprintf("JO%s", util.RandomStr(30)),
		SellerID:           data.SellerID,
		BuyerID:            data.BuyerID,
		GigTitle:           data.GigTitle,
		GigDescription:     data.GigDescription,
		Price:              data.Price,
		Status:             types.AWAITING_PAYMENT,
		ServiceFee:         uint(math.Ceil((25 / 1000) * float64(data.Price))),
		PaymentIntentID:    data.PaymentIntentID,
		StripeClientSecret: data.StripeClientSecret,
		StartDate:          startDate,
		Deadline:           startDate.AddDate(0, 0, data.Deadline),
		InvoiceID:          fmt.Sprintf("JI%s", util.RandomStr(30)),
	}

	result := tx.
		Model(&types.Order{}).
		Create(&newOrder)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}

	result = tx.
		Model(&types.OrderEvent{}).
		Create(&types.OrderEvent{
			OrderID:   newOrder.ID,
			Event:     fmt.Sprintf("Buyer Has Purchased Your Gig With PaymentID: %s", data.PaymentIntentID),
			CreatedAt: time.Now(),
		})
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}

	result = tx.Commit()

	return &newOrder, result.Error
}

func (os *OrderService) OrderDeliveredResponse(ctx context.Context, o types.Order, r *types.BuyerResponseOrderDelivered) (*types.Order, error) {
	tx := os.db.
		Debug().
		WithContext(ctx).
		Begin(&sql.TxOptions{
			Isolation: sql.LevelSerializable,
		})

	result := tx.
		Model(&types.DeliveredHistory{}).
		Where("id = ?", r.ID).
		Update("buyer_note", r.BuyerNote)
	if result.Error != nil {
		tx.Rollback()
		return &o, result.Error
	}

	result = tx.
		Model(&types.OrderEvent{}).
		Create(&types.OrderEvent{
			OrderID:   o.ID,
			Event:     fmt.Sprintf("Buyer Responded Your Order Delivered Progress"),
			CreatedAt: time.Now(),
		})
	if result.Error != nil {
		tx.Rollback()
		return &o, result.Error
	}

	result = tx.
		Model(&types.Order{}).
		Preload("DeliveredHistories").
		Preload("OrderEvents").
		Where("id = ?", o.ID).
		First(&o)
	if result.Error != nil {
		tx.Rollback()
		return &o, result.Error
	}

	result = tx.Commit()

	return &o, result.Error

}

func (os *OrderService) DeliveringOrder(ctx context.Context, o types.Order, dh types.DeliveredHistory) (*types.Order, error) {
	tx := os.db.
		Debug().
		WithContext(ctx).
		Begin(&sql.TxOptions{
			Isolation: sql.LevelSerializable,
		})

	result := tx.
		Model(&types.DeliveredHistory{}).
		Create(&types.DeliveredHistory{
			OrderID:       o.ID,
			ResultURL:     dh.ResultURL,
			BuyerNote:     dh.BuyerNote,
			ProgressNote:  dh.ProgressNote,
			DeliveredDate: time.Now(),
		})
	if result.Error != nil {
		tx.Rollback()
		return &o, result.Error
	}

	result = tx.
		Model(&types.OrderEvent{}).
		Create(&types.OrderEvent{
			OrderID:   o.ID,
			Event:     fmt.Sprintf("Seller Delivering The Order"),
			CreatedAt: time.Now(),
		})
	if result.Error != nil {
		tx.Rollback()
		return &o, result.Error
	}

	result = tx.
		Model(&types.Order{}).
		Preload("DeliveredHistories").
		Preload("OrderEvents").
		Where("id = ?", o.ID).
		First(&o)
	if result.Error != nil {
		tx.Rollback()
		return &o, result.Error
	}

	result = tx.Commit()

	return &o, result.Error
}

func (os *OrderService) RequestDeadlineExtension(ctx context.Context, o types.Order, data *types.DeadlineExtensionRequest) error {
	tx := os.db.
		Debug().
		WithContext(ctx).
		Begin()

	deadline := o.Deadline.AddDate(0, 0, data.NumberOfDays)
	result := tx.
		Model(&types.OrderEvent{}).
		Create(&types.OrderEvent{
			OrderID:   o.ID,
			Event:     fmt.Sprintf("Seller Requesting To Extend The Deadline To Become [%v] With Reason:\n%s", deadline, data.Reason),
			CreatedAt: time.Now(),
		})
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	result = tx.Commit()

	return result.Error
}

func (os *OrderService) DeadlineExtensionResponse(ctx context.Context, o types.Order, status types.DeadlineExtensionStatus, data *types.DeadlineExtensionRequest) (string, error) {
	tx := os.db.
		Debug().
		WithContext(ctx).
		Begin()

	switch status {
	case types.ACCEPTED:
		deadline := o.Deadline.AddDate(0, 0, data.NumberOfDays)
		msg := fmt.Sprintf("Buyer Accepted Seller Order Deadline Extension From [%v] To [%v]", o.Deadline, deadline)
		result := tx.
			Model(&types.Order{}).
			Where("id = ?", o.ID).
			Update("deadline", deadline)
		if result.Error != nil {
			tx.Rollback()
			return "", result.Error
		}

		result = tx.
			Model(&types.OrderEvent{}).
			Create(&types.OrderEvent{
				OrderID:   o.ID,
				Event:     msg,
				CreatedAt: time.Now(),
			})
		if result.Error != nil {
			tx.Rollback()
			return "", result.Error
		}

		result = tx.Commit()

		return msg, result.Error
	case types.REJECTED:
		msg := fmt.Sprintf("Buyer Rejected Your Order Deadline Extension With Reason:\n%s", data.Reason)
		result := tx.
			Model(&types.OrderEvent{}).
			Create(&types.OrderEvent{
				OrderID:   o.ID,
				Event:     msg,
				CreatedAt: time.Now(),
			})
		if result.Error != nil {
			tx.Rollback()
			return "", result.Error
		}

		result = tx.Commit()

		return msg, result.Error
	default:
		return "", fmt.Errorf("Unknown deadline extension response")
	}
}

func (os *OrderService) FindMyOrderNotifications(ctx context.Context, userID string) ([]types.OrderNotificationDTO, error) {
	var orders []types.OrderNotificationDTO
	result := os.db.
		Debug().
		WithContext(ctx).
		Model(&types.Order{}).
		Where("buyer_id = ?", userID).
		Find(&orders)

	return orders, result.Error
}

func (os *OrderService) MarkReadsMyOrderNotifications(ctx context.Context, userID string) error {
	result := os.db.
		Debug().
		WithContext(ctx).
		Model(&types.Order{}).
		Where("buyer_id = ?", userID).
		Update("unread = ?", false)

	return result.Error
}

func (os *OrderService) FindOrderByIDAndSellerID(ctx context.Context, id string, sellerId string) (*types.Order, error) {
	var o types.Order
	result := os.db.
		Debug().
		WithContext(ctx).
		Model(&types.Order{}).
		Preload("DeliveredHistories").
		Preload("OrderEvents").
		Where("id = ? AND seller_id = ?", id, sellerId).
		First(&o)

	return &o, result.Error
}

func (os *OrderService) FindOrderByID(ctx context.Context, id string) (*types.Order, error) {
	var o types.Order
	result := os.db.
		Debug().
		WithContext(ctx).
		Model(&types.Order{}).
		Preload("DeliveredHistories").
		Preload("OrderEvents").
		Where("id = ?", id).
		First(&o)

	return &o, result.Error
}

func (os *OrderService) FindOrderByPaymentIntentID(ctx context.Context, id string) (*types.Order, error) {
	var o types.Order
	result := os.db.
		Debug().
		WithContext(ctx).
		Model(&types.Order{}).
		Preload("DeliveredHistories").
		Preload("OrderEvents").
		Where("payment_intent_id = ?", id).
		First(&o)

	log.Println(o)
	return &o, result.Error
}

func (os *OrderService) FindOrdersByBuyerID(ctx context.Context, id string) ([]types.Order, error) {
	var o []types.Order
	result := os.db.
		Debug().
		WithContext(ctx).
		Model(&types.Order{}).
		Where("buyer_id = ?", id).
		Find(&o)

	return o, result.Error
}

func (os *OrderService) FindOrdersBySellerID(ctx context.Context, id string) ([]types.Order, error) {
	var o []types.Order
	result := os.db.
		Debug().
		WithContext(ctx).
		Model(&types.Order{}).
		Where("seller_id = ?", id).
		Find(&o)

	return o, result.Error
}
