package account

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/shopspring/decimal"
)

type Service struct {
	mu       sync.Mutex
	accounts map[string]*Account
	ledger   []LedgerEntry
}

func New() *Service {
	return &Service{accounts: make(map[string]*Account), ledger: []LedgerEntry{}}
}

func (s *Service) OpenAccount(userID string) *Account {
	s.mu.Lock()
	defer s.mu.Unlock()

	acc := &Account{
		ID:       uuid.New().String(),
		UserID:   userID,
		Balance:  decimal.Zero,
		Currency: "USDT",
		CreateAt: time.Now(),
	}
	s.accounts[acc.ID] = acc
	return acc
}

func (s *Service) GetBalance(id string) (decimal.Decimal, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	acc, ok := s.accounts[id]
	if !ok {
		return decimal.Decimal{}, fmt.Errorf("account not found: %s", id)
	}
	return acc.Balance, nil
}

func (s *Service) Deposit(id string, amt decimal.Decimal, note string) (*LedgerEntry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	acc, ok := s.accounts[id]
	if !ok {
		return nil, fmt.Errorf("account not found: %s", id)
	}
	// 充值金额必须为正：decimal 用 IsNegative/IsZero 判断，别用 < 0
	if amt.IsNegative() || amt.IsZero() {
		return nil, fmt.Errorf("deposit amount must be positive: %s", amt)
	}
	acc.Balance = acc.Balance.Add(amt)

	ledger := LedgerEntry{
		ID:        uuid.New().String(),
		AccountID: id,
		Amount:    amt,
		Type:      "Deposit",
		Note:      note,
		CreateAt:  time.Now(),
	}
	s.ledger = append(s.ledger, ledger)

	return &ledger, nil
}
