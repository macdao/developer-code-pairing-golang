package domain

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/shopspring/decimal"
)

// OrderStatus 订单状态
type OrderStatus string

const (
	OrderStatusPendingPayment OrderStatus = "PENDING_PAYMENT"
	OrderStatusPaid           OrderStatus = "PAID"
	OrderStatusCancelled      OrderStatus = "CANCELLED"
)

// 业务常量
var (
	DefaultPackagingFee = decimal.NewFromFloat(1.00)
	DefaultDeliveryFee  = decimal.NewFromFloat(3.00)
)

// Order 订单聚合根
type Order struct {
	OrderNumber string
	UserID      uint64
	MerchantID  string
	Status      OrderStatus
	Pricing     Pricing
	Delivery    DeliveryInfo
	Remark      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Items       []OrderItem
}

// Pricing 价格信息值对象
type Pricing struct {
	ItemsTotal   decimal.Decimal
	PackagingFee decimal.Decimal
	DeliveryFee  decimal.Decimal
	FinalAmount  decimal.Decimal
}

// DeliveryInfo 配送信息值对象
type DeliveryInfo struct {
	RecipientName  string
	RecipientPhone string
	Address        string
}

// OrderItem 订单项实体
type OrderItem struct {
	DishID   string
	DishName string
	Quantity int
	Price    decimal.Decimal
}

// NewOrder 创建新订单（工厂方法）
func NewOrder(userID uint64, merchantID string, items []OrderItem, delivery DeliveryInfo, remark string) *Order {
	now := time.Now()
	
	order := &Order{
		OrderNumber: generateOrderNumber(),
		UserID:      userID,
		MerchantID:  merchantID,
		Status:      OrderStatusPendingPayment,
		Items:       items,
		Delivery:    delivery,
		Remark:      remark,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	
	order.calculatePricing()
	return order
}

// calculatePricing 计算订单价格（私有方法，创建时自动调用）
func (o *Order) calculatePricing() {
	// 计算餐品总价
	o.Pricing.ItemsTotal = decimal.Zero
	for _, item := range o.Items {
		itemTotal := item.Price.Mul(decimal.NewFromInt(int64(item.Quantity)))
		o.Pricing.ItemsTotal = itemTotal
	}
	
	// 设置固定费用
	o.Pricing.PackagingFee = DefaultPackagingFee
	o.Pricing.DeliveryFee = DefaultDeliveryFee
	
	// 计算最终金额
	o.Pricing.FinalAmount = o.Pricing.ItemsTotal.
		Add(o.Pricing.PackagingFee).
		Add(o.Pricing.DeliveryFee)
}

// generateOrderNumber 生成订单号（格式：yyyyMMddHHmmss + 6位随机数）
func generateOrderNumber() string {
	timestamp := time.Now().Format("20060102150405")
	random := rand.Intn(1000000)
	return fmt.Sprintf("%s%06d", timestamp, random)
}
