package web

// CreateOrderRequest Web 层创建订单请求
type CreateOrderRequest struct {
	MerchantID   string              `json:"merchantId"`
	Items        []OrderItemRequest  `json:"items"`
	DeliveryInfo DeliveryInfoRequest `json:"deliveryInfo"`
	Remark       string              `json:"remark"`
}

// OrderItemRequest Web 层订单项请求
type OrderItemRequest struct {
	DishID   string  `json:"dishId"`
	DishName string  `json:"dishName"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

// DeliveryInfoRequest Web 层配送信息请求
type DeliveryInfoRequest struct {
	RecipientName  string `json:"recipientName"`
	RecipientPhone string `json:"recipientPhone"`
	Address        string `json:"address"`
}

// CreateOrderResponse 创建订单响应
type CreateOrderResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    *OrderData  `json:"data,omitempty"`
}

// OrderData 订单数据
type OrderData struct {
	OrderNumber string      `json:"orderNumber"`
	Status      string      `json:"status"`
	Pricing     PricingInfo `json:"pricing"`
	CreatedAt   string      `json:"createdAt"`
}

// PricingInfo 价格信息
type PricingInfo struct {
	ItemsTotal   string `json:"itemsTotal"`
	PackagingFee string `json:"packagingFee"`
	DeliveryFee  string `json:"deliveryFee"`
	FinalAmount  string `json:"finalAmount"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
}
