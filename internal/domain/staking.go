package domain

import (
	"dex/pkg/decimal"
	"time"
)

// Stake представляет собой позицию стейкинга (блокировки токенов ради получения APY).
type Stake struct {
	ID        string          `json:"id"`
	AccountID string          `json:"accountId"`
	Asset     string          `json:"asset"`
	Amount    decimal.Decimal `json:"amount"`
	APY       decimal.Decimal `json:"apy"` // Годовая процентная доходность (например, 0.12 для 12%)
	CreatedAt time.Time       `json:"createdAt"`
	UnlocksAt time.Time       `json:"unlocksAt"`
	IsActive  bool            `json:"isActive"`
}

// NewStake создает новую позицию стейкинга.
func NewStake(id, accountID, asset string, amount, apy decimal.Decimal, lockDuration time.Duration) *Stake {
	now := time.Now()
	return &Stake{
		ID:        id,
		AccountID: accountID,
		Asset:     asset,
		Amount:    amount,
		APY:       apy,
		CreatedAt: now,
		UnlocksAt: now.Add(lockDuration),
		IsActive:  true,
	}
}

// CalculateReward рассчитывает примерную прибыль на текущий момент времени.
func (s *Stake) CalculateReward() decimal.Decimal {
	if !s.IsActive {
		return decimal.Zero()
	}
	
	now := time.Now()
	if now.After(s.UnlocksAt) {
		now = s.UnlocksAt // Не начисляем сверх срока
	}
	
	durationSec := now.Sub(s.CreatedAt).Seconds()
	yearSec := (365 * 24 * time.Hour).Seconds()
	
	// Прибыль = Сумма * APY * (Прошедшее время / Год)
	// Для HFT систем APY может рассчитываться по секундам (Continuous compounding), но здесь линейно.
	fraction := decimal.NewFromFloat(durationSec / yearSec)
	reward := s.Amount.Mul(s.APY).Mul(fraction)
	return reward
}
