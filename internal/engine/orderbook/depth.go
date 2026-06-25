package orderbook

// Depth представляет собой снимок текущего состояния стакана.
// Отправляется на клиентскую часть (Web/Mobile) через REST API или WebSocket (глубина рынка).
type Depth struct {
	Bids []PriceLevel `json:"bids"` // Покупки
	Asks []PriceLevel `json:"asks"` // Продажи
}

// PriceLevel — это упрощенное представление уровня (без указателей на очереди) для JSON ответа.
type PriceLevel struct {
	Price  string `json:"price"`  // Цена (форматируется как string для исключения потери точности в JS клиентах)
	Volume string `json:"volume"` // Доступный объем на этой цене
}

// GetDepth собирает верхние N уровней стакана для отправки клиентам.
// limit задает глубину выборки (например, 50 или 100 уровней).
func (b *Book) GetDepth(limit int) Depth {
	depth := Depth{
		Bids: make([]PriceLevel, 0, limit),
		Asks: make([]PriceLevel, 0, limit),
	}

	// Копируем верхние лимиты по Bids
	for i := 0; i < len(b.Bids) && i < limit; i++ {
		depth.Bids = append(depth.Bids, PriceLevel{
			Price:  b.Bids[i].Price.String(),
			Volume: b.Bids[i].Volume.String(),
		})
	}

	// Копируем нижние лимиты по Asks (самые лучшие цены)
	for i := 0; i < len(b.Asks) && i < limit; i++ {
		depth.Asks = append(depth.Asks, PriceLevel{
			Price:  b.Asks[i].Price.String(),
			Volume: b.Asks[i].Volume.String(),
		})
	}

	return depth
}
