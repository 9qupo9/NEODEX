package decimal

import (
	"fmt"
	"math/big"
)

// Decimal — это обертка вокруг math/big.Rat для высокоточных финансовых вычислений.
// Использование float64 недопустимо в биржах из-за потери точности.
// TODO: Добавить пулы для big.Rat (sync.Pool), чтобы снизить нагрузку на Garbage Collector при HFT-нагрузках.
type Decimal struct {
	r *big.Rat
}

// Zero возвращает новый Decimal со значением 0.
// Используется для инициализации пустых балансов и объемов.
func Zero() Decimal {
	return Decimal{r: big.NewRat(0, 1)}
}

// NewFromString создает Decimal из строки (например, "1.5" или "0.0001").
// Возвращает ошибку, если формат строки неверный.
func NewFromString(s string) (Decimal, error) {
	r := new(big.Rat)
	if _, ok := r.SetString(s); !ok {
		return Decimal{}, fmt.Errorf("неверный формат числа: %s", s)
	}
	return Decimal{r: r}, nil
}

// MustParse парсит строку и паникует при ошибке.
// Полезно использовать только при инициализации констант при запуске приложения.
func MustParse(s string) Decimal {
	d, err := NewFromString(s)
	if err != nil {
		panic(err)
	}
	return d
}

// NewFromInt создает Decimal из обычного int64.
// TODO: Расширить поддержку для uint64 и big.Int для совместимости со смарт-контрактами (ERC20).
func NewFromInt(i int64) Decimal {
	return Decimal{r: big.NewRat(i, 1)}
}

// Add прибавляет другое значение Decimal и возвращает новый Decimal.
func (d Decimal) Add(other Decimal) Decimal {
	res := new(big.Rat).Add(d.r, other.r)
	return Decimal{r: res}
}

// Sub вычитает другое значение Decimal и возвращает новый Decimal.
func (d Decimal) Sub(other Decimal) Decimal {
	res := new(big.Rat).Sub(d.r, other.r)
	return Decimal{r: res}
}

// Mul умножает на другое значение Decimal и возвращает новый Decimal.
func (d Decimal) Mul(other Decimal) Decimal {
	res := new(big.Rat).Mul(d.r, other.r)
	return Decimal{r: res}
}

// Div делит на другое значение Decimal.
// ВАЖНО: Если делитель равен нулю, функция вызывает панику. Защита от panic должна быть на уровне бизнес-логики.
func (d Decimal) Div(other Decimal) Decimal {
	if other.r.Sign() == 0 {
		panic("деление на ноль")
	}
	res := new(big.Rat).Quo(d.r, other.r)
	return Decimal{r: res}
}

// Cmp сравнивает два числа Decimal.
// Возвращает: -1 если d < other, 0 если d == other, +1 если d > other.
func (d Decimal) Cmp(other Decimal) int {
	return d.r.Cmp(other.r)
}

// Sign возвращает знак числа.
// Возвращает: -1 если число отрицательное, 0 если ноль, +1 если положительное.
func (d Decimal) Sign() int {
	return d.r.Sign()
}

// String форматирует Decimal в строку с фиксированным количеством знаков после запятой (по умолчанию 8).
// TODO: Сделать динамическое количество знаков после запятой в зависимости от настроек торговой пары (Pair.QtyStep).
func (d Decimal) String() string {
	return d.r.FloatString(8)
}

// IsZero возвращает true, если значение равно нулю.
func (d Decimal) IsZero() bool {
	return d.r.Sign() == 0
}

// Copy создает глубокую копию Decimal, чтобы избежать мутаций при передаче по ссылке.
func (d Decimal) Copy() Decimal {
	res := new(big.Rat).Set(d.r)
	return Decimal{r: res}
}
