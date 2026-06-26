package storage

import "dex/internal/domain"

// SavePosition сохраняет или обновляет фьючерсную позицию
func (s *MemoryStore) SavePosition(position *domain.Position) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.positions[position.ID] = position
	return nil
}

// GetPosition возвращает позицию по ID
func (s *MemoryStore) GetPosition(id string) (*domain.Position, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if p, ok := s.positions[id]; ok {
		return p, nil
	}
	return nil, ErrOrderNotFound // Позже можно добавить ErrPositionNotFound
}

// GetPositionsByAccount возвращает список всех позиций пользователя
func (s *MemoryStore) GetPositionsByAccount(address string) ([]*domain.Position, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	var userPositions []*domain.Position
	for _, p := range s.positions {
		if p.AccountID == address {
			userPositions = append(userPositions, p)
		}
	}
	return userPositions, nil
}
