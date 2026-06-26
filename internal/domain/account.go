package domain

import (
	"dex/pkg/decimal"
	"errors"
	"sync"
)

var (
	// ErrInsufficientFunds возвращается, когда у пользователя недостаточно средств для выполнения операции.
	ErrInsufficientFunds = errors.New("недостаточно средств")
)

// Account представляет кошелек (балансы) пользователя на бирже.
// Включает в себя свободные средства и средства, замороженные в активных ордерах.
type Account struct {
	Address string       // Уникальный адрес пользователя (например, публичный ключ или адрес из EVM)
	mu      sync.RWMutex // Мьютекс для потокобезопасного доступа к балансам

	// Balances хранит доступные средства. Ключ — тикер актива (например, "USDT" или "BTC").
	Balances map[string]decimal.Decimal
	
	// Locked хранит замороженные средства (те, что сейчас находятся в ордерах).
	Locked map[string]decimal.Decimal
}

// NewAccount инициализирует новый кошелек пользователя с пустыми балансами.
func NewAccount(address string) *Account {
	return &Account{
		Address:  address,
		Balances: make(map[string]decimal.Decimal),
		Locked:   make(map[string]decimal.Decimal),
	}
}

// Deposit зачисляет доступные средства на аккаунт.
// TODO: Добавить логирование депозитов (Audit Trail) для финансового комплаенса.
func (a *Account) Deposit(asset string, amount decimal.Decimal) {
	a.mu.Lock()
	defer a.mu.Unlock()
	
	if current, exists := a.Balances[asset]; exists {
		a.Balances[asset] = current.Add(amount)
	} else {
		a.Balances[asset] = amount
	}
}

// Withdraw списывает доступные средства с аккаунта.
// Возвращает ErrInsufficientFunds, если баланс меньше суммы списания.
func (a *Account) Withdraw(asset string, amount decimal.Decimal) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	
	current, exists := a.Balances[asset]
	if !exists || current.Cmp(amount) < 0 {
		return ErrInsufficientFunds
	}
	
	a.Balances[asset] = current.Sub(amount)
	return nil
}

// LockFunds перемещает средства из доступного баланса (Balances) в замороженный (Locked).
// Вызывается в момент создания лимитного ордера, чтобы юзер не мог дважды потратить эти деньги.
func (a *Account) LockFunds(asset string, amount decimal.Decimal) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	current, exists := a.Balances[asset]
	if !exists || current.Cmp(amount) < 0 {
		return ErrInsufficientFunds
	}

	a.Balances[asset] = current.Sub(amount)
	
	if locked, exists := a.Locked[asset]; exists {
		a.Locked[asset] = locked.Add(amount)
	} else {
		a.Locked[asset] = amount
	}
	
	return nil
}

// UnlockFunds возвращает средства из заморозки (Locked) обратно в доступные (Balances).
// Вызывается при отмене ордера пользователем.
func (a *Account) UnlockFunds(asset string, amount decimal.Decimal) {
	a.mu.Lock()
	defer a.mu.Unlock()

	locked, exists := a.Locked[asset]
	if !exists || locked.Cmp(amount) < 0 {
		return // Защита от критических ошибок. В нормальной системе сюда не дойдет.
	}

	a.Locked[asset] = locked.Sub(amount)
	a.Balances[asset] = a.Balances[asset].Add(amount)
}

// DeductLocked перманентно списывает замороженные средства.
// Вызывается движком (Settler) после успешного исполнения ордера (совершения сделки).
func (a *Account) DeductLocked(asset string, amount decimal.Decimal) {
	a.mu.Lock()
	defer a.mu.Unlock()

	locked, exists := a.Locked[asset]
	if !exists || locked.Cmp(amount) < 0 {
		return 
	}
	a.Locked[asset] = locked.Sub(amount)
}

// SettleTrade атомарно применяет изменения балансов по сделке к аккаунту.
func (a *Account) SettleTrade(depositAsset string, depositAmount decimal.Decimal, deductAsset string, deductAmount decimal.Decimal) {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Зачисляем
	if current, exists := a.Balances[depositAsset]; exists {
		a.Balances[depositAsset] = current.Add(depositAmount)
	} else {
		a.Balances[depositAsset] = depositAmount
	}

	// Списываем замороженные
	if locked, exists := a.Locked[deductAsset]; exists {
		if locked.Cmp(deductAmount) >= 0 {
			a.Locked[deductAsset] = locked.Sub(deductAmount)
		}
	}
}

// GetBalance возвращает копию текущего доступного баланса по указанному активу.
func (a *Account) GetBalance(asset string) decimal.Decimal {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if b, ok := a.Balances[asset]; ok {
		return b.Copy()
	}
	return decimal.Zero()
}
