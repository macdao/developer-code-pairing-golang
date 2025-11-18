package application

import (
	"context"
	"fmt"
	"testing"

	"order-service/internal/domain"

	"github.com/stretchr/testify/assert"
)

// MockOrderRepository 模拟 Repository
type MockOrderRepository struct {
	orders map[string]*domain.Order
}

func NewMockOrderRepository() OrderRepository {
	return &MockOrderRepository{
		orders: make(map[string]*domain.Order),
	}
}

func (m *MockOrderRepository) Create(ctx context.Context, order *domain.Order) error {
	if _, exists := m.orders[order.OrderNumber]; exists {
		return fmt.Errorf("order already exists")
	}
	m.orders[order.OrderNumber] = order
	return nil
}

func (m *MockOrderRepository) FindByOrderNumber(ctx context.Context, orderNumber string) (*domain.Order, error) {
	order, exists := m.orders[orderNumber]
	if !exists {
		return nil, NewNotFoundError("order not found")
	}
	return order, nil
}

func TestOrderService_CreateOrder_Success(t *testing.T) {
	// Arrange
	repo := NewMockOrderRepository().(*MockOrderRepository)
	service := NewOrderService(repo)
	ctx := context.Background()

	req := &CreateOrderRequest{
		MerchantID: "merchant_001",
		Items: []OrderItemRequest{
			{
				DishID:   "dish_001",
				DishName: "宫保鸡丁",
				Quantity: 2,
				Price:    28.00,
			},
		},
		DeliveryInfo: DeliveryInfoRequest{
			RecipientName:  "张三",
			RecipientPhone: "13800138000",
			Address:        "北京市朝阳区xxx街道xxx号",
		},
		Remark: "少辣，多放醋",
	}

	// Act
	orderData, err := service.CreateOrder(ctx, 1001, req)

	// Assert - 验证返回值
	assert.NoError(t, err)
	assert.NotNil(t, orderData)
	assert.NotEmpty(t, orderData.OrderNumber)
	assert.Equal(t, "PENDING_PAYMENT", orderData.Status)
	assert.Equal(t, "56.00", orderData.Pricing.ItemsTotal)
	assert.Equal(t, "1.00", orderData.Pricing.PackagingFee)
	assert.Equal(t, "3.00", orderData.Pricing.DeliveryFee)
	assert.Equal(t, "60.00", orderData.Pricing.FinalAmount)
	assert.NotEmpty(t, orderData.CreatedAt)

	// Assert - 验证状态：订单确实被保存了
	savedOrder, err := repo.FindByOrderNumber(ctx, orderData.OrderNumber)
	assert.NoError(t, err)
	assert.NotNil(t, savedOrder)
	assert.Equal(t, orderData.OrderNumber, savedOrder.OrderNumber)
	assert.Equal(t, uint64(1001), savedOrder.UserID)
	assert.Equal(t, "merchant_001", savedOrder.MerchantID)
}

func TestOrderService_CreateOrder_ValidationError(t *testing.T) {
	// Arrange
	repo := NewMockOrderRepository()
	service := NewOrderService(repo)
	ctx := context.Background()

	// 无效的请求（空商家ID）
	req := &CreateOrderRequest{
		MerchantID: "", // 验证失败
		Items: []OrderItemRequest{
			{DishID: "dish_001", DishName: "宫保鸡丁", Quantity: 1, Price: 28.00},
		},
		DeliveryInfo: DeliveryInfoRequest{
			RecipientName:  "张三",
			RecipientPhone: "13800138000",
			Address:        "北京市朝阳区xxx",
		},
	}

	// Act
	orderData, err := service.CreateOrder(ctx, 1001, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, orderData)
	assert.IsType(t, &ValidationError{}, err)
}
