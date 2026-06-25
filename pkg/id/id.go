package id

import (
	"crypto/rand"
	"encoding/hex"
	"sync/atomic"
	"time"
)

// counter используется как атомарный счетчик для гарантии уникальности ID в пределах одной миллисекунды.
var counter uint64

// New генерирует псевдо-уникальный ID, объединяя текущее время, атомарный счетчик и случайные байты.
// Это очень быстрый генератор, подходящий для внутренних нужд DEX (ордера, трейды).
// Без использования сторонних библиотек вроде UUID или Snowflake.
//
// TODO: В будущем, при масштабировании на несколько серверов (кластер), необходимо добавить
// идентификатор "Worker ID" или "Node ID", чтобы избежать коллизий между разными машинами.
func New() string {
	c := atomic.AddUint64(&counter, 1)
	t := time.Now().UnixNano()
	
	// 4 байта случайных данных для повышения энтропии
	b := make([]byte, 4)
	_, _ = rand.Read(b)
	
	// Формат сборки: {Время в hex} - {Счетчик в hex} - {Рандом в hex}
	return hex.EncodeToString(append(uint64ToBytes(uint64(t)), append(uint64ToBytes(c), b...)...))
}

// uint64ToBytes конвертирует uint64 в массив из 8 байт (BigEndian).
// Используется вручную для максимальной производительности без рефлексии (пакет binary).
func uint64ToBytes(v uint64) []byte {
	b := make([]byte, 8)
	b[0] = byte(v >> 56)
	b[1] = byte(v >> 48)
	b[2] = byte(v >> 40)
	b[3] = byte(v >> 32)
	b[4] = byte(v >> 24)
	b[5] = byte(v >> 16)
	b[6] = byte(v >> 8)
	b[7] = byte(v)
	return b
}
