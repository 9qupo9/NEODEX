package domain

import (
	"dex/pkg/decimal"
	"time"
)

// Side определяет направление ордера: покупка или продажа.
type Side string

// OrderType определяет тип ордера (например, по рынку или по лимитной цене).
type OrderType string

// OrderStatus отражает текущее состояние ордера в жизненном цикле биржи.
type OrderStatus string

const (
	Buy  Side = "BUY"
	Sell Side = "SELL"

	Limit  OrderType = "LIMIT"
	Market OrderType = "MARKET"
	StopLimit  OrderType = "STOP_LIMIT"
	TakeProfit OrderType = "TAKE_PROFIT"

	// TimeInForce определяет как долго ордер будет активен
	GTC string = "GTC" // Good Til Cancelled
	IOC string = "IOC" // Immediate Or Cancel
	FOK string = "FOK" // Fill Or Kill

	StatusNew             OrderStatus = "NEW"              // Ордер только что создан и добавлен в стакан
	StatusPartiallyFilled OrderStatus = "PARTIALLY_FILLED" // Ордер исполнен частично
	StatusFilled          OrderStatus = "FILLED"           // Ордер исполнен полностью
	StatusCanceled        OrderStatus = "CANCELED"         // Ордер отменен пользователем
	StatusRejected        OrderStatus = "REJECTED"         // Ордер отклонен движком (например, из-за нехватки средств)
)

// Order представляет намерение пользователя купить или продать определенный актив.
type Order struct {
	ID           string          // Уникальный идентификатор ордера (генерируется pkg/id)
	AccountID    string          // ID пользователя, который выставил ордер
	Pair         Pair            // Торговая пара (например, BTC/USDT)
	Side         Side            // BUY или SELL
	Type         OrderType       // LIMIT, MARKET, STOP_LIMIT, TAKE_PROFIT
	TimeInForce  string          // GTC, IOC, FOK
	Price        decimal.Decimal // Цена исполнения. Для MARKET ордеров будет 0
	TriggerPrice decimal.Decimal // Триггерная цена для StopLimit и TakeProfit
	Qty          decimal.Decimal // Первоначально запрошенное количество (объем)
	FilledQty    decimal.Decimal // Количество, которое уже было исполнено
	Status       OrderStatus     // Текущий статус
	CreatedAt    int64           // Время создания в Unix nanoseconds для точной сортировки FIFO в стакане

	// Специфичные поля для Фьючерсов (Derivatives)
	IsFutures    bool       // true, если это фьючерсный ордер
	Leverage     int        // Плечо (от 1 до 100)
	MarginMode   MarginMode // ISOLATED или CROSS
	ReduceOnly   bool       // Если true, ордер может только уменьшать позицию, но не увеличивать/открывать новую
}

// UnfilledQty возвращает объем ордера, который еще ожидает исполнения.
// Формула: Общий объем - Исполненный объем
func (o *Order) UnfilledQty() decimal.Decimal {
	return o.Qty.Sub(o.FilledQty)
}

// IsFilled проверяет, был ли ордер полностью исполнен.
// Возвращает true, если исполненный объем больше или равен запрошенному.
func (o *Order) IsFilled() bool {
	return o.FilledQty.Cmp(o.Qty) >= 0
}

// NewOrder — конструктор для создания нового экземпляра ордера.
func NewOrder(id, accountID string, pair Pair, side Side, oType OrderType, price, qty decimal.Decimal, triggerPrice decimal.Decimal, timeInForce string) *Order {
	if timeInForce == "" {
		timeInForce = GTC
	}
	return &Order{
		ID:           id,
		AccountID:    accountID,
		Pair:         pair,
		Side:         side,
		Type:         oType,
		TimeInForce:  timeInForce,
		Price:        price,
		TriggerPrice: triggerPrice,
		Qty:          qty,
		FilledQty:    decimal.Zero(),
		Status:       StatusNew,
		CreatedAt:    time.Now().UnixNano(),
	}
}
