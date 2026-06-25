package chain

// Client — это интерфейс, описывающий контракт взаимодействия DEX
// с внешней блокчейн-сетью (нодой Ethereum, Bitcoin Core и т.д.).
// Интерфейс позволяет легко подменить реальный блокчейн моком (MockChain) для тестов.
type Client interface {
	// VerifyTransaction проверяет, действительно ли транзакция включена в блок.
	VerifyTransaction(txHash string) (bool, error)
	
	// GetBalance запрашивает on-chain баланс пользователя.
	GetBalance(address, asset string) (string, error)
	
	// SubmitTransaction отправляет подписанную транзакцию вывода (Withdrawal) 
	// из DEX обратно в блокчейн.
	SubmitTransaction(signedTx []byte) (string, error)
}
