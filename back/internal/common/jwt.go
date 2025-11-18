package common

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JwtSecret = []byte(getJWTSecret())

// getJWTSecret 获取JWT密钥
func getJWTSecret() string {
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		return secret
	}
	// 默认密钥，生产环境中应该使用环境变量
	return "KJH87sd98HDS8df79SDF98sd8F7SDF8sd7FSD8fsd8F7sd8f7"
}

// 自定义 JWT payload
type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// 生成 JWT
func GenerateToken(userID uint) (string, error) {
	now := time.Now()
	expireTime := now.Add(24 * time.Hour)

	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtSecret)
}

// 解析 Token
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return JwtSecret, nil
		})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}
