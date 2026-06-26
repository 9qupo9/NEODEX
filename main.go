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
)

// main — главная точка входа, где "собирается" вся архитектура (Dependency Injection).
func main() {
	log.Println("Запуск DEX Ноды...")

	// 1. Загрузка конфигурации
	cfg, _ := config.LoadConfig("config.json")

	// 2. Инициализация хранилища (балансы, стейт)
	store, err := storage.NewFileStore("dex_state.aof")
	if err != nil {
		log.Fatalf("Ошибка загрузки хранилища: %v", err)
	}

	// 3. Создаем торговую пару с правилами (Tick Size, Lot Size)
	pair := domain.Pair{
		BaseAsset:  "BTC",
		QuoteAsset: "USDT",
		PriceStep:  decimal.MustParse("0.01"),
		QtyStep:    decimal.MustParse("0.0001"),
	}

	// 4. Инициализация Торговых Ядер и Маршрутизатора
	router := engine.NewEngineRouter()

	// Добавляем Broadcaster для рассылки рыночных данных по WebSocket
	fromEngine := make(chan []byte, 1024)
	eng := engine.NewEngine(pair, store, fromEngine)
	futuresEng := engine.NewFuturesEngine(pair, store, fromEngine)
	
	router.RegisterEngine(eng) // Регистрация BTC_USDT в роутере

	// Можно добавить еще пару для тестов, например ETH/USDT
	ethPair := domain.Pair{
		BaseAsset:  "ETH",
		QuoteAsset: "USDT",
		PriceStep:  decimal.MustParse("0.01"),
		QtyStep:    decimal.MustParse("0.0001"),
	}
	ethEng := engine.NewEngine(ethPair, store, fromEngine)
	router.RegisterEngine(ethEng)

	// -- МОК ДАННЫХ ДЛЯ ТЕСТИРОВАНИЯ --
	// Печатаем немного тестовых денег пользователю, чтобы он мог отправить ордер
	acc, _ := store.CreateAccount("test_user_1")
	acc.Deposit("USDT", decimal.MustParse("100000"))
	acc.Deposit("BTC", decimal.MustParse("10"))
	acc.Deposit("ETH", decimal.MustParse("50"))
	
	// ----------------------------------

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

	// 7. Инициализация HTTP REST сервера (передаем и ядро, и WS хаб)
	httpHandlers := http.NewHandlers(router, futuresEng, store)
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
