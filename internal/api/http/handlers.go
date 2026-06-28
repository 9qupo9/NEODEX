package http

import (
	"dex/internal/api/tcp"
	"dex/internal/api/ws"
	"dex/internal/domain"
	"dex/internal/engine"
	"dex/internal/storage"
	"dex/pkg/decimal"
	"dex/pkg/id"
	"dex/pkg/syslog"
	"encoding/json"
	"net/http"
	"strings"
	"sync/atomic"
)

// Handlers объединяет все методы для HTTP REST API.
type Handlers struct {
	Router        *engine.EngineRouter  // Маршрутизатор спотовых движков
	FuturesEngine *engine.FuturesEngine // Ссылка на фьючерсное ядро (пока одно для тестов)
	Store         storage.Store
	TCPHandler    *tcp.Handler   // Ссылка на TCP обработчик для метрик
	WSHub         *ws.ShardedHub // Ссылка на WS хаб для метрик
}

// NewHandlers создает экземпляр обработчиков.
func NewHandlers(r *engine.EngineRouter, f *engine.FuturesEngine, store storage.Store, th *tcp.Handler, wh *ws.ShardedHub) *Handlers {
	return &Handlers{Router: r, FuturesEngine: f, Store: store, TCPHandler: th, WSHub: wh}
}

// HandleAdminMetrics возвращает JSON с системными метриками для админки.
func (h *Handlers) HandleAdminMetrics(w http.ResponseWriter, r *http.Request) {
	users, vol, rev := h.Store.GetSystemMetrics()
	tcpBots := atomic.LoadInt32(&h.TCPHandler.ActiveConnections)
	latency := h.Router.GetLastLatencyMs()
	activeOrders := h.Router.GetTotalActiveOrders()
	wsClients := h.WSHub.GetTotalClients()
	logs := syslog.GlobalBuffer.GetLogs()
	isHalted := h.Router.IsHalted()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"users":         users,
		"volume":        vol,
		"revenue":       rev,
		"tcp_bots":      tcpBots,
		"ws_clients":    wsClients,
		"active_orders": activeOrders,
		"latency_ms":    latency,
		"is_halted":     isHalted,
		"logs":          logs,
	})
}

// HandleAdminAction обрабатывает команды из админ-панели (снапшот, очистка кэша, остановка торгов).
func (h *Handlers) HandleAdminAction(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Action string `json:"action"`
		Halt   bool   `json:"halt"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	switch req.Action {
	case "snapshot":
		if err := h.Store.Snapshot(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "clear_cache":
		syslog.GlobalBuffer.Clear()
	case "toggle_trading":
		h.Router.SetHalted(req.Halt)
	default:
		http.Error(w, "unknown action", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// OrderRequest — структура входящего JSON для создания ордера.
type OrderRequest struct {
	AccountID    string `json:"accountId"`    // Кто выставляет
	Base         string `json:"base"`         // Базовый актив
	Quote        string `json:"quote"`        // Котируемый актив
	Side         string `json:"side"`         // BUY или SELL
	Type         string `json:"type"`         // LIMIT, MARKET, STOP_LIMIT, TAKE_PROFIT
	TimeInForce  string `json:"timeInForce"`  // GTC, IOC, FOK
	Price        string `json:"price"`        // Цена за 1 лот
	TriggerPrice string `json:"triggerPrice"` // Триггер для стоп-ордеров
	Qty          string `json:"qty"`          // Общее количество
	Signature    string `json:"signature"`    // Web3 подпись (EIP-712)
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
	triggerPrice, _ := decimal.NewFromString(req.TriggerPrice)

	pair := domain.Pair{BaseAsset: req.Base, QuoteAsset: req.Quote}
	order := domain.NewOrder(id.New(), req.AccountID, pair, domain.Side(req.Side), domain.OrderType(req.Type), price, qty, triggerPrice, req.TimeInForce)

	// Передаем ордер в маршрутизатор движков (EngineRouter)
	trades, err := h.Router.PlaceOrder(order)
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

// HandleGetOrderbook (GET /api/v1/orderbook?symbol=...) возвращает срез текущего стакана (глубину рынка).
func (h *Handlers) HandleGetOrderbook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	symbol := r.URL.Query().Get("symbol")
	if symbol == "" {
		symbol = "BTC_USDT" // Default for backward compatibility
	}

	parts := strings.Split(symbol, "_")
	if len(parts) != 2 {
		http.Error(w, "invalid symbol format, expected BASE_QUOTE", http.StatusBadRequest)
		return
	}
	base := parts[0]
	quote := parts[1]

	// Получаем топ 50 уровней
	depth, err := h.Router.GetDepth(base, quote, 50)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

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

	orders, err := h.Store.GetOrdersByAccount(accountId)
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

	acc, err := h.Store.GetAccount(accountId)
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

	if err := h.Router.CancelOrder(orderId, h.Store); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "canceled"})
}

// HandleAdminUsers возвращает список всех пользователей системы.
func (h *Handlers) HandleAdminUsers(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.Store.GetAllAccounts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accounts)
}

// HandleAdminBlockUser блокирует или разблокирует пользователя.
func (h *Handlers) HandleAdminBlockUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Address string `json:"address"`
		Block   bool   `json:"block"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	acc, err := h.Store.GetAccount(req.Address)
	if err != nil {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}
	acc.SetBlocked(req.Block)
	w.WriteHeader(http.StatusOK)
}
