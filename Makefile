.PHONY: help build run test test-coverage clean deps fmt vet lint install-tools generate-token

# 变量定义
BINARY_NAME=order-service
MAIN_PATH=./cmd/server
BUILD_DIR=./bin
GO=go
GOTEST=$(GO) test
GOVET=$(GO) vet
GOFMT=gofmt

# 默认目标
.DEFAULT_GOAL := help

# 帮助信息
help: ## 显示帮助信息
	@echo "可用的 make 命令："
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

# 安装依赖
deps: ## 下载项目依赖
	@echo "正在下载依赖..."
	$(GO) mod download
	$(GO) mod tidy

# 构建
build: ## 编译项目
	@echo "正在编译..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "编译完成: $(BUILD_DIR)/$(BINARY_NAME)"

# 运行
run: ## 运行服务
	@echo "正在启动服务..."
	$(GO) run $(MAIN_PATH)

# 测试
test: ## 运行所有测试
	@echo "正在运行测试..."
	$(GOTEST) -v ./...

# 测试覆盖率
test-coverage: ## 运行测试并生成覆盖率报告
	@echo "正在生成测试覆盖率报告..."
	$(GOTEST) -cover ./...
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "覆盖率报告已生成: coverage.html"

# 代码格式化
fmt: ## 格式化代码
	@echo "正在格式化代码..."
	$(GOFMT) -w .

# 代码检查
vet: ## 运行 go vet
	@echo "正在运行 go vet..."
	$(GOVET) ./...

# 清理
clean: ## 清理构建产物
	@echo "正在清理..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "清理完成"

# 生成测试 Token
generate-token: ## 生成测试 JWT Token (使用方式: make generate-token USER_ID=1001)
	@if [ -z "$(USER_ID)" ]; then \
		echo "请指定 USER_ID，例如: make generate-token USER_ID=1001"; \
	else \
		$(GO) run tools/generate_token.go $(USER_ID); \
	fi

# 完整构建流程
all: clean deps fmt vet test build ## 执行完整的构建流程

# 开发模式（热重载需要额外工具）
dev: ## 开发模式运行（需要安装 air）
	@which air > /dev/null || (echo "请先安装 air: go install github.com/air-verse/air@latest" && exit 1)
	air
