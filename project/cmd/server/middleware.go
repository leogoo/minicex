package main

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1) h := c.GetHeader("Authorization")，去掉 "Bearer " 前缀
		// 2) uid, err := ParseToken(token)；err != nil → c.AbortWithStatusJSON(401, ...) 并 return
		// 3) c.Set("user_id", uid)   // 注入请求上下文
		// 4) c.Next()                // 放行给下一个 handler

		raw := c.GetHeader("Authorization")
		const prefix = "Bearer "
		if !strings.HasPrefix(raw, prefix) {
			c.AbortWithStatusJSON(401, gin.H{"error": "missing or malformed token"})
			return
		}
		token := raw[len(prefix):] // 去掉 "Bearer " 前缀
		uid, err := ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": err.Error()})
			return
		}
		c.Set("user_id", uid)
		c.Next()
	}
}
