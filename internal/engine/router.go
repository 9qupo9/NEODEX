package engine

import (
	"dex/internal/domain"
	"dex/internal/storage"
	"errors"
	"fmt"
	"sync"
)

var ErrEngineNotFound = errors.New("engine not found for this pair")

// EngineRouter manages multiple matching engines, one for each trading pair.
type EngineRouter struct {
	engines map[string]*Engine
	mu      sync.RWMutex
}

// NewEngineRouter creates a new router instance.
func NewEngineRouter() *EngineRouter {
	return &EngineRouter{
		engines: make(map[string]*Engine),
	}
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

// PlaceOrder routes an order to the correct engine.
func (r *EngineRouter) PlaceOrder(order *domain.Order) ([]*domain.Trade, error) {
	e, err := r.GetEngine(order.Pair.BaseAsset, order.Pair.QuoteAsset)
	if err != nil {
		return nil, err
	}
	return e.PlaceOrder(order)
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
