# 订单服务 (Order Service)

基于 Go 和 Echo 框架的订餐系统订单创建服务，采用六边形架构和 DDD 原则实现。

## 技术栈

- Go 1.25+
- Echo v4 (Web 框架)
- JWT (认证)
- go-playground/validator (验证)
- shopspring/decimal (精度计算)
- testify (测试)
- 内存存储

## 架构设计

采用六边形架构（Hexagonal Architecture）：

- **Domain Layer**: 核心业务逻辑和领域模型
- **Application Layer**: 应用服务和端口定义
- **Adapter Layer**: 外部适配器（Web、持久化等）

依赖方向：Adapter → Application → Domain

### 项目结构

```
order-service/
├── cmd/server/main.go           # 应用入口
├── internal/
│   ├── domain/                  # 领域层
│   ├── application/             # 应用层
│   └── adapter/                 # 适配器层
│       ├── web/                 # Web 适配器
│       └── persistence/         # 持久化适配器
├── tools/                       # 工具脚本
└── README.md
```

## 快速开始

### 1. 安装依赖

```bash
go mod download
```

### 2. 运行服务

```bash
go run ./cmd/server
```

服务启动在 `http://localhost:8080`

### 3. 运行测试

```bash
# 运行所有测试
go test ./...

# 查看覆盖率
go test ./... -cover
```

## API 使用

### 1. 生成测试 Token

```bash
go run tools/generate_token.go 1001
```

### 2. 创建订单

```bash
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "merchantId": "merchant_001",
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
    },
    "remark": "少辣"
  }'
```

### 3. 使用测试脚本

```bash
# 启动服务
go run ./cmd/server

# 运行测试脚本（另一个终端）
./test_api.sh
```

## 测试

```bash
# 运行所有测试
go test ./...

# 查看覆盖率
go test ./... -cover

# 详细输出
go test ./... -v
```
