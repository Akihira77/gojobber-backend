package service

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/Akihira77/gojobber/services/7-order/types"
	"gorm.io/gorm"
)

type OrderService struct {
	db *gorm.DB
}

type OrderServiceImpl interface {
	FindOrderByID(ctx context.Context, id string) (*types.Order, error)
	FindOrdersByBuyerID(ctx context.Context, id string) ([]types.Order, error)
	FindOrdersBySellerID(ctx context.Context, id string) ([]types.Order, error)
	CreateOrder(ctx context.Context) error
	CreatePaymentIntent(ctx context.Context) error
	ChangeOrderStatus(ctx context.Context, o *types.Order, newStatus string) error
	ExtendingDeadline(ctx context.Context, o *types.Order, addDays int) error
	DeliveringOrder(ctx context.Context, o *types.Order, dh *types.DeliveredHistory) error
	FindMyOrderNotifications(ctx context.Context, sellerId string) error
}

func NewOrderService(db *gorm.DB) OrderServiceImpl {
	return &OrderService{
		db: db,
	}
}

func (os *OrderService) ChangeOrderStatus(ctx context.Context, o *types.Order, newStatus string) error {
	if !slices.Contains(types.OrderStatuses, newStatus) {
		return fmt.Errorf("Unknown order status")
	}

	o.Status = types.OrderStatus(newStatus)
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

// TODO: ADD STRIPE FIRST
func (os *OrderService) CreateOrder(ctx context.Context) error {
	panic("unimplemented")
}

// TODO: ADD STRIPE FIRST
func (os *OrderService) CreatePaymentIntent(ctx context.Context) error {
	panic("unimplemented")
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

func (os *OrderService) ExtendingDeadline(ctx context.Context, o *types.Order, addDays int) error {
	o.Deadline = o.Deadline.Add(time.Duration(addDays))
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
