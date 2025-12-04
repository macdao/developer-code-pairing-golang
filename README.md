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

```bash
# 查看所有可用命令
make help

# 安装依赖
make deps

# 运行服务（启动在 http://localhost:8080）
make run

# 运行测试
make test

# 查看测试覆盖率
make test-coverage

# 构建二进制文件
make build

# 生成测试 Token
make generate-token USER_ID=1001

# 完整构建流程（清理、依赖、格式化、检查、测试、构建）
make all
```

## API 使用

### 1. 生成测试 Token

```bash
make generate-token USER_ID=1001
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
make run

# 运行测试脚本（另一个终端）
./test_api.sh
```

## Makefile 命令说明

| 命令 | 说明 |
|------|------|
| `make help` | 显示所有可用命令 |
| `make deps` | 下载项目依赖 |
| `make build` | 编译项目到 bin/ 目录 |
| `make run` | 运行服务 |
| `make test` | 运行所有测试 |
| `make test-coverage` | 生成测试覆盖率报告 |
| `make fmt` | 格式化代码 |
| `make vet` | 运行 go vet 检查 |
| `make clean` | 清理构建产物 |
| `make generate-token` | 生成测试 JWT Token |
| `make all` | 执行完整构建流程 |
