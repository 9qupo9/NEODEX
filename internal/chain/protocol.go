package chain

// TxRequest описывает структуру запроса на создание транзакции.
// Используется для парсинга JSON-сообщений от внешних нод блокчейна.
// TODO: Заменить string на кастомный тип decimal.Decimal для поля Amount, добавив MarshalJSON/UnmarshalJSON.
type TxRequest struct {
	From   string `json:"from"`   // Адрес отправителя (User Wallet)
	To     string `json:"to"`     // Адрес получателя (DEX Hot Wallet)
	Amount string `json:"amount"` // Сумма перевода
	Asset  string `json:"asset"`  // Тикер актива (например, ETH или USDC)
	Sig    string `json:"sig"`    // Криптографическая подпись (Hex)
}

// TxResponse описывает ответ биржи на ончейн-транзакцию.
type TxResponse struct {
	TxHash string `json:"txHash"` // Хэш транзакции в блокчейне
	Status string `json:"status"` // Статус (Pending, Success, Failed)
}
