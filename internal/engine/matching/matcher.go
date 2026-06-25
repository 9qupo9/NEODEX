package matching

import (
	"dex/internal/domain"
	"dex/internal/engine/orderbook"
	"dex/pkg/decimal"
	"dex/pkg/id"
)

// Matcher — это основной движок сведения заявок.
// Он пересекает входящие рыночные/лимитные заявки с ликвидностью в стакане (OrderBook).
type Matcher struct {
	Book *orderbook.Book
}

// NewMatcher инициализирует движок сведения для конкретного стакана.
func NewMatcher(book *orderbook.Book) *Matcher {
	return &Matcher{Book: book}
}

// ProcessLimitOrder принимает входящий лимитный ордер и пытается немедленно свести его
// со встречными заявками в стакане. Оставшийся неисполненным объем (остаток) добавляется в стакан.
// Возвращает массив совершенных сделок (Trades), которые передаются дальше в Клиринг (Settlement).
// TODO: Выделить логику исполнения рыночных (Market) ордеров в отдельный метод.
func (m *Matcher) ProcessLimitOrder(order *domain.Order) []*domain.Trade {
	var trades []*domain.Trade

	// Пока ордер не заполнен, ищем встречные предложения в стакане
	for order.UnfilledQty().Cmp(decimal.Zero()) > 0 {
		var bestLevel *orderbook.Level
		var matchPrice decimal.Decimal

		// Ищем лучшую цену встречного ордера (для покупателя — самая дешевая продажа Asks[0])
		if order.Side == domain.Buy {
			if len(m.Book.Asks) == 0 {
				break // Продавцов нет
			}
			bestLevel = m.Book.Asks[0]
			matchPrice = bestLevel.Price
			if matchPrice.Cmp(order.Price) > 0 {
				break // Цена лучшего продавца выше нашей лимитной цены, сведение невозможно
			}
		} else {
			if len(m.Book.Bids) == 0 {
				break // Покупателей нет
			}
			bestLevel = m.Book.Bids[0]
			matchPrice = bestLevel.Price
			if matchPrice.Cmp(order.Price) < 0 {
				break // Цена лучшего покупателя ниже нашей лимитной, сведение невозможно
			}
		}

		// Берем самый первый ордер в очереди на этом ценовом уровне (по принципу FIFO)
		headNode := bestLevel.Orders.Head
		if headNode == nil {
			break
		}
		makerOrder := headNode.Order

		// Высчитываем объем сделки (поглощаем минимум из доступного у мейкера и требуемого у тейкера)
		tradeQty := order.UnfilledQty()
		if makerOrder.UnfilledQty().Cmp(tradeQty) < 0 {
			tradeQty = makerOrder.UnfilledQty()
		}

		// Генерируем объект сделки (Trade)
		trade := domain.NewTrade(
			id.New(),
			makerOrder.ID,
			order.ID,
			makerOrder.AccountID,
			order.AccountID,
			m.Book.Pair,
			matchPrice,
			tradeQty,
			order.Side,
		)
		trades = append(trades, trade)

		// Обновляем состояния исполненного объема у обоих ордеров
		order.FilledQty = order.FilledQty.Add(tradeQty)
		makerOrder.FilledQty = makerOrder.FilledQty.Add(tradeQty)

		// Если Мейкер исполнен полностью, убираем его из стакана и словаря
		if makerOrder.IsFilled() {
			makerOrder.Status = domain.StatusFilled
			bestLevel.Remove(headNode)
			delete(m.Book.OrdersMap, makerOrder.ID)
		} else {
			// Иначе он частично исполнен, нужно обновить только объем уровня
			makerOrder.Status = domain.StatusPartiallyFilled
			bestLevel.Volume = bestLevel.Volume.Sub(tradeQty)
		}

		// Если ценовой уровень полностью опустел (выкупили всю ликвидность), удаляем уровень из стакана
		if bestLevel.Orders.Size == 0 {
			if order.Side == domain.Buy {
				m.Book.Asks = m.Book.Asks[1:]
			} else {
				m.Book.Bids = m.Book.Bids[1:]
			}
		}
	}

	// Обрабатываем остаток входящего (Тейкер) ордера
	if order.IsFilled() {
		order.Status = domain.StatusFilled
	} else if order.FilledQty.Cmp(decimal.Zero()) > 0 {
		order.Status = domain.StatusPartiallyFilled
		// Остаток ставим в стакан как лимитную заявку
		m.Book.AddOrder(order)
	} else {
		// Ничего не сматчилось, ордер целиком идет в стакан
		m.Book.AddOrder(order)
	}

	return trades
}
