package storage

import (
	"dex/internal/domain"
	"errors"
)

var (
	ErrAccountNotFound = errors.New("аккаунт не найден")
	ErrOrderNotFound   = errors.New("ордер не найден")
)

// Store определяет интерфейс для всех видов хранилищ (in-memory, file, database).
// Использование интерфейса позволяет легко подменять реализацию без изменения бизнес-логики.
// TODO: Расширить интерфейс методами выборки списка ордеров (GetOpenOrdersByAccount).
type Store interface {
	// Accounts (Аккаунты)
	
	// GetAccount возвращает аккаунт по его адресу. Вернет ошибку, если аккаунт не существует.
	GetAccount(address string) (*domain.Account, error)
	
	// GetAllAccounts возвращает список всех аккаунтов.
	GetAllAccounts() ([]*domain.Account, error)
	
	// CreateAccount создает новый пустой аккаунт в хранилище.
	CreateAccount(address string) (*domain.Account, error)
	
	// SaveAccount сохраняет изменения аккаунта (например, после изменения балансов).
	SaveAccount(account *domain.Account) error

	// Orders (Ордера)
	
	// SaveOrder сохраняет новый ордер.
	SaveOrder(order *domain.Order) error
	
	// GetOrder возвращает существующий ордер по ID.
	GetOrder(id string) (*domain.Order, error)
	
	// GetOrdersByAccount возвращает список всех ордеров пользователя.
	GetOrdersByAccount(address string) ([]*domain.Order, error)
	
	// UpdateOrder обновляет стейт ордера (заполнение объема, изменение статуса).
	UpdateOrder(order *domain.Order) error

	// Trades (Сделки)
	
	// SaveTrade сохраняет завершенную сделку в историю.
	SaveTrade(trade *domain.Trade) error

	// SettleTradeBalances атомарно сохраняет сделку и обновляет балансы
	SettleTradeBalances(trade *domain.Trade) error

	// Positions (Позиции для фьючерсов)
	GetPositionsByAccount(address string) ([]*domain.Position, error)
	GetPosition(id string) (*domain.Position, error)
	SavePosition(position *domain.Position) error

	// Metrics (Админ-панель)
	GetSystemMetrics() (usersCount int, totalVolume string, protocolRevenue string)

	// Snapshot сохраняет дамп БД (если поддерживается реализацией).
	Snapshot() error
}
