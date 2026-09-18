package match

import (
	"time"

	"github.com/shopspring/decimal"
)

type Side string

const (
	Buy  Side = "buy"
	Sell Side = "sell"
)

type Order struct {
	ID        string
	UserID    string
	Side      Side
	Price     decimal.Decimal // 限价单的挂单价格
	Amount    decimal.Decimal // 剩余未成交数量（撮合中会被扣减）
	CreatedAt time.Time
}

type Trade struct {
	ID        string
	BuyerID   string
	SellerID  string
	Price     decimal.Decimal // 成交价（取挂单方价格，price-time 优先）
	Amount    decimal.Decimal // 成交量
	TakerSide Side            // 主动方方向，记录谁发动了这笔成交
}
