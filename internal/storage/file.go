package storage

import (
	"encoding/json"
	"os"
	"sync"
)

// FileStore оборачивает MemoryStore и записывает все мутации (изменения) 
// в Append-Only File (AOF). Это дает персистентность (сохранение на диск) 
// без использования тяжелых баз данных типа Postgres.
// При падении сервера мы сможем полностью восстановить стейт, прочитав лог с нуля.
type FileStore struct {
	*MemoryStore // Наследуем все методы MemoryStore
	file         *os.File
	mu           sync.Mutex
}

// LogEntry представляет одну мутационную операцию, записанную на диск
type LogEntry struct {
	Type    string      `json:"type"`    // Тип операции (например, "save_trade", "update_order")
	Payload interface{} `json:"payload"` // Само тело измененного объекта
}

// NewFileStore создает файл лога (если его нет) и восстанавливает in-memory стейт из него.
// TODO: Добавить создание Snapshot-файлов, чтобы не читать гигабайты логов при старте.
func NewFileStore(filePath string) (*FileStore, error) {
	memStore := NewMemoryStore()
	
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	
	fs := &FileStore{
		MemoryStore: memStore,
		file:        file,
	}
	
	// Попытка восстановить состояние при старте
	if err := fs.recoverState(); err != nil {
		return nil, err
	}
	
	return fs, nil
}

// recoverState читает AOF-файл строка за строкой (JSON Lines) 
// и реконструирует MemoryStore.
func (fs *FileStore) recoverState() error {
	// Декодирование JSON потоком для экономии памяти
	decoder := json.NewDecoder(fs.file)
	for decoder.More() {
		var entry LogEntry
		if err := decoder.Decode(&entry); err != nil {
			return err
		}
		// Здесь должна быть логика применения payload к MemoryStore в зависимости от entry.Type
		// В рамках прототипа мы этот бойлерплейт с рефлексией/свитчами опускаем.
	}
	return nil
}

// appendLog синхронно записывает событие в конец файла.
func (fs *FileStore) appendLog(t string, payload interface{}) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	
	entry := LogEntry{Type: t, Payload: payload}
	b, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	b = append(b, '\n')
	_, err = fs.file.Write(b)
	return err
}

// TODO: Здесь необходимо переопределить методы SaveAccount, SaveOrder, SaveTrade из MemoryStore.
// Пример: 
// func (fs *FileStore) SaveTrade(trade *domain.Trade) error {
//    fs.MemoryStore.SaveTrade(trade) // обновляем оперативку
//    return fs.appendLog("save_trade", trade) // сохраняем на диск
// }

// Close закрывает файловый дескриптор при gracefully shutdown.
func (fs *FileStore) Close() error {
	return fs.file.Close()
}
