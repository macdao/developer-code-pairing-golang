package main

import (
	"log"

	"order-service/internal/adapter/persistence"
	"order-service/internal/adapter/web"
	"order-service/internal/application"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// 1. 初始化 Repository
	repo := persistence.NewInMemoryOrderRepository()

	// 2. 初始化 Application Service
	orderService := application.NewOrderService(repo)

	// 3. 初始化 Handler
	orderHandler := web.NewOrderHandler(orderService)

	// 4. 创建 Echo 实例
	e := echo.New()

	// 5. 配置中间件
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// 6. 注册路由
	api := e.Group("/api/v1")
	api.POST("/orders", orderHandler.CreateOrder, web.AuthMiddleware)

	// 7. 启动服务器
	log.Println("Starting server on :8080")
	if err := e.Start(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
