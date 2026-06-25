package config

import (
	"encoding/json"
	"os"
)

// Config хранит основные настройки для старта ноды DEX.
type Config struct {
	HttpAddr string `json:"httpAddr"` // Адрес для REST и WebSocket (например, ":8080")
	TcpAddr  string `json:"tcpAddr"`  // Адрес для HFT TCP P2P протокола (например, ":9000")
	DataDir  string `json:"dataDir"`  // Директория для хранения AOF (логов транзакций)
}

// LoadConfig пытается прочитать настройки из файла config.json.
// Если файл не найден (например, при первом запуске), возвращает настройки по умолчанию.
func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		// Возвращаем настройки по умолчанию, если файл не существует
		return &Config{
			HttpAddr: ":8080",
			TcpAddr:  ":9000",
			DataDir:  "./data",
		}, nil
	}
	defer file.Close()

	var cfg Config
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
