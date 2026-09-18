package main

import (
	"minicex/internal/account"
	"minicex/internal/match"

	"github.com/gin-gonic/gin"
)

func main() {
	svc := account.New()
	book := match.New()
	r := gin.Default()

	// 临时登录：拿到 token 用，后面接真实用户库
	r.POST("/login", func(c *gin.Context) {
		var req struct {
			UserID string `json:"user_id"`
		}
		c.ShouldBindJSON(&req)
		token, _ := GenerateToken(req.UserID)
		c.JSON(200, gin.H{"token": token})
	})

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "pong"})
	})

	// 受保护组：下面所有路由都要过鉴权
	auth := r.Group("/api")
	auth.Use(AuthRequired())
	account.RegisterRoutes(auth, svc)
	match.RegisterRoutes(auth, book)
	r.Run(":8088")
}
