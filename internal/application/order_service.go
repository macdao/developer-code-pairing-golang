package application

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
	"order-service/internal/domain"
)

// orderService 应用服务实现
type orderService struct {
	repo OrderRepository
}

// NewOrderService 创建应用服务实例
func NewOrderService(repo OrderRepository) OrderService {
	return &orderService{repo: repo}
}

// CreateOrder 实现 OrderService 接口
func (s *orderService) CreateOrder(ctx context.Context, userID uint64, req *CreateOrderRequest) (*OrderData, error) {
	// 1. 验证请求数据（使用 validator）
	if err := Validator.Struct(req); err != nil {
		// 转换 validator 错误为应用层错误
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, e := range validationErrors {
				return nil, NewValidationError(e.Field(), e.Error())
			}
		}
		return nil, NewValidationError("", err.Error())
	}

	// 2. 转换 DTO 到领域对象
	items := s.convertToOrderItems(req.Items)
	delivery := domain.DeliveryInfo{
		RecipientName:  req.DeliveryInfo.RecipientName,
		RecipientPhone: req.DeliveryInfo.RecipientPhone,
		Address:        req.DeliveryInfo.Address,
	}

	// 3. 创建订单（领域对象负责初始化所有状态）
	order := domain.NewOrder(userID, req.MerchantID, items, delivery, req.Remark)

	// 4. 保存订单
	if err := s.repo.Create(ctx, order); err != nil {
		return nil, NewInternalError("failed to create order", err)
	}

	// 5. 返回结果
	return s.convertToDTO(order), nil
}

// convertToOrderItems 转换订单项
func (s *orderService) convertToOrderItems(items []OrderItemRequest) []domain.OrderItem {
	result := make([]domain.OrderItem, len(items))
	for i, item := range items {
		result[i] = domain.OrderItem{
			DishID:   item.DishID,
			DishName: item.DishName,
			Quantity: item.Quantity,
			Price:    decimal.NewFromFloat(item.Price),
		}
	}
	return result
}

// convertToDTO 转换领域对象到 DTO
func (s *orderService) convertToDTO(order *domain.Order) *OrderData {
	return &OrderData{
		OrderNumber: order.OrderNumber,
		Status:      string(order.Status),
		Pricing: PricingInfo{
			ItemsTotal:   order.Pricing.ItemsTotal.StringFixed(2),
			PackagingFee: order.Pricing.PackagingFee.StringFixed(2),
			DeliveryFee:  order.Pricing.DeliveryFee.StringFixed(2),
			FinalAmount:  order.Pricing.FinalAmount.StringFixed(2),
		},
		CreatedAt: order.CreatedAt.Format(time.RFC3339),
	}
}
