package persistence

import (
	"context"
	"fmt"

	"order-service/internal/application"
	"order-service/internal/domain"
)

// InMemoryOrderRepository 内存订单仓储实现
type InMemoryOrderRepository struct {
	orders map[string]*domain.Order // 按 OrderNumber 索引
}

// NewInMemoryOrderRepository 创建内存仓储实例
func NewInMemoryOrderRepository() application.OrderRepository {
	return &InMemoryOrderRepository{
		orders: make(map[string]*domain.Order),
	}
}

// Create 创建订单
func (r *InMemoryOrderRepository) Create(ctx context.Context, order *domain.Order) error {
	// 检查订单号唯一性
	if _, exists := r.orders[order.OrderNumber]; exists {
		return fmt.Errorf("order number %s already exists", order.OrderNumber)
	}

	r.orders[order.OrderNumber] = order
	return nil
}

// FindByOrderNumber 根据订单号查询订单
func (r *InMemoryOrderRepository) FindByOrderNumber(ctx context.Context, orderNumber string) (*domain.Order, error) {
	order, exists := r.orders[orderNumber]
	if !exists {
		return nil, application.NewNotFoundError(fmt.Sprintf("order %s not found", orderNumber))
	}

	return order, nil
}
