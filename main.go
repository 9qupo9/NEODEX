package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"dex/internal/api/http"
	"dex/internal/api/tcp"
	"dex/internal/api/ws"
	"dex/internal/config"
	"dex/internal/domain"
	"dex/internal/engine"
	"dex/internal/storage"
	"dex/pkg/decimal"
	"dex/pkg/syslog"
)

// main — главная точка входа, где "собирается" вся архитектура (Dependency Injection).
func main() {
	// Перехватываем стандартный логгер для вывода в админ-панель
	log.SetOutput(syslog.GlobalBuffer)
	log.SetFlags(0) // Убираем стандартные префиксы времени (будем добавлять свои или парсить)

	log.Println("[SYSTEM] Запуск DEX Ноды...")

	// 1. Загрузка конфигурации
	cfg, _ := config.LoadConfig("config.json")

	// 2. Инициализация хранилища (балансы, стейт)
	store, err := storage.NewFileStore("dex_state.aof")
	if err != nil {
		log.Fatalf("Ошибка загрузки хранилища: %v", err)
	}

	// 4. Создаем канал Broadcaster для рассылки рыночных данных по WebSocket
	fromEngine := make(chan []byte, 1024)

	// 5. Инициализация Маршрутизатора
	router := engine.NewEngineRouter(store, fromEngine)

	// 6. Создаем торговую пару с правилами (Tick Size, Lot Size)
	pair := domain.Pair{
		BaseAsset:  "BTC",
		QuoteAsset: "USDT",
		PriceStep:  decimal.MustParse("0.01"),
		QtyStep:    decimal.MustParse("0.0001"),
	}

	router.CreateEngine(pair) // Создаем и регистрируем BTC/USDT

	// Для фьючерсов оставим жесткую инициализацию (или тоже можно через роутер в будущем)
	futuresEng := engine.NewFuturesEngine(pair, store, fromEngine)


	// Можно добавить еще пару для тестов, например ETH/USDT
	ethPair := domain.Pair{
		BaseAsset:  "ETH",
		QuoteAsset: "USDT",
		PriceStep:  decimal.MustParse("0.01"),
		QtyStep:    decimal.MustParse("0.0001"),
	}
	router.CreateEngine(ethPair)


	// 5. Инициализация WebSocket Хаба (Шардированный для масштабирования)
	wsHub := ws.NewShardedHub(256) // 256 шардов для максимальной утилизации CPU
	wsHub.Run() // Запуск Event Loops для каждого шарда
	
	// Горутина-мост: перекладываем JSON-сообщения из Engine Broadcaster в WS Hub
	go func() {
		for msg := range fromEngine {
			wsHub.Broadcast(msg)
		}
	}()

	// 6. Инициализация низкоуровневого TCP сервера (для HFT ботов)
	tcpHandler := tcp.NewHandler(router)
	tcpServer := tcp.NewServer(cfg.TcpAddr, tcpHandler)
	go func() {
		if err := tcpServer.Start(); err != nil {
			log.Fatalf("Ошибка запуска TCP сервера: %v", err)
		}
	}()

	// 7. Инициализация HTTP REST сервера (передаем и ядро, и WS хаб, и tcpHandler для метрик)
	httpHandlers := http.NewHandlers(router, futuresEng, store, tcpHandler, wsHub)
	httpServer := http.NewServer(cfg.HttpAddr, httpHandlers, wsHub)
	go func() {
		if err := httpServer.Start(); err != nil {
			log.Fatalf("Ошибка запуска HTTP сервера: %v", err)
		}
	}()

	// 8. Graceful Shutdown (Изящное завершение)
	// Перехватываем сигналы остановки (Ctrl+C, SIGTERM от Docker), 
	// чтобы успеть записать AOF лог на диск перед смертью процесса.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Изящное завершение работы (Graceful shutdown)...")
}
