package domain

import "dex/pkg/decimal"

// Pair представляет торговую пару на бирже (например, BTC/USDT).
// Содержит правила торговли (Trading Rules) для защиты стакана и движка от спама.
type Pair struct {
	BaseAsset  string // Базовый актив (то, что мы покупаем/продаем). Пример: BTC
	QuoteAsset string // Котируемый актив (то, чем мы расплачиваемся). Пример: USDT
	
	// Ограничения по цене для предотвращения ошибочных ордеров (Fat Finger protection)
	// TODO: Сделать динамический MinPrice/MaxPrice на основе текущей рыночной цены +- 20%.
	MinPrice decimal.Decimal
	MaxPrice decimal.Decimal
	
	// Ограничения по объему (lot size limits)
	MinQty   decimal.Decimal
	MaxQty   decimal.Decimal
	
	// Шаг цены (Tick Size). Например, 0.01 значит, что цена 1.234 недопустима, только 1.23.
	PriceStep decimal.Decimal 
	
	// Шаг объема (Lot Size). Ограничивает количество знаков после запятой в количестве актива.
	QtyStep   decimal.Decimal 
}

// String возвращает строковое представление пары. Пример: "BTC-USDT"
// Используется для логирования и маршрутизации.
func (p Pair) String() string {
	return p.BaseAsset + "-" + p.QuoteAsset
}

// Equals проверяет, идентичны ли две торговые пары по базовому и котируемому активу.
func (p Pair) Equals(other Pair) bool {
	return p.BaseAsset == other.BaseAsset && p.QuoteAsset == other.QuoteAsset
}
