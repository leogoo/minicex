package match

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Handler struct{ book *Book }

func NewHandler(b *Book) *Handler { return &Handler{book: b} }

func (h *Handler) submitOrder(c *gin.Context) {
	// 1) 解析 JSON：user_id, side, price, amount
	// 2) 构造 *match.Order（ID 用 uuid，Amount 用 decimal.NewFromString）
	// 3) trades := h.book.Submit(order)
	// 4) c.JSON(200, gin.H{"order": order, "trades": trades})
	var req struct {
		UserID string `json:"user_id"`
		Side   Side   `json:"side"`
		Price  string `json:"price"`
		Amount string `json:"amount"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}
	price, err := decimal.NewFromString(req.Price)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid price"})
		return
	}
	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid price"})
		return
	}
	order := &Order{
		ID:        uuid.New().String(),
		UserID:    req.UserID,
		Side:      req.Side,
		Price:     price,
		Amount:    amount,
		CreatedAt: time.Now(),
	}
	trades := h.book.Submit(order)
	c.JSON(200, gin.H{"order": order, "trades": trades})
}

func (h *Handler) getBook(c *gin.Context) {
	bids, asks := h.book.Snapshot()
	c.JSON(200, gin.H{"bids": bids, "asks": asks})
}

func RegisterRoutes(r gin.IRouter, b *Book) {
	h := NewHandler(b)
	r.POST("/orders", h.submitOrder)
	r.GET("/orderbook", h.getBook)
}
