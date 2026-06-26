package engine

import (
	"encoding/json"
	"errors"
	"dex/internal/domain"
	"dex/internal/engine/matching"
	"dex/internal/engine/orderbook"
	"dex/internal/storage"
	"dex/pkg/decimal"
	"sync"
)

var ErrInsufficientMargin = errors.New("insufficient margin")

// FuturesEngine — специализированное торговое ядро для деривативов.
// Отвечает за логику блокировки маржи с учетом плеча и обновление позиций.
type FuturesEngine struct {
	Pair     domain.Pair
	Store    storage.Store
	Book     *orderbook.Book
	Matcher  *matching.Matcher
	mu       sync.Mutex
	wsOut    chan []byte
}

func NewFuturesEngine(pair domain.Pair, store storage.Store, wsOut chan []byte) *FuturesEngine {
	book := orderbook.NewBook(pair)
	matcher := matching.NewMatcher(book)

	return &FuturesEngine{
		Pair:    pair,
		Store:   store,
		Book:    book,
		Matcher: matcher,
		wsOut:   wsOut,
	}
}

// PlaceOrder размещает фьючерсный ордер, проверяет обеспечение с учетом плеча.
func (e *FuturesEngine) PlaceOrder(order *domain.Order) ([]*domain.Trade, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	acc, err := e.Store.GetAccount(order.AccountID)
	if err != nil {
		return nil, err
	}

	// 1. Проверяем маржу (простейшая логика начальной маржи)
	// Номинальная стоимость = Цена * Количество
	// Требуемая маржа = Номинальная стоимость / Плечо
	notionalValue := order.Price.Mul(order.Qty)
	leverageDec := decimal.NewFromInt(int64(order.Leverage))
	requiredMargin := notionalValue.Div(leverageDec)

	// Фьючерсы торгуются и обеспечены QuoteAsset (USDT/USDC)
	if err := acc.LockFunds(e.Pair.QuoteAsset, requiredMargin); err != nil {
		return nil, ErrInsufficientMargin
	}

	// 2. Отправляем в движок сведения (Matcher)
	trades := e.Matcher.ProcessLimitOrder(order)
	_ = e.Store.SaveOrder(order)

	// 3. Обработка трейдов (Обновление позиций вместо расчетов спота)
	// Для упрощения: мы пока просто замораживаем маржу. Полноценный клиринг 
	// требует открытия/закрытия позиции и начисления реализованного PnL.
	for _, t := range trades {
		// Сохраняем трейд в историю
		e.Store.SaveTrade(t)
		
		// TODO: Интегрировать логику создания/усреднения domain.Position
		// В этой итерации мы фокусируемся на возможности выставления ордера.
	}

	// 4. WebSocket стакан (Futures)
	if e.wsOut != nil {
		depth := e.Book.GetDepth(50)
		msg := map[string]interface{}{
			"type":   "futures_orderbook",
			"symbol": "BTC_USDT_PERP",
			"data":   depth,
		}
		if b, err := json.Marshal(msg); err == nil {
			select {
			case e.wsOut <- b:
			default:
			}
		}
	}

	return trades, nil
}
