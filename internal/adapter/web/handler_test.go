package web

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"order-service/internal/application"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// testValidator 测试用验证器（不做实际验证）
type testValidator struct{}

func (tv *testValidator) Validate(i interface{}) error {
	return nil
}

// MockOrderService 模拟 OrderService
type MockOrderService struct {
	mock.Mock
}

func (m *MockOrderService) CreateOrder(ctx context.Context, userID uint64, req *application.CreateOrderRequest) (*application.OrderData, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*application.OrderData), args.Error(1)
}

func TestOrderHandler_CreateOrder_Success(t *testing.T) {
	e := echo.New()
	e.Validator = &testValidator{}
	mockService := new(MockOrderService)
	handler := NewOrderHandler(mockService)

	// 准备请求
	reqBody := CreateOrderRequest{
		MerchantID: "merchant1",
		Items: []OrderItemRequest{
			{
				DishID:   "dish1",
				DishName: "宫保鸡丁",
				Quantity: 2,
				Price:    28.00,
			},
		},
		DeliveryInfo: DeliveryInfoRequest{
			RecipientName:  "张三",
			RecipientPhone: "13800138000",
			Address:        "北京市朝阳区xxx",
		},
		Remark: "少辣",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set(UserIDKey, uint64(1001))

	// 设置 mock 期望
	expectedOrderData := &application.OrderData{
		OrderNumber: "20241117120000123456",
		Status:      "PENDING_PAYMENT",
		Pricing: application.PricingInfo{
			ItemsTotal:   "56.00",
			PackagingFee: "1.00",
			DeliveryFee:  "3.00",
			FinalAmount:  "60.00",
		},
		CreatedAt: "2024-11-17T12:00:00Z",
	}
	mockService.On("CreateOrder", mock.Anything, uint64(1001), mock.AnythingOfType("*application.CreateOrderRequest")).Return(expectedOrderData, nil)

	// 执行
	err := handler.CreateOrder(c)

	// 验证
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var response CreateOrderResponse
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, http.StatusCreated, response.Code)
	assert.Equal(t, "20241117120000123456", response.Data.OrderNumber)
	mockService.AssertExpectations(t)
}

func TestOrderHandler_CreateOrder_ValidationError(t *testing.T) {
	e := echo.New()
	e.Validator = &testValidator{}
	mockService := new(MockOrderService)
	handler := NewOrderHandler(mockService)

	// 准备请求（缺少必填字段）
	reqBody := CreateOrderRequest{
		Items: []OrderItemRequest{},
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set(UserIDKey, uint64(1001))

	// 设置 mock 期望
	mockService.On("CreateOrder", mock.Anything, uint64(1001), mock.AnythingOfType("*application.CreateOrderRequest")).
		Return(nil, application.NewValidationError("merchant_id", "商家ID不能为空"))

	// 执行
	err := handler.CreateOrder(c)

	// 验证
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var response ErrorResponse
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Contains(t, response.Message, "商家ID不能为空")
}

func TestOrderHandler_CreateOrder_Unauthorized(t *testing.T) {
	e := echo.New()
	mockService := new(MockOrderService)
	handler := NewOrderHandler(mockService)

	reqBody := CreateOrderRequest{
		MerchantID: "merchant1",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// 不设置 UserID

	// 执行
	err := handler.CreateOrder(c)

	// 验证
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestOrderHandler_CreateOrder_InternalError(t *testing.T) {
	e := echo.New()
	e.Validator = &testValidator{}
	mockService := new(MockOrderService)
	handler := NewOrderHandler(mockService)

	reqBody := CreateOrderRequest{
		MerchantID: "merchant1",
		Items: []OrderItemRequest{
			{
				DishID:   "dish1",
				DishName: "宫保鸡丁",
				Quantity: 2,
				Price:    28.00,
			},
		},
		DeliveryInfo: DeliveryInfoRequest{
			RecipientName:  "张三",
			RecipientPhone: "13800138000",
			Address:        "北京市朝阳区xxx",
		},
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set(UserIDKey, uint64(1001))

	// 设置 mock 期望
	mockService.On("CreateOrder", mock.Anything, uint64(1001), mock.AnythingOfType("*application.CreateOrderRequest")).
		Return(nil, application.NewInternalError("database error", nil))

	// 执行
	err := handler.CreateOrder(c)

	// 验证
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
