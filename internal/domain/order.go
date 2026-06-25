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
	// TODO: Добавить типы StopLimit, TakeProfit, TrailingStop (как на MEXC).

	StatusNew             OrderStatus = "NEW"              // Ордер только что создан и добавлен в стакан
	StatusPartiallyFilled OrderStatus = "PARTIALLY_FILLED" // Ордер исполнен частично
	StatusFilled          OrderStatus = "FILLED"           // Ордер исполнен полностью
	StatusCanceled        OrderStatus = "CANCELED"         // Ордер отменен пользователем
	StatusRejected        OrderStatus = "REJECTED"         // Ордер отклонен движком (например, из-за нехватки средств)
)

// Order представляет намерение пользователя купить или продать определенный актив.
// TODO: Добавить TimeInForce (GTC - Good-Til-Cancelled, IOC - Immediate-Or-Cancel, FOK - Fill-Or-Kill).
type Order struct {
	ID           string          // Уникальный идентификатор ордера (генерируется pkg/id)
	AccountID    string          // ID пользователя, который выставил ордер
	Pair         Pair            // Торговая пара (например, BTC/USDT)
	Side         Side            // BUY или SELL
	Type         OrderType       // LIMIT или MARKET
	Price        decimal.Decimal // Цена исполнения. Для MARKET ордеров будет 0
	Qty          decimal.Decimal // Первоначально запрошенное количество (объем)
	FilledQty    decimal.Decimal // Количество, которое уже было исполнено
	Status       OrderStatus     // Текущий статус
	CreatedAt    int64           // Время создания в Unix nanoseconds для точной сортировки FIFO в стакане
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
// Инициализирует ордер со статусом NEW, нулевым заполненным объемом и текущим временем.
func NewOrder(id, accountID string, pair Pair, side Side, oType OrderType, price, qty decimal.Decimal) *Order {
	return &Order{
		ID:        id,
		AccountID: accountID,
		Pair:      pair,
		Side:      side,
		Type:      oType,
		Price:     price,
		Qty:       qty,
		FilledQty: decimal.Zero(),
		Status:    StatusNew,
		CreatedAt: time.Now().UnixNano(),
	}
}
