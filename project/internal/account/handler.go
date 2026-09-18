package account

import (
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// 模块自己把路由挂上去，main.go 只调用这一行
func RegisterRoutes(r gin.IRouter, svc *Service) {
	h := NewHandler(svc)
	r.POST("/accounts", h.OpenAccount)
	r.GET("/accounts/:id/balance", h.GetBalance)
	r.POST("/accounts/:id/deposit", h.Deposit)
}

func (h *Handler) OpenAccount(c *gin.Context) {
	// 1) 解析请求体
	var req struct {
		UserID string `json:"user_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	acc := h.svc.OpenAccount(req.UserID)
	c.JSON(200, acc)
}
func (h *Handler) GetBalance(c *gin.Context) {
	id := c.Param("id")
	balance, err := h.svc.GetBalance(id)
	if err != nil {
		c.JSON(200, balance)
	}
	c.JSON(400, gin.H{"error": err.Error()})
}
func (h *Handler) Deposit(c *gin.Context) {
	var req struct {
		id   string          `json:"id"`
		amt  decimal.Decimal `json: "amt"`
		note string          `json:"note"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}
	entry, err := h.svc.Deposit(req.id, req.amt, req.note)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}
	c.JSON(400, entry)
}
