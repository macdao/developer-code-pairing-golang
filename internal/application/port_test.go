package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCreateOrderRequest_Validate_Success 测试验证成功
func TestCreateOrderRequest_Validate_Success(t *testing.T) {
	// Arrange
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
			Address:        "北京市朝阳区xxx",
		},
		Remark: "少辣",
	}

	// Act
	err := Validator.Struct(req)

	// Assert
	assert.NoError(t, err)
}

// TestCreateOrderRequest_Validate_EmptyMerchantID 测试空商家ID
func TestCreateOrderRequest_Validate_EmptyMerchantID(t *testing.T) {
	// Arrange
	req := &CreateOrderRequest{
		MerchantID: "", // 空商家ID
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
	err := Validator.Struct(req)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "MerchantID")
}

// TestCreateOrderRequest_Validate_EmptyItems 测试空订单项
func TestCreateOrderRequest_Validate_EmptyItems(t *testing.T) {
	// Arrange
	req := &CreateOrderRequest{
		MerchantID: "merchant_001",
		Items:      []OrderItemRequest{}, // 空订单项
		DeliveryInfo: DeliveryInfoRequest{
			RecipientName:  "张三",
			RecipientPhone: "13800138000",
			Address:        "北京市朝阳区xxx",
		},
	}

	// Act
	err := Validator.Struct(req)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Items")
}

// TestCreateOrderRequest_Validate_LongRemark 测试超长备注
func TestCreateOrderRequest_Validate_LongRemark(t *testing.T) {
	// Arrange
	longRemark := ""
	for i := 0; i < 201; i++ {
		longRemark += "a"
	}

	req := &CreateOrderRequest{
		MerchantID: "merchant_001",
		Items: []OrderItemRequest{
			{DishID: "dish_001", DishName: "宫保鸡丁", Quantity: 1, Price: 28.00},
		},
		DeliveryInfo: DeliveryInfoRequest{
			RecipientName:  "张三",
			RecipientPhone: "13800138000",
			Address:        "北京市朝阳区xxx",
		},
		Remark: longRemark,
	}

	// Act
	err := Validator.Struct(req)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Remark")
}

// TestCreateOrderRequest_Validate_EmptyRemark 测试空备注（应该成功）
func TestCreateOrderRequest_Validate_EmptyRemark(t *testing.T) {
	// Arrange
	req := &CreateOrderRequest{
		MerchantID: "merchant_001",
		Items: []OrderItemRequest{
			{DishID: "dish_001", DishName: "宫保鸡丁", Quantity: 1, Price: 28.00},
		},
		DeliveryInfo: DeliveryInfoRequest{
			RecipientName:  "张三",
			RecipientPhone: "13800138000",
			Address:        "北京市朝阳区xxx",
		},
		Remark: "", // 空备注
	}

	// Act
	err := Validator.Struct(req)

	// Assert
	assert.NoError(t, err)
}

// TestOrderItemRequest_Validate_EmptyDishID 测试空餐品ID
func TestOrderItemRequest_Validate_EmptyDishID(t *testing.T) {
	// Arrange
	req := &CreateOrderRequest{
		MerchantID: "merchant_001",
		Items: []OrderItemRequest{
			{DishID: "", DishName: "宫保鸡丁", Quantity: 1, Price: 28.00}, // 空DishID
		},
		DeliveryInfo: DeliveryInfoRequest{
			RecipientName:  "张三",
			RecipientPhone: "13800138000",
			Address:        "北京市朝阳区xxx",
		},
	}

	// Act
	err := Validator.Struct(req)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "DishID")
}

// TestOrderItemRequest_Validate_InvalidQuantity 测试无效数量
func TestOrderItemRequest_Validate_InvalidQuantity(t *testing.T) {
	// Arrange
	req := &CreateOrderRequest{
		MerchantID: "merchant_001",
		Items: []OrderItemRequest{
			{DishID: "dish_001", DishName: "宫保鸡丁", Quantity: 0, Price: 28.00}, // 数量为0
		},
		DeliveryInfo: DeliveryInfoRequest{
			RecipientName:  "张三",
			RecipientPhone: "13800138000",
			Address:        "北京市朝阳区xxx",
		},
	}

	// Act
	err := Validator.Struct(req)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Quantity")
}

