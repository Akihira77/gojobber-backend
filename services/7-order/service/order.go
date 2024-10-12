package service

import (
	"context"
	"fmt"
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
	FindOrdersByBuyerID(ctx context.Context, id string) ([]types.Order, error)
	FindOrdersBySellerID(ctx context.Context, id string) ([]types.Order, error)
	CreateOrder(ctx context.Context, data *types.CreateOrderDTO) (*types.Order, error)
	RefundingOrder(ctx context.Context) error
	ChangeOrderStatus(ctx context.Context, o *types.Order, newStatus types.OrderStatus) error
	ExtendingDeadline(ctx context.Context, o *types.Order, numberOfDays int) error
	DeliveringOrder(ctx context.Context, o *types.Order, dh *types.DeliveredHistory) error
	FindMyOrderNotifications(ctx context.Context, sellerId string) error
}

func NewOrderService(db *gorm.DB) OrderServiceImpl {
	return &OrderService{
		db: db,
	}
}

// TODO: REFUNDING ORDER
func (os *OrderService) RefundingOrder(ctx context.Context) error {
	panic("unimplemented")
}

func (os *OrderService) ChangeOrderStatus(ctx context.Context, o *types.Order, newStatus types.OrderStatus) error {
	o.Status = newStatus
	o.OrderEvents = append(o.OrderEvents, types.OrderEvent{
		Event:     fmt.Sprintf("Order Status Changed To [%s]", newStatus),
		CreatedAt: time.Now(),
	})

	result := os.db.
		Debug().
		WithContext(ctx).
		Save(o)

	return result.Error
}

func (os *OrderService) CreateOrder(ctx context.Context, data *types.CreateOrderDTO) (*types.Order, error) {
	var orderEvents types.OrderEvents
	orderEvents = append(orderEvents, types.OrderEvent{
		Event:     fmt.Sprintf("Buyer Has Purchased Your Gig With PaymentID: %s", data.PaymentIntentID),
		CreatedAt: time.Now(),
	})

	startDate := time.Now()
	newOrder := types.Order{
		SellerID:           data.SellerID,
		BuyerID:            data.BuyerID,
		GigTitle:           data.GigTitle,
		GigDescription:     data.GigDescription,
		Price:              data.Price,
		Status:             types.OrderStatuses[string(types.PENDING)],
		ServiceFee:         uint(math.Ceil((25 / 100) * float64(data.Price))),
		PaymentIntentID:    data.PaymentIntentID,
		StartDate:          startDate,
		Deadline:           startDate.AddDate(0, 0, data.Deadline),
		InvoiceID:          fmt.Sprintf("JI%s", util.RandomStr(30)),
		OrderEvents:        orderEvents,
		DeliveredHistories: []types.DeliveredHistory{},
	}

	result := os.db.
		Debug().
		WithContext(ctx).
		Model(&types.Order{}).
		Create(&newOrder)

	return &newOrder, result.Error
}

func (os *OrderService) DeliveringOrder(ctx context.Context, o *types.Order, dh *types.DeliveredHistory) error {
	tx := os.db.
		Debug().
		WithContext(ctx)

	result := tx.
		Model(&types.DeliveredHistory{}).
		Create(dh)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	o.OrderEvents = append(o.OrderEvents, types.OrderEvent{
		Event:     fmt.Sprintf("Seller Delivering The Order"),
		CreatedAt: time.Now(),
	})
	result = tx.Save(o)
	if result.Error != nil {
		return result.Error
	}

	result = tx.Commit()

	return result.Error
}

func (os *OrderService) ExtendingDeadline(ctx context.Context, o *types.Order, numberOfDays int) error {
	o.Deadline = o.Deadline.Add(time.Duration(numberOfDays))
	o.OrderEvents = append(o.OrderEvents, types.OrderEvent{
		Event:     "Buyer Approved Seller Extending Date Request",
		CreatedAt: time.Now(),
	})

	result := os.db.
		Debug().
		WithContext(ctx).
		Save(o)

	return result.Error
}

func (os *OrderService) FindMyOrderNotifications(ctx context.Context, sellerId string) error {
	panic("unimplemented")
}

func (os *OrderService) FindOrderByID(ctx context.Context, id string) (*types.Order, error) {
	var o types.Order
	result := os.db.
		Debug().
		WithContext(ctx).
		Model(&types.Order{}).
		Where("id = ?", id).
		First(&o)

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
