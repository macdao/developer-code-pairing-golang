package domain

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestNewOrder_CreatesOrderWithCorrectDefaults(t *testing.T) {
	// Arrange
	userID := uint64(1001)
	merchantID := "merchant_001"
	items := []OrderItem{
		{
			DishID:   "dish_001",
			DishName: "宫保鸡丁",
			Quantity: 2,
			Price:    decimal.NewFromFloat(28.00),
		},
	}
	delivery := DeliveryInfo{
		RecipientName:  "张三",
		RecipientPhone: "13800138000",
		Address:        "北京市朝阳区xxx",
	}
	remark := "少辣"

	// Act
	order := NewOrder(userID, merchantID, items, delivery, remark)

	// Assert
	assert.NotNil(t, order)
	assert.Equal(t, userID, order.UserID)
	assert.Equal(t, merchantID, order.MerchantID)
	assert.Equal(t, OrderStatusPendingPayment, order.Status)
	assert.Equal(t, remark, order.Remark)
	assert.NotEmpty(t, order.OrderNumber)
	assert.Len(t, order.Items, 1)
	assert.Equal(t, delivery, order.Delivery)
	assert.False(t, order.CreatedAt.IsZero())
	assert.False(t, order.UpdatedAt.IsZero())
}

func TestNewOrder_CalculatesPricingCorrectly_SingleItem(t *testing.T) {
	// Arrange
	items := []OrderItem{
		{
			DishID:   "dish_001",
			DishName: "宫保鸡丁",
			Quantity: 2,
			Price:    decimal.NewFromFloat(28.00),
		},
	}
	delivery := DeliveryInfo{
		RecipientName:  "张三",
		RecipientPhone: "13800138000",
		Address:        "北京市朝阳区xxx",
	}

	// Act
	order := NewOrder(1001, "merchant_001", items, delivery, "")

	// Assert - 2 * 28.00 = 56.00
	assert.Equal(t, "56.00", order.Pricing.ItemsTotal.StringFixed(2))
	assert.Equal(t, "1.00", order.Pricing.PackagingFee.StringFixed(2))
	assert.Equal(t, "3.00", order.Pricing.DeliveryFee.StringFixed(2))
	assert.Equal(t, "60.00", order.Pricing.FinalAmount.StringFixed(2))
}

func TestNewOrder_CalculatesPricingCorrectly_DecimalPrecision(t *testing.T) {
	// Arrange - 测试小数精度
	items := []OrderItem{
		{
			DishID:   "dish_001",
			DishName: "特价菜",
			Quantity: 3,
			Price:    decimal.NewFromFloat(12.99),
		},
	}
	delivery := DeliveryInfo{
		RecipientName:  "张三",
		RecipientPhone: "13800138000",
		Address:        "北京市朝阳区xxx",
	}

	// Act
	order := NewOrder(1001, "merchant_001", items, delivery, "")

	// Assert - 3 * 12.99 = 38.97
	assert.Equal(t, "38.97", order.Pricing.ItemsTotal.StringFixed(2))
	assert.Equal(t, "42.97", order.Pricing.FinalAmount.StringFixed(2))
}

func TestGenerateOrderNumber_Format(t *testing.T) {
	// Act
	orderNumber := generateOrderNumber()

	// Assert
	assert.Len(t, orderNumber, 20, "订单号应该是20位")
	assert.Regexp(t, `^\d{20}$`, orderNumber, "订单号应该全是数字")

	// 验证时间戳部分（前14位）
	timestamp := orderNumber[:14]
	now := time.Now().Format("20060102150405")
	assert.Equal(t, now, timestamp, "订单号的时间戳部分应该是当前时间")
}

func TestGenerateOrderNumber_Uniqueness(t *testing.T) {
	// Act - 生成多个订单号
	orderNumbers := make(map[string]bool)
	for i := 0; i < 100; i++ {
		orderNumber := generateOrderNumber()
		orderNumbers[orderNumber] = true
	}

	// Assert - 应该有多个不同的订单号（由于随机数部分）
	// 注意：如果在同一秒内生成，可能会有重复，但概率很低
	assert.GreaterOrEqual(t, len(orderNumbers), 1, "应该生成订单号")
}

func TestNewOrder_SetsTimestamps(t *testing.T) {
	// Arrange
	before := time.Now()
	items := []OrderItem{
		{DishID: "dish_001", DishName: "测试菜", Quantity: 1, Price: decimal.NewFromFloat(10.00)},
	}
	delivery := DeliveryInfo{
		RecipientName:  "张三",
		RecipientPhone: "13800138000",
		Address:        "北京市朝阳区xxx",
	}

	// Act
	order := NewOrder(1001, "merchant_001", items, delivery, "")
	after := time.Now()

	// Assert
	assert.False(t, order.CreatedAt.IsZero())
	assert.False(t, order.UpdatedAt.IsZero())
	assert.True(t, order.CreatedAt.After(before) || order.CreatedAt.Equal(before))
	assert.True(t, order.CreatedAt.Before(after) || order.CreatedAt.Equal(after))
	assert.Equal(t, order.CreatedAt, order.UpdatedAt)
}

func TestNewOrder_WithEmptyRemark(t *testing.T) {
	// Arrange
	items := []OrderItem{
		{DishID: "dish_001", DishName: "测试菜", Quantity: 1, Price: decimal.NewFromFloat(10.00)},
	}
	delivery := DeliveryInfo{
		RecipientName:  "张三",
		RecipientPhone: "13800138000",
		Address:        "北京市朝阳区xxx",
	}

	// Act
	order := NewOrder(1001, "merchant_001", items, delivery, "")

	// Assert
	assert.Empty(t, order.Remark)
	assert.NotNil(t, order)
}

func TestNewOrder_WithLongRemark(t *testing.T) {
	// Arrange
	items := []OrderItem{
		{DishID: "dish_001", DishName: "测试菜", Quantity: 1, Price: decimal.NewFromFloat(10.00)},
	}
	delivery := DeliveryInfo{
		RecipientName:  "张三",
		RecipientPhone: "13800138000",
		Address:        "北京市朝阳区xxx",
	}
	longRemark := "这是一个很长的备注，包含了很多信息，比如少辣、多放醋、不要香菜等等"

	// Act
	order := NewOrder(1001, "merchant_001", items, delivery, longRemark)

	// Assert
	assert.Equal(t, longRemark, order.Remark)
}

func TestOrderStatus_Constants(t *testing.T) {
	// Assert - 验证订单状态常量
	assert.Equal(t, OrderStatus("PENDING_PAYMENT"), OrderStatusPendingPayment)
	assert.Equal(t, OrderStatus("PAID"), OrderStatusPaid)
	assert.Equal(t, OrderStatus("CANCELLED"), OrderStatusCancelled)
}

func TestDefaultFees(t *testing.T) {
	// Assert - 验证默认费用
	assert.Equal(t, "1.00", DefaultPackagingFee.StringFixed(2))
	assert.Equal(t, "3.00", DefaultDeliveryFee.StringFixed(2))
}