// TestOrderItemRequest_Validate_InvalidPrice 测试无效价格
func TestOrderItemRequest_Validate_InvalidPrice(t *testing.T) {
	// Arrange
	req := &CreateOrderRequest{
		MerchantID: "merchant_001",
		Items: []OrderItemRequest{
			{DishID: "dish_001", DishName: "宫保鸡丁", Quantity: 1, Price: 0}, // 价格为0
		},
		DeliveryInfo: DeliveryInfoRequest{
			RecipientName:  "张三",
			RecipientPhone: "13800138000",
			Address:        "北京市朝阳区xxx",
		},
	}

	// Act
	err := Validator.Struct(req)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Price")
}

// TestDeliveryInfoRequest_Validate_EmptyRecipientName 测试空收货人姓名
func TestDeliveryInfoRequest_Validate_EmptyRecipientName(t *testing.T) {
	// Arrange
	req := &CreateOrderRequest{
		MerchantID: "merchant_001",
		Items: []OrderItemRequest{
			{DishID: "dish_001", DishName: "宫保鸡丁", Quantity: 1, Price: 28.00},
		},
		DeliveryInfo: DeliveryInfoRequest{
			RecipientName:  "", // 空姓名
			RecipientPhone: "13800138000",
			Address:        "北京市朝阳区xxx",
		},
	}

	// Act
	err := Validator.Struct(req)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "RecipientName")
}

// TestDeliveryInfoRequest_Validate_InvalidPhone 测试无效手机号
func TestDeliveryInfoRequest_Validate_InvalidPhone(t *testing.T) {
	testCases := []struct {
		name  string
		phone string
	}{
		{"短手机号", "12345"},
		{"长手机号", "138001380001"},
		{"不以1开头", "23800138000"},
		{"空手机号", ""},
		{"第二位是0", "10800138000"},
		{"第二位是1", "11800138000"},
		{"第二位是2", "12800138000"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			req := &CreateOrderRequest{
				MerchantID: "merchant_001",
				Items: []OrderItemRequest{
					{DishID: "dish_001", DishName: "宫保鸡丁", Quantity: 1, Price: 28.00},
				},
				DeliveryInfo: DeliveryInfoRequest{
					RecipientName:  "张三",
					RecipientPhone: tc.phone,
					Address:        "北京市朝阳区xxx",
				},
			}

			// Act
			err := Validator.Struct(req)

			// Assert
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "RecipientPhone")
		})
	}
}

// TestDeliveryInfoRequest_Validate_ValidPhone 测试有效手机号
func TestDeliveryInfoRequest_Validate_ValidPhone(t *testing.T) {
	testCases := []struct {
		name  string
		phone string
	}{
		{"13开头", "13800138000"},
		{"14开头", "14800138000"},
		{"15开头", "15800138000"},
		{"16开头", "16800138000"},
		{"17开头", "17800138000"},
		{"18开头", "18800138000"},
		{"19开头", "19800138000"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			req := &CreateOrderRequest{
				MerchantID: "merchant_001",
				Items: []OrderItemRequest{
					{DishID: "dish_001", DishName: "宫保鸡丁", Quantity: 1, Price: 28.00},
				},
				DeliveryInfo: DeliveryInfoRequest{
					RecipientName:  "张三",
					RecipientPhone: tc.phone,
					Address:        "北京市朝阳区xxx",
				},
			}

			// Act
			err := Validator.Struct(req)

			// Assert
			assert.NoError(t, err)
		})
	}
}

// TestDeliveryInfoRequest_Validate_EmptyAddress 测试空地址
func TestDeliveryInfoRequest_Validate_EmptyAddress(t *testing.T) {
	// Arrange
	req := &CreateOrderRequest{
		MerchantID: "merchant_001",
		Items: []OrderItemRequest{
			{DishID: "dish_001", DishName: "宫保鸡丁", Quantity: 1, Price: 28.00},
		},
		DeliveryInfo: DeliveryInfoRequest{
			RecipientName:  "张三",
			RecipientPhone: "13800138000",
			Address:        "", // 空地址
		},
	}

	// Act
	err := Validator.Struct(req)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Address")
}

// TestCreateOrderRequest_Validate_MultipleErrors 测试多个验证错误
func TestCreateOrderRequest_Validate_MultipleErrors(t *testing.T) {
	// Arrange
	req := &CreateOrderRequest{
		MerchantID: "", // 错误1：空商家ID
		Items:      []OrderItemRequest{}, // 错误2：空订单项
		DeliveryInfo: DeliveryInfoRequest{
			RecipientName:  "", // 错误3：空姓名
			RecipientPhone: "123", // 错误4：无效手机号
			Address:        "", // 错误5：空地址
		},
	}

	// Act
	err := Validator.Struct(req)

	// Assert
	assert.Error(t, err)
	// 应该包含多个错误
}
