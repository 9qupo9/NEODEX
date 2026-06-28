package engine

import (
	"encoding/json"
	"dex/internal/domain"
	"dex/internal/engine/matching"
	"dex/internal/engine/orderbook"
	"dex/internal/engine/settlement"
	"dex/internal/storage"
	"dex/pkg/decimal"
	"sync"
)

// Engine — это основной фасад и координатор всей торговой пары на бирже.
// В HFT-системах каждый Engine (каждая пара, например BTC/USDT) работает в отдельной горутине, 
// чтобы не блокировать другие пары.
type Engine struct {
	Pair     domain.Pair                  // Торговая пара этого движка
	Store    storage.Store                // Ссылка на глобальное хранилище балансов
	Book     *orderbook.Book              // Книга ордеров (стакан)
	Matcher      *matching.Matcher            // Движок сведения
	Settler      *settlement.Settler          // Движок расчетов (клиринга)
	mu           sync.Mutex                   // Глобальный мьютекс движка для строгой детерминированности обработки
	wsOut        chan []byte                  // Канал для отправки JSON-сообщений в WebSocket Hub
	conditionals []*domain.Order              // Очередь условных ордеров (StopLimit, TakeProfit)
	lastPrice    decimal.Decimal              // Последняя цена сделки (рыночная цена)
}

// NewEngine инициализирует и запускает все внутренние подсистемы торговой пары.
func NewEngine(pair domain.Pair, store storage.Store, wsOut chan []byte) *Engine {
	book := orderbook.NewBook(pair)
	matcher := matching.NewMatcher(book)
	
	// Очередь на 10,000 сделок (буфер для асинхронного клиринга)
	q := settlement.NewQueue(10000)
	settler := settlement.NewSettler(store, q)
	
	// Запускаем фоновый воркер клиринга
	settler.Start()

	return &Engine{
		Pair:         pair,
		Store:        store,
		Book:         book,
		Matcher:      matcher,
		Settler:      settler,
		wsOut:        wsOut,
		conditionals: make([]*domain.Order, 0),
		lastPrice:    decimal.Zero(),
	}
}

// PlaceOrder принимает новый ордер, валидирует средства, замораживает их и отправляет в Matcher.
// Это главная точка входа для всех торговых запросов (из HTTP, WS или TCP).
// Возвращает массив сделок (Trades), если они состоялись, и ошибку в случае нехватки средств.
func (e *Engine) PlaceOrder(order *domain.Order) ([]*domain.Trade, error) {
	// Блокируем движок. В однопоточной (для конкретной пары) архитектуре это гарантирует 
	// абсолютную хронологическую правильность исполнения заявок без состояния гонки.
	e.mu.Lock()
	defer e.mu.Unlock()

	// 1. Проверяем баланс пользователя
	acc, err := e.Store.GetAccount(order.AccountID)
	if err != nil {
		return nil, err
	}

	// 2. Замораживаем средства до исполнения или отмены
	if order.Side == domain.Buy {
		// Для покупки нужен QuoteAsset (например, USDT). Замораживаем (Цена * Количество)
		requiredQuote := order.Price.Mul(order.Qty)
		if err := acc.LockFunds(e.Pair.QuoteAsset, requiredQuote); err != nil {
			return nil, err
		}
	} else {
		// Для продажи нужен BaseAsset (например, BTC). Замораживаем Количество (Qty)
		if err := acc.LockFunds(e.Pair.BaseAsset, order.Qty); err != nil {
			return nil, err
		}
	}

	// 3. Если это условный ордер, добавляем в очередь и выходим
	if order.Type == domain.StopLimit || order.Type == domain.TakeProfit {
		e.conditionals = append(e.conditionals, order)
		_ = e.Store.SaveOrder(order)
		return nil, nil // Сделок пока нет
	}

	// 4. Отправляем в Matcher (сведение)
	trades := e.Matcher.ProcessLimitOrder(order)
	
	// Обновляем последнюю цену
	if len(trades) > 0 {
		e.lastPrice = trades[len(trades)-1].Price
		e.triggerConditionals() // Проверяем отложенные ордера
	}

	// 5. Сохраняем текущее состояние ордера (Filled/PartiallyFilled/New)
	_ = e.Store.SaveOrder(order)
	
	// 6. Все произошедшие сделки отправляем в фоновый клиринг (списание балансов)
	for _, t := range trades {
		e.Settler.Queue.Enqueue(t)
	}

	// 6. Рассылаем новый стакан по WebSocket
	if e.wsOut != nil {
		depth := e.Book.GetDepth(50)
		msg := map[string]interface{}{
			"type":   "orderbook",
			"symbol": "BTC_USDT",
			"data":   depth,
		}
		if b, err := json.Marshal(msg); err == nil {
			// Неблокирующая отправка
			select {
			case e.wsOut <- b:
			default:
			}
		}
	}

	return trades, nil
}

