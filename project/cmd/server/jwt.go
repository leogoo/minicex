package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 密钥先放常量，生产要进环境变量/Secret（见 ⑤）
var jwtSecret = []byte("dev-only-change-me")

type claims struct {
	UserID string `json:"uid"`
	jwt.RegisteredClaims
}

func GenerateToken(userID string) (string, error) {
	// 1) 构造 claims：塞入 userID + 过期时间(如 2h)
	// 2) jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 3) .SignedString(jwtSecret)

	now := time.Now()
	c := claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(2000 * time.Hour)),
			Issuer:    "minicex",
		},
	}
	// 用 HS256 + 密钥签名
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	signed, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err) // %w 包装，保留错误链
	}
	return signed, nil
}

func ParseToken(token string) (string, error) {
	c := &claims{}
	parsed, err := jwt.ParseWithClaims(
		token,
		c,
		func(t *jwt.Token) (interface{}, error) {
			// 防 alg=none 降级攻击：只允许 HMAC 系列算法
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return jwtSecret, nil
		},
	)
	if err != nil {
		return "", fmt.Errorf("parse token: %w", err)
	}
	if !parsed.Valid {
		return "", errors.New("invalid token")
	}
	return c.UserID, nil
}
