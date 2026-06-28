package engine

import (
	"dex/internal/domain"
	"dex/internal/storage"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

var ErrEngineNotFound = errors.New("engine not found for this pair")

// EngineRouter manages multiple matching engines, one for each trading pair.
type EngineRouter struct {
	store         storage.Store
	wsOut         chan []byte
	engines       map[string]*Engine
	mu            sync.RWMutex
	lastLatencyNs int64 // Последнее время сведения в наносекундах
	isHalted      int32 // 1 если торги остановлены, 0 если активны
}

// NewEngineRouter creates a new router instance.
func NewEngineRouter(store storage.Store, wsOut chan []byte) *EngineRouter {
	return &EngineRouter{
		store:   store,
		wsOut:   wsOut,
		engines: make(map[string]*Engine),
	}
}

// CreateEngine динамически создает новую торговую пару и запускает ее ядро.
func (r *EngineRouter) CreateEngine(pair domain.Pair) {
	eng := NewEngine(pair, r.store, r.wsOut)
	r.RegisterEngine(eng)
}

// GetAllEngines возвращает список всех активных ядер.
func (r *EngineRouter) GetAllEngines() []*Engine {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var all []*Engine
	for _, e := range r.engines {
		all = append(all, e)
	}
	return all
}

// RegisterEngine adds a new engine to the router.
func (r *EngineRouter) RegisterEngine(e *Engine) {
	r.mu.Lock()
	defer r.mu.Unlock()
	symbol := fmt.Sprintf("%s_%s", e.Pair.BaseAsset, e.Pair.QuoteAsset)
	r.engines[symbol] = e
}

// GetEngine returns the engine associated with the given pair.
func (r *EngineRouter) GetEngine(base, quote string) (*Engine, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	symbol := fmt.Sprintf("%s_%s", base, quote)
	if e, exists := r.engines[symbol]; exists {
		return e, nil
	}
	return nil, ErrEngineNotFound
}

// PlaceOrder находит нужный Engine и отправляет ордер в его стакан.
func (r *EngineRouter) PlaceOrder(order *domain.Order) ([]*domain.Trade, error) {
	if r.IsHalted() {
		return nil, errors.New("trading halted by admin")
	}

	start := time.Now()

	e, err := r.GetEngine(order.Pair.BaseAsset, order.Pair.QuoteAsset)
	if err != nil {
		return nil, err
	}

	trades, err := e.PlaceOrder(order)

	// Сохраняем метрику задержки
	atomic.StoreInt64(&r.lastLatencyNs, time.Since(start).Nanoseconds())

	return trades, err
}

// GetLastLatencyMs возвращает последнюю задержку в миллисекундах.
func (r *EngineRouter) GetLastLatencyMs() float64 {
	ns := atomic.LoadInt64(&r.lastLatencyNs)
	return float64(ns) / 1e6
}

// SetHalted останавливает или возобновляет торги (создание новых ордеров)
func (r *EngineRouter) SetHalted(halt bool) {
	if halt {
		atomic.StoreInt32(&r.isHalted, 1)
	} else {
		atomic.StoreInt32(&r.isHalted, 0)
	}
}

// IsHalted проверяет, остановлены ли торги
func (r *EngineRouter) IsHalted() bool {
	return atomic.LoadInt32(&r.isHalted) == 1
}

// GetTotalActiveOrders возвращает сумму всех активных ордеров по всем стаканам.
func (r *EngineRouter) GetTotalActiveOrders() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	total := 0
	for _, e := range r.engines {
		total += e.GetActiveOrdersCount()
	}
	return total
}

// CancelOrder routes a cancel request.
func (r *EngineRouter) CancelOrder(orderId string, store storage.Store) error {
	order, err := store.GetOrder(orderId)
	if err != nil {
		return err
	}
	e, err := r.GetEngine(order.Pair.BaseAsset, order.Pair.QuoteAsset)
	if err != nil {
		return err
	}
	return e.CancelOrder(orderId)
}

// GetDepth routes depth request to correct engine.
func (r *EngineRouter) GetDepth(base, quote string, limit int) (interface{}, error) {
	e, err := r.GetEngine(base, quote)
	if err != nil {
		return nil, err
	}
	return e.GetDepth(limit), nil
}
