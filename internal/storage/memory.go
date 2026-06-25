package storage

import (
	"dex/internal/domain"
	"sync"
)

// MemoryStore — реализация хранилища Store полностью в оперативной памяти (in-memory).
// Обеспечивает максимальную скорость работы, необходимую для HFT (High-Frequency Trading) движков.
// Идеально для биржи уровня MEXC, где задержки ввода-вывода (I/O) к базе данных недопустимы в основном цикле.
type MemoryStore struct {
	mu       sync.RWMutex                  // Мьютекс для защиты от состояний гонки (Race conditions) при конкурентных картах
	accounts map[string]*domain.Account    // Карта пользователей
	orders   map[string]*domain.Order      // Карта ордеров
	trades   map[string]*domain.Trade      // Карта трейдов
}

// NewMemoryStore инициализирует пустые карты для in-memory хранилища.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		accounts: make(map[string]*domain.Account),
		orders:   make(map[string]*domain.Order),
		trades:   make(map[string]*domain.Trade),
	}
}

// GetAccount безопасно (с R-Lock) читает данные аккаунта из карты.
func (s *MemoryStore) GetAccount(address string) (*domain.Account, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	if acc, ok := s.accounts[address]; ok {
		return acc, nil
	}
	return nil, ErrAccountNotFound
}

// CreateAccount безопасно создает аккаунт, если он еще не существует.
func (s *MemoryStore) CreateAccount(address string) (*domain.Account, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if acc, ok := s.accounts[address]; ok {
		return acc, nil // Уже существует, возвращаем старый
	}
	
	acc := domain.NewAccount(address)
	s.accounts[address] = acc
	return acc, nil
}

// SaveAccount записывает измененный аккаунт. В memory-реализации указатель уже лежит в мапе, 
// но метод нужен для соответствия интерфейсу и триггера AOF лога в FileStore.
func (s *MemoryStore) SaveAccount(account *domain.Account) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.accounts[account.Address] = account
	return nil
}

// SaveOrder добавляет новый ордер в in-memory карту.
func (s *MemoryStore) SaveOrder(order *domain.Order) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.orders[order.ID] = order
	return nil
}

// GetOrder возвращает ордер по ID.
func (s *MemoryStore) GetOrder(id string) (*domain.Order, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if o, ok := s.orders[id]; ok {
		return o, nil
	}
	return nil, ErrOrderNotFound
}

// GetOrdersByAccount возвращает список всех ордеров пользователя.
func (s *MemoryStore) GetOrdersByAccount(address string) ([]*domain.Order, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	var userOrders []*domain.Order
	for _, o := range s.orders {
		if o.AccountID == address {
			userOrders = append(userOrders, o)
		}
	}
	// TODO: Сортировать по времени создания (но в in-memory пока не храним timestamp на ордере, так что просто отдаем).
	return userOrders, nil
}

// UpdateOrder обновляет существу ордер.
func (s *MemoryStore) UpdateOrder(order *domain.Order) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.orders[order.ID]; !ok {
		return ErrOrderNotFound
	}
	s.orders[order.ID] = order
	return nil
}

// SaveTrade записывает завершенную сделку.
func (s *MemoryStore) SaveTrade(trade *domain.Trade) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.trades[trade.ID] = trade
	return nil
}
