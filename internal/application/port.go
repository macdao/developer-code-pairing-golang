package application

import (
	"context"
	"regexp"

	"github.com/go-playground/validator/v10"
	"order-service/internal/domain"
)

// Validator 全局验证器实例
var Validator = validator.New()

// 手机号正则表达式：1开头，第二位是3-9，后面9位数字
var phoneRegex = regexp.MustCompile(`^1[3-9]\d{9}$`)

func init() {
	// 注册自定义手机号验证函数
	_ = Validator.RegisterValidation("phone", validatePhone)
}

// validatePhone 验证中国大陆手机号格式
func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	return phoneRegex.MatchString(phone)
}

// OrderService 定义应用服务接口（输入端口）
// Web 适配器通过此接口调用核心业务逻辑
type OrderService interface {
	CreateOrder(ctx context.Context, userID uint64, req *CreateOrderRequest) (*OrderData, error)
}

// OrderRepository 定义数据持久化接口（输出端口）
// 应用服务通过此接口访问数据存储
type OrderRepository interface {
	Create(ctx context.Context, order *domain.Order) error
	FindByOrderNumber(ctx context.Context, orderNumber string) (*domain.Order, error)
}

// CreateOrderRequest 创建订单请求（应用层 DTO）
type CreateOrderRequest struct {
	MerchantID   string              `validate:"required"`
	Items        []OrderItemRequest  `validate:"required,min=1,dive"`
	DeliveryInfo DeliveryInfoRequest `validate:"required"`
	Remark       string              `validate:"omitempty,max=200"`
}

// OrderItemRequest 订单项请求
type OrderItemRequest struct {
	DishID   string  `validate:"required"`
	DishName string  `validate:"required"`
	Quantity int     `validate:"required,gt=0"`
	Price    float64 `validate:"required,gt=0"`
}

// DeliveryInfoRequest 配送信息请求
type DeliveryInfoRequest struct {
	RecipientName  string `validate:"required"`
	RecipientPhone string `validate:"required,phone"`
	Address        string `validate:"required"`
}

// OrderData 订单数据（应用层 DTO）
type OrderData struct {
	OrderNumber string
	Status      string
	Pricing     PricingInfo
	CreatedAt   string
}

// PricingInfo 价格信息
type PricingInfo struct {
	ItemsTotal   string
	PackagingFee string
	DeliveryFee  string
	FinalAmount  string
}
