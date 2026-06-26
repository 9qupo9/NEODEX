package http

import (
	"dex/internal/domain"
	"dex/pkg/decimal"
	"dex/pkg/id"
	"encoding/json"
	"net/http"
)

// FuturesOrderRequest — структура входящего JSON для создания фьючерсного ордера.
type FuturesOrderRequest struct {
	AccountID  string `json:"accountId"` // Кто выставляет
	Base       string `json:"base"`      // Базовый актив
	Quote      string `json:"quote"`     // Котируемый актив
	Side       string `json:"side"`      // BUY или SELL (Long/Short)
	Type       string `json:"type"`      // LIMIT или MARKET
	Price      string `json:"price"`     // Цена
	Qty        string `json:"qty"`       // Общее количество
	Leverage   int    `json:"leverage"`  // Плечо
	MarginMode string `json:"marginMode"`// ISOLATED или CROSS
}

// HandlePlaceFuturesOrder (POST /api/v1/futures/order)
func (h *Handlers) HandlePlaceFuturesOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req FuturesOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	price, _ := decimal.NewFromString(req.Price)
	qty, _ := decimal.NewFromString(req.Qty)

	pair := domain.Pair{BaseAsset: req.Base, QuoteAsset: req.Quote}
	order := domain.NewOrder(id.New(), req.AccountID, pair, domain.Side(req.Side), domain.OrderType(req.Type), price, qty)
	
	order.IsFutures = true
	order.Leverage = req.Leverage
	order.MarginMode = domain.MarginMode(req.MarginMode)

	// Передаем ордер во фьючерсный движок
	trades, err := h.FuturesEngine.PlaceOrder(order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"orderId": order.ID,
		"status":  order.Status,
		"trades":  len(trades),
	})
}

// HandleGetFuturesPositions (GET /api/v1/futures/positions?accountId=...)
func (h *Handlers) HandleGetFuturesPositions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	accountId := r.URL.Query().Get("accountId")
	if accountId == "" {
		http.Error(w, "accountId required", http.StatusBadRequest)
		return
	}

	positions, err := h.FuturesEngine.Store.GetPositionsByAccount(accountId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(positions)
}
