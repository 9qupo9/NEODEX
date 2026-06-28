package http

import (
	"dex/internal/domain"
	"dex/internal/storage"
	"dex/pkg/decimal"
	"dex/pkg/id"
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

// StakingManager (временно хранится тут для упрощения архитектуры)
var (
	stakesMu sync.RWMutex
	stakes   = make(map[string]*domain.Stake)
)

// StakingHandler обрабатывает запросы на стейкинг (DeFi).
// Принцип SOLID: Single Responsibility Principle.
type StakingHandler struct {
	Store storage.Store
}

func NewStakingHandler(store storage.Store) *StakingHandler {
	return &StakingHandler{Store: store}
}

type StakeRequest struct {
	AccountID string `json:"accountId"`
	Asset     string `json:"asset"`
	Amount    string `json:"amount"`
	Days      int    `json:"days"` // Количество дней блокировки
}

// HandleStake обрабатывает POST-запрос на стейкинг средств.
func (h *StakingHandler) HandleStake(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req StakeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	amount, _ := decimal.NewFromString(req.Amount)
	if amount.IsZero() {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}

	acc, err := h.Store.GetAccount(req.AccountID)
	if err != nil {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	// Списываем средства с баланса (помещаем в стейкинг)
	if err := acc.Withdraw(req.Asset, amount); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.Store.SaveAccount(acc)

	// Стандартная доходность: 10% годовых для USDT, 5% для BTC
	apy := decimal.MustParse("0.10")
	if req.Asset == "BTC" || req.Asset == "ETH" {
		apy = decimal.MustParse("0.05")
	}

	lockDuration := time.Duration(req.Days) * 24 * time.Hour
	stake := domain.NewStake(id.New(), req.AccountID, req.Asset, amount, apy, lockDuration)

	stakesMu.Lock()
	stakes[stake.ID] = stake
	stakesMu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stake)
}

// HandleGetStakes возвращает активные стейки пользователя.
func (h *StakingHandler) HandleGetStakes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	accountID := r.URL.Query().Get("accountId")
	if accountID == "" {
		http.Error(w, "accountId required", http.StatusBadRequest)
		return
	}

	stakesMu.RLock()
	var userStakes []*domain.Stake
	for _, s := range stakes {
		if s.AccountID == accountID && s.IsActive {
			userStakes = append(userStakes, s)
		}
	}
	stakesMu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	
	// Конвертируем в DTO с расчетом текущей награды
	var response []map[string]interface{}
	for _, s := range userStakes {
		response = append(response, map[string]interface{}{
			"id":        s.ID,
			"asset":     s.Asset,
			"amount":    s.Amount.String(),
			"apy":       s.APY.String(),
			"reward":    s.CalculateReward().String(),
			"unlocksAt": s.UnlocksAt,
		})
	}
	if response == nil {
		response = []map[string]interface{}{}
	}
	
	json.NewEncoder(w).Encode(response)
}

// HandleUnstake забирает стейк и награду, если прошел период лока.
func (h *StakingHandler) HandleUnstake(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	stakeID := r.URL.Query().Get("id")
	if stakeID == "" {
		http.Error(w, "id required", http.StatusBadRequest)
		return
	}

	stakesMu.Lock()
	stake, ok := stakes[stakeID]
	if !ok || !stake.IsActive {
		stakesMu.Unlock()
		http.Error(w, "Stake not found or inactive", http.StatusNotFound)
		return
	}

	if time.Now().Before(stake.UnlocksAt) {
		stakesMu.Unlock()
		http.Error(w, "Stake is still locked", http.StatusBadRequest)
		return
	}

	stake.IsActive = false
	stakesMu.Unlock()

	acc, _ := h.Store.GetAccount(stake.AccountID)
	reward := stake.CalculateReward()
	
	// Возвращаем тело депозита + награду
	acc.Deposit(stake.Asset, stake.Amount)
	acc.Deposit(stake.Asset, reward)
	h.Store.SaveAccount(acc)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
		"reward": reward.String(),
	})
}
