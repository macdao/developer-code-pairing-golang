package web

import (
	"net/http"

	"order-service/internal/application"

	"github.com/labstack/echo/v4"
)

// OrderHandler 订单 HTTP 处理器
type OrderHandler struct {
	orderService application.OrderService
}

// NewOrderHandler 创建订单处理器
func NewOrderHandler(orderService application.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// CreateOrder 创建订单 HTTP 处理器
func (h *OrderHandler) CreateOrder(c echo.Context) error {
	// 1. 从 Context 获取用户ID
	userID, ok := c.Get(UserIDKey).(uint64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "user not authenticated",
		})
	}

	// 2. 解析请求体
	var webReq CreateOrderRequest
	if err := c.Bind(&webReq); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "invalid request body",
		})
	}

	// 3. 转换 Web DTO 到应用层 DTO
	appReq := h.convertToApplicationDTO(&webReq)

	// 4. 调用应用服务（验证在应用层完成）
	orderData, err := h.orderService.CreateOrder(c.Request().Context(), userID, appReq)
	if err != nil {
		return h.handleError(c, err)
	}

	// 5. 返回成功响应
	return c.JSON(http.StatusCreated, CreateOrderResponse{
		Code:    http.StatusCreated,
		Message: "order created successfully",
		Data: &OrderData{
			OrderNumber: orderData.OrderNumber,
			Status:      orderData.Status,
			Pricing: PricingInfo{
				ItemsTotal:   orderData.Pricing.ItemsTotal,
				PackagingFee: orderData.Pricing.PackagingFee,
				DeliveryFee:  orderData.Pricing.DeliveryFee,
				FinalAmount:  orderData.Pricing.FinalAmount,
			},
			CreatedAt: orderData.CreatedAt,
		},
	})
}

// convertToApplicationDTO 转换 Web DTO 到应用层 DTO
func (h *OrderHandler) convertToApplicationDTO(webReq *CreateOrderRequest) *application.CreateOrderRequest {
	items := make([]application.OrderItemRequest, len(webReq.Items))
	for i, item := range webReq.Items {
		items[i] = application.OrderItemRequest{
			DishID:   item.DishID,
			DishName: item.DishName,
			Quantity: item.Quantity,
			Price:    item.Price,
		}
	}

	return &application.CreateOrderRequest{
		MerchantID: webReq.MerchantID,
		Items:      items,
		DeliveryInfo: application.DeliveryInfoRequest{
			RecipientName:  webReq.DeliveryInfo.RecipientName,
			RecipientPhone: webReq.DeliveryInfo.RecipientPhone,
			Address:        webReq.DeliveryInfo.Address,
		},
		Remark: webReq.Remark,
	}
}

// handleError 处理不同类型的错误
func (h *OrderHandler) handleError(c echo.Context, err error) error {
	switch e := err.(type) {
	case *application.ValidationError:
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: e.Message,
			Field:   e.Field,
		})
	case *application.NotFoundError:
		return c.JSON(http.StatusNotFound, ErrorResponse{
			Code:    http.StatusNotFound,
			Message: e.Message,
		})
	case *application.InternalError:
		// 记录详细错误日志（生产环境应使用日志库）
		c.Logger().Error(e)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
	default:
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
	}
}
