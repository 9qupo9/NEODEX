package storage

import (
	"dex/internal/domain"
	"dex/pkg/decimal"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
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
	Type    string          `json:"type"`    // Тип операции (например, "save_trade", "update_order")
	Payload json.RawMessage `json:"payload"` // Само тело измененного объекта
}

// NewFileStore создает файл лога (если его нет) и восстанавливает in-memory стейт из него.
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
	// Сбрасываем указатель файла в начало
	fs.file.Seek(0, 0)
	decoder := json.NewDecoder(fs.file)
	for decoder.More() {
		var entry LogEntry
		if err := decoder.Decode(&entry); err != nil {
			return err
		}
		
		switch entry.Type {
		case "create_account":
			var address string
			if err := json.Unmarshal(entry.Payload, &address); err == nil {
				fs.MemoryStore.CreateAccount(address)
			}
		case "save_order":
			var order domain.Order
			if err := json.Unmarshal(entry.Payload, &order); err == nil {
				fs.MemoryStore.SaveOrder(&order)
			}
		case "update_order":
			var order domain.Order
			if err := json.Unmarshal(entry.Payload, &order); err == nil {
				fs.MemoryStore.UpdateOrder(&order)
			}
		case "settle_trade":
			var trade domain.Trade
			if err := json.Unmarshal(entry.Payload, &trade); err == nil {
				fs.MemoryStore.SettleTradeBalances(&trade)
			}
		case "save_position":
			var position domain.Position
			if err := json.Unmarshal(entry.Payload, &position); err == nil {
				fs.MemoryStore.SavePosition(&position)
			}
		case "deposit":
			var payload struct {
				Address string `json:"address"`
				Asset   string `json:"asset"`
				Amount  string `json:"amount"`
			}
			if err := json.Unmarshal(entry.Payload, &payload); err == nil {
				if acc, err := fs.MemoryStore.GetAccount(payload.Address); err == nil {
					if amt, err := decimal.NewFromString(payload.Amount); err == nil {
						acc.Deposit(payload.Asset, amt)
					}
				}
			}
		}
	}
	return nil
}

// appendLog синхронно записывает событие в конец файла.
func (fs *FileStore) appendLog(t string, payload interface{}) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	
	bPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	
	entry := LogEntry{Type: t, Payload: bPayload}
	b, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	b = append(b, '\n')
	_, err = fs.file.Write(b)
	return err
}

func (fs *FileStore) CreateAccount(address string) (*domain.Account, error) {
	acc, err := fs.MemoryStore.CreateAccount(address)
	if err == nil {
		fs.appendLog("create_account", address)
	}
	return acc, err
}

func (fs *FileStore) SaveOrder(order *domain.Order) error {
	err := fs.MemoryStore.SaveOrder(order)
	if err == nil {
		fs.appendLog("save_order", order)
	}
	return err
}

func (fs *FileStore) UpdateOrder(order *domain.Order) error {
	err := fs.MemoryStore.UpdateOrder(order)
	if err == nil {
		fs.appendLog("update_order", order)
	}
	return err
}

func (fs *FileStore) SettleTradeBalances(trade *domain.Trade) error {
	err := fs.MemoryStore.SettleTradeBalances(trade)
	if err == nil {
		fs.appendLog("settle_trade", trade)
	}
	return err
}

func (fs *FileStore) SavePosition(position *domain.Position) error {
	err := fs.MemoryStore.SavePosition(position)
	if err == nil {
		fs.appendLog("save_position", position)
	}
	return err
}

// Close закрывает файловый дескриптор при gracefully shutdown.
func (fs *FileStore) Close() error {
	return fs.file.Close()
}

// Snapshot делает полную перезапись AOF файла текущим состоянием in-memory.
// Это компактизация данных.
func (fs *FileStore) Snapshot() error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	// Получаем текущее состояние
	accounts, orders, positions := fs.MemoryStore.DumpState()

	tempPath := fs.file.Name() + ".tmp"
	tempFile, err := os.OpenFile(tempPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	// Функция для записи в темповый файл
	writeEntry := func(t string, payload interface{}) error {
		bPayload, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		entry := LogEntry{Type: t, Payload: bPayload}
		b, err := json.Marshal(entry)
		if err != nil {
			return err
		}
		b = append(b, '\n')
		_, err = tempFile.Write(b)
		return err
	}

	// Записываем балансы (как депозиты)
	for _, acc := range accounts {
		writeEntry("create_account", acc.Address)
		
		for asset, amount := range acc.GetBalances() {
			if amount.Cmp(decimal.Zero()) > 0 {
				payload := struct {
					Address string `json:"address"`
					Asset   string `json:"asset"`
					Amount  string `json:"amount"`
				}{
					Address: acc.Address,
					Asset:   asset,
					Amount:  amount.String(),
				}
				writeEntry("deposit", payload)
			}
		}
	}

	// Записываем активные ордера (со статусами)
	for _, ord := range orders {
		if ord.Status == domain.StatusNew || ord.Status == domain.StatusPartiallyFilled {
			writeEntry("save_order", ord)
		}
	}

	// Записываем позиции
	for _, pos := range positions {
		writeEntry("save_position", pos)
	}

	// Закрываем темповый файл
	if err := tempFile.Sync(); err != nil {
		tempFile.Close()
		return err
	}
	if err := tempFile.Close(); err != nil {
		return err
	}

	// Подменяем старый файл
	oldPath := fs.file.Name()
	if err := fs.file.Close(); err != nil {
		return err
	}
	if err := os.Rename(tempPath, oldPath); err != nil {
		// Попытка переоткрыть старый файл
		fs.file, _ = os.OpenFile(oldPath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
		return err
	}

	// Открываем новый файл
	newFile, err := os.OpenFile(oldPath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	fs.file = newFile

	// ДОПОЛНИТЕЛЬНО: Создаем красивую копию в папку snapshots/ для администратора
	snapshotDir := "snapshots"
	if err := os.MkdirAll(snapshotDir, 0755); err == nil {
		timestamp := time.Now().Format("2006-01-02_15-04-05")
		backupPath := filepath.Join(snapshotDir, fmt.Sprintf("snapshot_%s.jsonl", timestamp))
		
		// Копируем сжатый старый файл как новый снапшот
		input, err := os.ReadFile(oldPath)
		if err == nil {
			os.WriteFile(backupPath, input, 0644)
		}
	}

	return nil
}
