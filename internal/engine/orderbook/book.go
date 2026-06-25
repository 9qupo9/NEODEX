package orderbook

import (
	"dex/internal/domain"
	"dex/pkg/decimal"
	"sort"
)

// Book представляет собой in-memory книгу лимитных ордеров (стакан) для одной торговой пары.
// Оптимизирована для сверхбыстрого добавления и O(1) отмены заявок.
type Book struct {
	Pair domain.Pair

	// Списки отсортированных ценовых уровней
	// Используем slices и бинарный поиск, так как в Go они быстрее деревьев из-за локальности кэша процессора (CPU Cache Locality).
	Bids []*Level // Покупки. Сортировка по УБЫВАНИЮ (самая высокая цена покупки сверху стакана)
	Asks []*Level // Продажи. Сортировка по ВОЗРАСТАНИЮ (самая низкая цена продажи снизу стакана)

	// OrdersMap хранит прямые ссылки на узлы ордеров по их ID.
	// Это позволяет отменять ордера за константное время O(1), не перебирая весь стакан.
	OrdersMap map[string]*OrderNode
}

// NewBook инициализирует пустую книгу для пары.
func NewBook(pair domain.Pair) *Book {
	return &Book{
		Pair:      pair,
		Bids:      make([]*Level, 0),
		Asks:      make([]*Level, 0),
		OrdersMap: make(map[string]*OrderNode),
	}
}

// AddOrder добавляет лимитную заявку в книгу как Мейкер (Maker).
// Функция определяет направление (Bid или Ask) и вызывает соответствующий метод.
func (b *Book) AddOrder(order *domain.Order) {
	node := &OrderNode{Order: order}
	b.OrdersMap[order.ID] = node

	if order.Side == domain.Buy {
		b.addBid(node)
	} else {
		b.addAsk(node)
	}
}

// CancelOrder находит заявку по ID в OrdersMap и мгновенно её удаляет.
// Если на ценовом уровне больше не осталось ордеров, сам уровень тоже удаляется из среза.
// Возвращает true, если отмена прошла успешно.
func (b *Book) CancelOrder(orderID string) bool {
	node, exists := b.OrdersMap[orderID]
	if !exists {
		return false
	}

	price := node.Order.Price
	if node.Order.Side == domain.Buy {
		idx := b.findBidLevel(price)
		if idx != -1 {
			b.Bids[idx].Remove(node)
			if b.Bids[idx].Orders.Size == 0 {
				b.Bids = append(b.Bids[:idx], b.Bids[idx+1:]...) // Удаляем пустой уровень
			}
		}
	} else {
		idx := b.findAskLevel(price)
		if idx != -1 {
			b.Asks[idx].Remove(node)
			if b.Asks[idx].Orders.Size == 0 {
				b.Asks = append(b.Asks[:idx], b.Asks[idx+1:]...)
			}
		}
	}

	delete(b.OrdersMap, orderID)
	return true
}

// addBid находит нужный уровень покупки с помощью бинарного поиска (O(log N)) и вставляет ордер.
// Если уровня нет, он создается в правильном (убывающем) порядке.
func (b *Book) addBid(node *OrderNode) {
	price := node.Order.Price
	idx := sort.Search(len(b.Bids), func(i int) bool {
		return b.Bids[i].Price.Cmp(price) <= 0
	})

	if idx < len(b.Bids) && b.Bids[idx].Price.Cmp(price) == 0 {
		b.Bids[idx].Append(node)
	} else {
		level := NewLevel(price)
		level.Append(node)
		// Вставляем с сохранением сортировки
		b.Bids = append(b.Bids, nil)
		copy(b.Bids[idx+1:], b.Bids[idx:])
		b.Bids[idx] = level
	}
}

// addAsk работает аналогично addBid, но сохраняет возрастающий порядок.
func (b *Book) addAsk(node *OrderNode) {
	price := node.Order.Price
	idx := sort.Search(len(b.Asks), func(i int) bool {
		return b.Asks[i].Price.Cmp(price) >= 0
	})

	if idx < len(b.Asks) && b.Asks[idx].Price.Cmp(price) == 0 {
		b.Asks[idx].Append(node)
	} else {
		level := NewLevel(price)
		level.Append(node)
		// Вставляем с сохранением сортировки
		b.Asks = append(b.Asks, nil)
		copy(b.Asks[idx+1:], b.Asks[idx:])
		b.Asks[idx] = level
	}
}

// findBidLevel ищет индекс уровня покупки. Возвращает -1, если уровень не найден.
func (b *Book) findBidLevel(price decimal.Decimal) int {
	idx := sort.Search(len(b.Bids), func(i int) bool {
		return b.Bids[i].Price.Cmp(price) <= 0
	})
	if idx < len(b.Bids) && b.Bids[idx].Price.Cmp(price) == 0 {
		return idx
	}
	return -1
}

// findAskLevel ищет индекс уровня продажи. Возвращает -1, если уровень не найден.
func (b *Book) findAskLevel(price decimal.Decimal) int {
	idx := sort.Search(len(b.Asks), func(i int) bool {
		return b.Asks[i].Price.Cmp(price) >= 0
	})
	if idx < len(b.Asks) && b.Asks[idx].Price.Cmp(price) == 0 {
		return idx
	}
	return -1
}
