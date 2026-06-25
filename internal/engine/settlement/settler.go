package settlement

import (
	"dex/internal/domain"
	"dex/internal/storage"
	"log"
)

// Settler (Клиринг) берет сделки из очереди и производит фактический расчет балансов.
// Он списывает деньги у покупателя и зачисляет продавцу, а также наоборот.
type Settler struct {
	Store *storage.MemoryStore // Использование MemoryStore напрямую для скорости (или интерфейс Store)
	Queue *Queue               // Очередь, откуда приходят сделки
}

// NewSettler создает клиринговый процессор.
func NewSettler(store *storage.MemoryStore, queue *Queue) *Settler {
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
	makerAcc, err := s.Store.GetAccount(trade.MakerAddress)
	if err != nil {
		log.Printf("Критическая ошибка клиринга: Мейкер не найден %s", trade.MakerAddress)
		return
	}
	takerAcc, err := s.Store.GetAccount(trade.TakerAddress)
	if err != nil {
		log.Printf("Критическая ошибка клиринга: Тейкер не найден %s", trade.TakerAddress)
		return
	}

	totalQuote := trade.Price.Mul(trade.Qty)

	if trade.TakerSide == domain.Buy {
		// Тейкер получает Base на доступный баланс
		takerAcc.Deposit(trade.Pair.BaseAsset, trade.Qty)
		// У тейкера списываются замороженные Quote
		takerAcc.DeductLocked(trade.Pair.QuoteAsset, totalQuote)

		// Мейкер получает Quote на доступный баланс
		makerAcc.Deposit(trade.Pair.QuoteAsset, totalQuote)
		// У мейкера списываются замороженные Base
		makerAcc.DeductLocked(trade.Pair.BaseAsset, trade.Qty)
	} else {
		// Тейкер получает Quote
		takerAcc.Deposit(trade.Pair.QuoteAsset, totalQuote)
		// Тейкер отдает замороженные Base
		takerAcc.DeductLocked(trade.Pair.BaseAsset, trade.Qty)

		// Мейкер получает Base
		makerAcc.Deposit(trade.Pair.BaseAsset, trade.Qty)
		// Мейкер отдает замороженные Quote
		makerAcc.DeductLocked(trade.Pair.QuoteAsset, totalQuote)
	}

	// Сохраняем сделку в истории хранилища
	_ = s.Store.SaveTrade(trade)
}
