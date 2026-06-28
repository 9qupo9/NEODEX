package syslog

import (
	"bytes"
	"io"
	"os"
	"sync"
	"time"
)

// LogEntry представляет одну запись лога.
type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
}

// RingBuffer перехватывает логи и сохраняет последние N записей в памяти.
// Также он дублирует вывод в стандартный os.Stdout.
type RingBuffer struct {
	mu       sync.Mutex
	capacity int
	entries  []LogEntry
	out      io.Writer
}

// NewRingBuffer создает новый буфер для логов.
func NewRingBuffer(capacity int) *RingBuffer {
	return &RingBuffer{
		capacity: capacity,
		entries:  make([]LogEntry, 0, capacity),
		out:      os.Stdout,
	}
}

// Write реализует интерфейс io.Writer.
func (b *RingBuffer) Write(p []byte) (n int, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Убираем финальный перевод строки, если он есть, для чистоты JSON
	msg := string(bytes.TrimSpace(p))
	
	entry := LogEntry{
		Timestamp: time.Now(),
		Message:   msg,
	}

	if len(b.entries) >= b.capacity {
		// Сдвигаем влево и добавляем в конец
		b.entries = append(b.entries[1:], entry)
	} else {
		b.entries = append(b.entries, entry)
	}

	// Дублируем в консоль
	return b.out.Write(p)
}

// GetLogs возвращает копию последних логов.
func (b *RingBuffer) GetLogs() []LogEntry {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	cpy := make([]LogEntry, len(b.entries))
	copy(cpy, b.entries)
	return cpy
}

// Clear очищает буфер логов.
func (b *RingBuffer) Clear() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.entries = make([]LogEntry, 0, b.capacity)
}

// Глобальный инстанс буфера логов (храним последние 50 записей).
var GlobalBuffer = NewRingBuffer(50)
