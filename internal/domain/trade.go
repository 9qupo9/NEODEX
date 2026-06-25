package domain

import (
	"dex/pkg/decimal"
	"time"
)

// Trade представляет успешно совершенную сделку (факт сведения двух ордеров).
// Это immutable (неизменяемая) сущность, которая записывается в историю навсегда.
// TODO: Добавить структуру для комиссий (FeeAmount, FeeAsset).
type Trade struct {
	ID            string          // Уникальный идентификатор сделки
	MakerOrderID  string          // ID ордера мейкера (того, чей ордер уже стоял в стакане)
	TakerOrderID  string          // ID ордера тейкера (того, кто инициировал сделку "по рынку" или ударил в лимит)
	MakerAddress  string          // Адрес пользователя-мейкера
	TakerAddress  string          // Адрес пользователя-тейкера
	Pair          Pair            // Торговая пара, по которой прошла сделка
	Price         decimal.Decimal // Цена исполнения (всегда равна цене мейкер-ордера)
	Qty           decimal.Decimal // Объем сделки (количество базового актива)
	TakerSide     Side            // Направление тейкера (Buy - тейкер купил, Sell - тейкер продал)
	Timestamp     int64           // Время сделки в наносекундах для агрегации свечей (Kline/OHLCV)
}

// NewTrade — конструктор для создания новой сделки.
// Вызывается внутри движка (Matcher) в момент пересечения Bids и Asks.
func NewTrade(id, makerID, takerID, makerAddr, takerAddr string, pair Pair, price, qty decimal.Decimal, takerSide Side) *Trade {
	return &Trade{
		ID:           id,
		MakerOrderID: makerID,
		TakerOrderID: takerID,
		MakerAddress: makerAddr,
		TakerAddress: takerAddr,
		Pair:         pair,
		Price:        price,
		Qty:          qty,
		TakerSide:    takerSide,
		Timestamp:    time.Now().UnixNano(),
	}
}
