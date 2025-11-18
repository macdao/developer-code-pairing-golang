package web

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

const (
	// UserIDKey Context 中存储用户ID的键
	UserIDKey = "userID"
)

// AuthMiddleware JWT 认证中间件
func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// 从 Authorization header 提取 token
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "missing authorization header",
			})
		}

		// 检查 Bearer 前缀
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.JSON(http.StatusUnauthorized, ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "invalid authorization header format",
			})
		}

		tokenString := parts[1]

		// 验证 token
		claims, err := ValidateToken(tokenString)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "invalid or expired token",
			})
		}

		// 将用户ID存入 Context
		c.Set(UserIDKey, claims.UserID)

		return next(c)
	}
}