// CancelOrder отменяет активный лимитный ордер, убирает его из стакана и размораживает остаток средств.
// TODO: Оптимизировать блокировки, возможно, разделив лок стакана и лок аккаунта.
func (e *Engine) CancelOrder(orderID string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	order, err := e.Store.GetOrder(orderID)
	if err != nil {
		return err
	}

	// Если ордер уже закрыт, ничего не делаем
	if order.Status == domain.StatusFilled || order.Status == domain.StatusCanceled {
		return nil
	}

	acc, err := e.Store.GetAccount(order.AccountID)
	if err != nil {
		return err
	}

	// Удаляем из стакана
	if e.Book.CancelOrder(orderID) {
		order.Status = domain.StatusCanceled
		_ = e.Store.UpdateOrder(order)
		
		// Возвращаем (размораживаем) незаполненный остаток
		if order.Side == domain.Buy {
			unfilledQuote := order.Price.Mul(order.UnfilledQty())
			acc.UnlockFunds(e.Pair.QuoteAsset, unfilledQuote)
		} else {
			acc.UnlockFunds(e.Pair.BaseAsset, order.UnfilledQty())
		}
		
		// Рассылаем новый стакан по WebSocket
		if e.wsOut != nil {
			depth := e.Book.GetDepth(50)
			msg := map[string]interface{}{
				"type":   "orderbook",
				"symbol": "BTC_USDT",
				"data":   depth,
			}
			if b, err := json.Marshal(msg); err == nil {
				select {
				case e.wsOut <- b:
				default:
				}
			}
		}

		return nil
	}
	
	return storage.ErrOrderNotFound
}

// GetDepth возвращает снимок стакана (срез по лучшим ценам).
func (e *Engine) GetDepth(limit int) orderbook.Depth {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.Book.GetDepth(limit)
}

// GetActiveOrdersCount возвращает количество активных ордеров в стакане.
func (e *Engine) GetActiveOrdersCount() int {
	e.mu.Lock()
	defer e.mu.Unlock()
	return len(e.Book.OrdersMap)
}

// triggerConditionals проверяет условные ордера и если цена достигла триггера, отправляет их в стакан.
func (e *Engine) triggerConditionals() {
	if len(e.conditionals) == 0 || e.lastPrice.IsZero() {
		return
	}

	var remaining []*domain.Order
	for _, condOrder := range e.conditionals {
		triggered := false
		
		switch condOrder.Type {
		case domain.StopLimit:
			// StopLimit срабатывает, когда цена падает (для продаж) или растет (для покупок) до триггера
			if condOrder.Side == domain.Sell && e.lastPrice.Cmp(condOrder.TriggerPrice) <= 0 {
				triggered = true
			} else if condOrder.Side == domain.Buy && e.lastPrice.Cmp(condOrder.TriggerPrice) >= 0 {
				triggered = true
			}
		case domain.TakeProfit:
			// TakeProfit срабатывает, когда цена растет (для продаж) или падает (для покупок) до триггера
			if condOrder.Side == domain.Sell && e.lastPrice.Cmp(condOrder.TriggerPrice) >= 0 {
				triggered = true
			} else if condOrder.Side == domain.Buy && e.lastPrice.Cmp(condOrder.TriggerPrice) <= 0 {
				triggered = true
			}
		}

		if triggered {
			// Конвертируем в обычный лимитный или маркет ордер
			condOrder.Type = domain.Limit
			if condOrder.Price.IsZero() {
				condOrder.Type = domain.Market
			}
			
			// Процессим его в стакане
			trades := e.Matcher.ProcessLimitOrder(condOrder)
			_ = e.Store.UpdateOrder(condOrder)
			
			for _, t := range trades {
				e.Settler.Queue.Enqueue(t)
			}
		} else {
			remaining = append(remaining, condOrder)
		}
	}
	e.conditionals = remaining
}
