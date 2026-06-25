package http

import (
	"encoding/json"
	"net/http"
	"dex/internal/domain"
	"dex/internal/engine"
	"dex/pkg/decimal"
	"dex/pkg/id"
)

// Handlers объединяет все методы для HTTP REST API.
type Handlers struct {
	Engine *engine.Engine // Ссылка на ядро
}

// NewHandlers создает экземпляр обработчиков.
func NewHandlers(e *engine.Engine) *Handlers {
	return &Handlers{Engine: e}
}

// OrderRequest — структура входящего JSON для создания ордера.
type OrderRequest struct {
	AccountID string `json:"accountId"` // Кто выставляет
	Base      string `json:"base"`      // Базовый актив
	Quote     string `json:"quote"`     // Котируемый актив
	Side      string `json:"side"`      // BUY или SELL
	Type      string `json:"type"`      // LIMIT или MARKET
	Price     string `json:"price"`     // Цена за 1 лот
	Qty       string `json:"qty"`       // Общее количество
}

// HandlePlaceOrder (POST /api/v1/order) обрабатывает запрос на создание ордера.
// Валидирует данные, парсит числа и отправляет команду в ядро (Engine).
func (h *Handlers) HandlePlaceOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req OrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	price, _ := decimal.NewFromString(req.Price)
	qty, _ := decimal.NewFromString(req.Qty)

	pair := domain.Pair{BaseAsset: req.Base, QuoteAsset: req.Quote}
	order := domain.NewOrder(id.New(), req.AccountID, pair, domain.Side(req.Side), domain.OrderType(req.Type), price, qty)

	// Передаем ордер в движок (блокирующе, но очень быстро, так как всё in-memory)
	trades, err := h.Engine.PlaceOrder(order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"orderId": order.ID,
		"status":  order.Status,
		"trades":  len(trades), // Количество совершенных сделок прямо во время исполнения
	})
}

// HandleGetOrderbook (GET /api/v1/orderbook) возвращает срез текущего стакана (глубину рынка).
func (h *Handlers) HandleGetOrderbook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// Получаем топ 50 уровней
	depth := h.Engine.GetDepth(50)
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(depth)
}

// HandleGetOrders (GET /api/v1/orders?accountId=...) возвращает все ордера пользователя.
func (h *Handlers) HandleGetOrders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	accountId := r.URL.Query().Get("accountId")
	if accountId == "" {
		http.Error(w, "accountId required", http.StatusBadRequest)
		return
	}

	orders, err := h.Engine.Store.GetOrdersByAccount(accountId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}


// HandleGetBalance (GET /api/v1/balance?accountId=...) возвращает доступные средства пользователя.
func (h *Handlers) HandleGetBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	accountId := r.URL.Query().Get("accountId")
	if accountId == "" {
		http.Error(w, "accountId required", http.StatusBadRequest)
		return
	}

	acc, err := h.Engine.Store.GetAccount(accountId)
	if err != nil {
		http.Error(w, "account not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	
	// Конвертируем балансы из Decimal в строки (чтобы не терять точность в JSON / JS)
	balMap := make(map[string]string)
	for k, v := range acc.Balances {
		balMap[k] = v.String()
	}
	json.NewEncoder(w).Encode(balMap)
}

// HandleCancelOrder (POST /api/v1/order/cancel?orderId=...) отменяет ордер и возвращает средства из заморозки.
func (h *Handlers) HandleCancelOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	orderId := r.URL.Query().Get("orderId")
	if orderId == "" {
		http.Error(w, "orderId required", http.StatusBadRequest)
		return
	}

	if err := h.Engine.CancelOrder(orderId); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "canceled"})
}
