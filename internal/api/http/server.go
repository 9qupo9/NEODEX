package http

import (
	"dex/internal/api/ws"
	"log"
	"net/http"
	"dex/internal/ui"
	"encoding/json"
	"io"
	"sync"
	"time"
)

var (
	cachedIcons map[string]string
	iconsMutex  sync.RWMutex
	lastFetch   time.Time
	iconBytesCache sync.Map
)

func getIconsMap() map[string]string {
	iconsMutex.RLock()
	cacheValid := time.Since(lastFetch) < 1*time.Hour
	iconsMutex.RUnlock()

	if cacheValid && cachedIcons != nil {
		iconsMutex.RLock()
		defer iconsMutex.RUnlock()
		return cachedIcons
	}

	// Fetch from Binance
	resp, err := http.Get("https://www.binance.com/bapi/asset/v2/public/asset/asset/get-all-asset")
	if err != nil {
		log.Println("Error fetching icons:", err)
		return map[string]string{}
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)
	var data map[string]interface{}
	json.Unmarshal(b, &data)
	
	newMap := make(map[string]string)
	if data["data"] != nil {
		assets := data["data"].([]interface{})
		for _, v := range assets {
			asset := v.(map[string]interface{})
			code, _ := asset["assetCode"].(string)
			url, _ := asset["logoUrl"].(string)
			if code != "" && url != "" {
				newMap[code] = url
			}
		}
	}

	iconsMutex.Lock()
	cachedIcons = newMap
	lastFetch = time.Now()
	iconsMutex.Unlock()

	return newMap
}

// Server — это стандартный REST HTTP сервер.
// Предоставляет внешнее API для Web-клиентов и мобильных приложений.
type Server struct {
	addr     string
	handlers *Handlers
	wsHub    *ws.Hub // Ссылка на хаб для апгрейда HTTP до WebSocket
}

// NewServer создает HTTP сервер с роутингом (Mux).
func NewServer(addr string, handlers *Handlers, wsHub *ws.Hub) *Server {
	return &Server{
		addr:     addr,
		handlers: handlers,
		wsHub:    wsHub,
	}
}

// Start запускает HTTP-сервер. Блокирующий вызов.
func (s *Server) Start() error {
	// прямо в паттернах путей, поэтому сторонние роутеры вроде Chi или Gorilla не нужны.
	mux := http.NewServeMux()

	// REST API маршруты
	mux.HandleFunc("/api/v1/order", s.handlers.HandlePlaceOrder)
	mux.HandleFunc("/api/v1/order/cancel", s.handlers.HandleCancelOrder)
	mux.HandleFunc("/api/v1/orderbook", s.handlers.HandleGetOrderbook)
	mux.HandleFunc("/api/v1/balance", s.handlers.HandleGetBalance)
	mux.HandleFunc("/api/v1/orders", s.handlers.HandleGetOrders)
	
	// Initialize New Listings Scanner
	InitScanner()

	// Binance Icons Proxy
	mux.HandleFunc("/api/v1/icon", func(w http.ResponseWriter, r *http.Request) {
		symbol := r.URL.Query().Get("s")
		if symbol == "" {
			http.Error(w, "missing symbol", http.StatusBadRequest)
			return
		}

		// Check memory cache
		if img, ok := iconBytesCache.Load(symbol); ok {
			w.Header().Set("Content-Type", "image/png")
			w.Header().Set("Cache-Control", "public, max-age=86400")
			w.Write(img.([]byte))
			return
		}

		icons := getIconsMap()
		url, exists := icons[symbol]
		if !exists || url == "" {
			http.NotFound(w, r)
			return
		}

		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("User-Agent", "Mozilla/5.0")
		resp, err := http.DefaultClient.Do(req)
		if err != nil || resp.StatusCode != 200 {
			http.NotFound(w, r)
			return
		}
		defer resp.Body.Close()

		imgBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "error reading image", http.StatusInternalServerError)
			return
		}

		iconBytesCache.Store(symbol, imgBytes)
		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Cache-Control", "public, max-age=86400")
		w.Write(imgBytes)
	})

	// New Listings API
	mux.HandleFunc("/api/v1/new-listings", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(GetNewListings())
	})

	// WebSocket маршрут (передача соединения из HTTP в WS)
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.Upgrade(w, r, s.wsHub)
	})

	// Redirect root to /markets
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, "/markets", http.StatusTemporaryRedirect)
	})

	// Markets Page
	mux.HandleFunc("/markets", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(ui.RenderMarketsPage()))
	})

	// Spot Trading Terminal
	mux.HandleFunc("/spot", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(ui.RenderSpotPage()))
	})

	// Wallet / Portfolio Dashboard
	mux.HandleFunc("/wallet", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(ui.RenderWalletPage()))
	})

	log.Printf("HTTP REST сервер запущен на %s", s.addr)
	return http.ListenAndServe(s.addr, mux)
}
