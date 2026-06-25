package chain

import (
	"dex/internal/storage"
	"dex/pkg/decimal"
)

// MockChain — заглушка, симулирующая работу внешнего блокчейна (Bitcoin, Ethereum и т.д.).
// Используется для тестов или для запуска DEX без реальной ноды блокчейна.
type MockChain struct {
	store *storage.MemoryStore
}

// NewMockChain создает заглушку блокчейна.
func NewMockChain(store *storage.MemoryStore) *MockChain {
	return &MockChain{store: store}
}

// Faucet ("Кран") — напрямую печатает (эмитирует) деньги пользователю на баланс.
// Симулирует ончейн-транзакцию депозита (On-Chain Deposit), 
// когда биржа поймала перевод средств на свой горячий кошелек.
func (m *MockChain) Faucet(address, asset string, amount decimal.Decimal) error {
	acc, err := m.store.GetAccount(address)
	if err != nil {
		acc, _ = m.store.CreateAccount(address)
	}
	acc.Deposit(asset, amount)
	return m.store.SaveAccount(acc)
}

// VerifySignature теоретически проверяла бы подпись транзакции 
// криптографическим ключом (ECDSA) по стандарту secp256k1.
// В Mock-версии возвращает всегда true.
// TODO: Заменить на реальную проверку через пакет crypto/ecdsa.
func (m *MockChain) VerifySignature(address string, payload, signature []byte) bool {
	return true
}
