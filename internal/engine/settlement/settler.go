package settlement

import (
	"dex/internal/domain"
	"dex/internal/storage"
	"log"
)

// Settler (Клиринг) берет сделки из очереди и производит фактический расчет балансов.
// Он списывает деньги у покупателя и зачисляет продавцу, а также наоборот.
type Settler struct {
	Store storage.Store        // Использование интерфейса Store
	Queue *Queue               // Очередь, откуда приходят сделки
}

// NewSettler создает клиринговый процессор.
func NewSettler(store storage.Store, queue *Queue) *Settler {
	return &Settler{
		Store: store,
		Queue: queue,
	}
}

// Start запускает бесконечный цикл обработки очереди сделок в отдельной горутине.
// TODO: Для масштабирования запустить пул горутин (worker pool) и шардировать их по ID аккаунтов, 
// чтобы избежать deadlocks и ожидания мьютексов.
func (s *Settler) Start() {
	go func() {
		for trade := range s.Queue.Trades {
			s.settleTrade(trade)
		}
	}()
}

// settleTrade — сердце клиринга. Пересчитывает балансы Мейкера и Тейкера.
// На этом этапе средства уже заморожены (Locked), поэтому мы только размораживаем/списываем
// и добавляем на основной баланс (Balances) встречному лицу.
func (s *Settler) settleTrade(trade *domain.Trade) {
	// SettleTradeBalances берет на себя атомарность и сохранение Trade
	if err := s.Store.SettleTradeBalances(trade); err != nil {
		log.Printf("Критическая ошибка клиринга: %v (Trade ID: %s)", err, trade.ID)
	}
}
