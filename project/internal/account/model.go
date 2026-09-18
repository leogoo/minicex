package account

import (
	"time"

	"github.com/shopspring/decimal"
)

type Account struct {
	ID       string
	UserID   string
	Balance  decimal.Decimal
	Currency string
	CreateAt time.Time
}

type LedgerEntry struct {
	ID        string
	AccountID string
	Amount    decimal.Decimal
	Type      string
	Note      string
	CreateAt  time.Time
}
