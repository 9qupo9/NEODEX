package orderbook

import (
	"dex/pkg/decimal"
)

// Level представляет собой один ценовой уровень в стакане ордеров.
// На одной цене могут стоять сотни пользователей, поэтому уровень агрегирует их в очередь (Orders).
type Level struct {
	Price  decimal.Decimal // Цена уровня (например, $60,000.00 за BTC)
	Orders *OrderQueue     // Очередь заявок пользователей на этой цене (FIFO)
	Volume decimal.Decimal // Суммарный объем всех неисполненных ордеров на этом уровне (для передачи по WebSocket)
}

// NewLevel создает пустой ценовой уровень.
func NewLevel(price decimal.Decimal) *Level {
	return &Level{
		Price:  price,
		Orders: NewOrderQueue(),
		Volume: decimal.Zero(),
	}
}

// Append добавляет ордер в очередь на этом уровне и увеличивает общий объем (Volume).
func (l *Level) Append(node *OrderNode) {
	l.Orders.Push(node.Order)
	l.Volume = l.Volume.Add(node.Order.UnfilledQty())
}

// Remove удаляет отмененный ордер с уровня и уменьшает общий объем (Volume).
func (l *Level) Remove(node *OrderNode) {
	l.Orders.Remove(node)
	l.Volume = l.Volume.Sub(node.Order.UnfilledQty())
}
