package domain

import (
	"dex/pkg/decimal"
)

// MarginMode определяет тип маржи для фьючерсной позиции
type MarginMode string

const (
	MarginIsolated MarginMode = "ISOLATED"
	MarginCross    MarginMode = "CROSS"
)

// Position представляет собой открытую фьючерсную позицию пользователя
type Position struct {
	ID               string          `json:"id"`
	AccountID        string          `json:"accountId"`
	Pair             Pair            `json:"pair"`             // Например, BTC/USDT
	Side             Side            `json:"side"`             // BUY (Long) или SELL (Short)
	MarginMode       MarginMode      `json:"marginMode"`       // ISOLATED или CROSS
	Leverage         int             `json:"leverage"`         // Плечо, например 10, 20, 50, 100
	EntryPrice       decimal.Decimal `json:"entryPrice"`       // Средняя цена входа
	Size             decimal.Decimal `json:"size"`             // Размер позиции в базовом активе (например, 1.5 BTC)
	Margin           decimal.Decimal `json:"margin"`           // Изолированная маржа (замороженные средства)
	LiquidationPrice decimal.Decimal `json:"liquidationPrice"` // Цена, при которой позиция будет принудительно ликвидирована
}

// NewPosition создает новую пустую позицию (вызывается при первом открытии)
func NewPosition(id, accountID string, pair Pair, side Side, marginMode MarginMode, leverage int) *Position {
	return &Position{
		ID:         id,
		AccountID:  accountID,
		Pair:       pair,
		Side:       side,
		MarginMode: marginMode,
		Leverage:   leverage,
		EntryPrice: decimal.Zero(),
		Size:       decimal.Zero(),
		Margin:     decimal.Zero(),
	}
}
