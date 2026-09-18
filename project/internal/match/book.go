package match

import (
	"slices"
	"sync"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Book struct {
	mu   sync.Mutex
	bids []*Order // 买盘：价格降序（出价高的在前）
	asks []*Order // 卖盘：价格升序（要价低的在前）
}

func New() *Book { return &Book{} }

// Snapshot 只读快照：持锁拷贝一份盘口交出去（handler 的 GET /orderbook 用它）
// 为什么不让外部直接读 b.bids / b.asks：① 封装（字段故意不导出）；② 读也要持锁；③ 别把内部切片头泄漏出去
func (b *Book) Snapshot() (bids, asks []*Order) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return slices.Clone(b.bids), slices.Clone(b.asks) // 需要 import "slices"
}

// 买盘降序
func (b *Book) insertBid(o *Order) {
	i := 0
	for i < len(b.bids) && b.bids[i].Price.GreaterThanOrEqual(o.Price) {
		i++
	}
	b.bids = slices.Insert(b.bids, i, o) // 等价于 append+copy+赋值，但不出错
}

// 卖盘升序
func (b *Book) insertAsk(o *Order) {
	i := 0
	for i < len(b.asks) && b.asks[i].Price.LessThanOrEqual(o.Price) {
		i++
	}
	b.asks = slices.Insert(b.asks, i, o)
}

// Submit 接收一笔新订单，内部撮合，返回本次产生的成交列表
func (b *Book) Submit(o *Order) []Trade {
	b.mu.Lock()
	defer b.mu.Unlock()

	var trades []Trade
	if o.Side == Buy {
		for len(b.asks) > 0 && o.Price.GreaterThanOrEqual(b.asks[0].Price) {
			best := b.asks[0]
			fill := decimal.Min(best.Amount, o.Amount)

			trades = append(trades, Trade{
				ID:        uuid.New().String(),
				BuyerID:   o.UserID,    // 主动方是买方
				SellerID:  best.UserID, // maker 是卖方
				Price:     best.Price,  // 成交价取挂单方价格
				Amount:    fill,
				TakerSide: o.Side,
			})

			o.Amount = o.Amount.Sub(fill)
			best.Amount = best.Amount.Sub(fill)
			if best.Amount.IsZero() {
				b.asks = b.asks[1:] // 吃完弹出队首
			}
			if o.Amount.IsZero() {
				break
			}
		}
		if !o.Amount.IsZero() {
			b.insertBid(o)
		}
	} else {
		for len(b.bids) > 0 && o.Price.LessThanOrEqual(b.bids[0].Price) {
			best := b.bids[0]
			fill := decimal.Min(best.Amount, o.Amount)

			trades = append(trades, Trade{
				ID:        uuid.New().String(),
				BuyerID:   best.UserID, // maker 是买方
				SellerID:  o.UserID,    // 主动方是卖方
				Price:     best.Price,
				Amount:    fill,
				TakerSide: o.Side,
			})

			o.Amount = o.Amount.Sub(fill)
			best.Amount = best.Amount.Sub(fill)
			if best.Amount.IsZero() {
				b.bids = b.bids[1:] // 吃完弹出队首
			}
			if o.Amount.IsZero() {
				break
			}
		}
		if !o.Amount.IsZero() {
			b.insertAsk(o)
		}
	}
	return trades
}
