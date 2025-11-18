package web

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	// JWTSecret JWT 密钥（生产环境应从配置读取）
	JWTSecret = "interview-demo-secret-key"
	// TokenExpiration token 过期时间
	TokenExpiration = 24 * time.Hour
)

// Claims JWT Claims
type Claims struct {
	UserID uint64 `json:"userId"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT token（用于测试）
func GenerateToken(userID uint64) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(JWTSecret))
}

// ValidateToken 验证 token 有效性
func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// ExtractUserID 从 token 提取用户ID
func ExtractUserID(tokenString string) (uint64, error) {
	claims, err := ValidateToken(tokenString)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}
