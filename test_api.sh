#!/bin/bash

# 订单服务 API 测试脚本

echo "=== 订单服务 API 测试 ==="
echo ""

# 生成测试 token
echo "1. 生成测试 JWT token (userID: 1001)..."
TOKEN=$(go run tools/generate_token.go 1001 | grep "eyJ" | head -1)
echo "Token: $TOKEN"
echo ""

# 测试创建订单
echo "2. 创建订单..."
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "merchantId": "merchant_001",
    "items": [
      {
        "dishId": "dish_001",
        "dishName": "宫保鸡丁",
        "quantity": 2,
        "price": 28.00
      },
      {
        "dishId": "dish_002",
        "dishName": "鱼香肉丝",
        "quantity": 1,
        "price": 26.00
      }
    ],
    "deliveryInfo": {
      "recipientName": "张三",
      "recipientPhone": "13800138000",
      "address": "北京市朝阳区xxx街道xxx号"
    },
    "remark": "少辣，多放醋"
  }' | jq .

echo ""
echo ""

# 测试验证失败
echo "3. 测试验证失败（缺少商家ID）..."
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "items": [
      {
        "dishId": "dish_001",
        "dishName": "宫保鸡丁",
        "quantity": 2,
        "price": 28.00
      }
    ],
    "deliveryInfo": {
      "recipientName": "张三",
      "recipientPhone": "13800138000",
      "address": "北京市朝阳区xxx街道xxx号"
    }
  }' | jq .

echo ""
echo ""

# 测试未授权
echo "4. 测试未授权（无 token）..."
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
    "merchantId": "merchant_001",
    "items": [
      {
        "dishId": "dish_001",
        "dishName": "宫保鸡丁",
        "quantity": 1,
        "price": 28.00
      }
    ],
    "deliveryInfo": {
      "recipientName": "张三",
      "recipientPhone": "13800138000",
      "address": "北京市朝阳区xxx街道xxx号"
    }
  }' | jq .

echo ""
echo "=== 测试完成 ==="
