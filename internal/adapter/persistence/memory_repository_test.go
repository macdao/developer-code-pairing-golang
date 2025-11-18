package persistence

import (
	"context"
	"testing"

	"order-service/internal/domain"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryOrderRepository_Create(t *testing.T) {
	repo := NewInMemoryOrderRepository()
	ctx := context.Background()

	order := createTestOrder("20241117120000123456")

	err := repo.Create(ctx, order)
	assert.NoError(t, err)

	// 验证可以查询到订单
	found, err := repo.FindByOrderNumber(ctx, order.OrderNumber)
	assert.NoError(t, err)
	assert.Equal(t, order.OrderNumber, found.OrderNumber)
	assert.Equal(t, order.UserID, found.UserID)
}

func TestInMemoryOrderRepository_Create_DuplicateOrderNumber(t *testing.T) {
	repo := NewInMemoryOrderRepository()
	ctx := context.Background()

	order1 := createTestOrder("20241117120000123456")
	order2 := createTestOrder("20241117120000123456")

	err := repo.Create(ctx, order1)
	assert.NoError(t, err)

	// 尝试创建相同订单号的订单
	err = repo.Create(ctx, order2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

func TestInMemoryOrderRepository_FindByOrderNumber_NotFound(t *testing.T) {
	repo := NewInMemoryOrderRepository()
	ctx := context.Background()

	_, err := repo.FindByOrderNumber(ctx, "nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}



// createTestOrder 创建测试订单
func createTestOrder(orderNumber string) *domain.Order {
	items := []domain.OrderItem{
		{
			DishID:   "dish1",
			DishName: "宫保鸡丁",
			Quantity: 2,
			Price:    decimal.NewFromFloat(28.00),
		},
	}

	delivery := domain.DeliveryInfo{
		RecipientName:  "张三",
		RecipientPhone: "13800138000",
		Address:        "北京市朝阳区xxx",
	}

	order := domain.NewOrder(1001, "merchant1", items, delivery, "少辣")
	// 覆盖订单号以便测试
	order.OrderNumber = orderNumber
	return order
}
