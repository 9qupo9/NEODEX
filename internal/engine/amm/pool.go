package amm

import (
	"dex/internal/domain"
	"dex/pkg/decimal"
	"errors"
	"sync"
)

// Pool представляет собой пул ликвидности для автоматического маркет-мейкера (AMM).
// Работает по формуле константного произведения (x * y = k), как Uniswap v2.
// В гибридном движке, если в стакане нет ликвидности, ордера могут быть маршрутизированы в этот пул.
type Pool struct {
	Pair     domain.Pair
	mu       sync.RWMutex
	BaseRes  decimal.Decimal // Резерв базового актива (например, количество BTC в пуле)
	QuoteRes decimal.Decimal // Резерв котируемого актива (например, количество USDT в пуле)
}

// NewPool создает пустой пул ликвидности.
func NewPool(pair domain.Pair) *Pool {
	return &Pool{
		Pair:     pair,
		BaseRes:  decimal.Zero(),
		QuoteRes: decimal.Zero(),
	}
}

// AddLiquidity добавляет ликвидность в пул (зачисление токенов провайдером ликвидности - LP).
// TODO: Возвращать LP токены пользователю (Share Tokens).
func (p *Pool) AddLiquidity(base, quote decimal.Decimal) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.BaseRes = p.BaseRes.Add(base)
	p.QuoteRes = p.QuoteRes.Add(quote)
}

// SwapExactIn рассчитывает количество получаемого актива (amountOut) за заданное количество входящего актива.
// Использует формулу y_out = y_res - (x_res * y_res) / (x_res + x_in).
// Для простоты прототипа комиссия LP провайдерам пока не учитывается (нет fee 0.3%).
// TODO: Добавить расчет Slippage (проскальзывания) и Price Impact.
func (p *Pool) SwapExactIn(isBaseIn bool, amountIn decimal.Decimal) (decimal.Decimal, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Если пул пустой, своп невозможен
	if p.BaseRes.IsZero() || p.QuoteRes.IsZero() {
		return decimal.Zero(), errors.New("недостаточно ликвидности в AMM пуле")
	}

	var amountOut decimal.Decimal

	if isBaseIn {
		// Пользователь отдает Base (например, продает BTC), получает Quote (USDT)
		k := p.BaseRes.Mul(p.QuoteRes)
		newBase := p.BaseRes.Add(amountIn)
		newQuote := k.Div(newBase)
		amountOut = p.QuoteRes.Sub(newQuote)

		// Обновляем резервы
		p.BaseRes = newBase
		p.QuoteRes = newQuote
	} else {
		// Пользователь отдает Quote (покупает за USDT), получает Base (BTC)
		k := p.BaseRes.Mul(p.QuoteRes)
		newQuote := p.QuoteRes.Add(amountIn)
		newBase := k.Div(newQuote)
		amountOut = p.BaseRes.Sub(newBase)

		// Обновляем резервы
		p.BaseRes = newBase
		p.QuoteRes = newQuote
	}

	return amountOut, nil
}
