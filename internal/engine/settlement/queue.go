package settlement

import "dex/internal/domain"

// Queue представляет собой буферизированный канал (очередь) для асинхронного клиринга.
// Зачем это нужно: Движок сведения (Matcher) не должен заниматься работой с БД и пересчетом балансов,
// иначе это убьет производительность. Matcher моментально сводит ордера и кидает трейды в эту очередь.
// А Settler читает из нее и в фоне обновляет аккаунты пользователей.
type Queue struct {
	Trades chan *domain.Trade
}

// NewQueue инициализирует очередь с заданным буфером (размер буфера важен при спайках на HFT).
func NewQueue(bufferSize int) *Queue {
	return &Queue{
		Trades: make(chan *domain.Trade, bufferSize),
	}
}

// Enqueue добавляет новую сделку в очередь на расчеты.
func (q *Queue) Enqueue(trade *domain.Trade) {
	q.Trades <- trade
}
