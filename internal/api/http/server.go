package http

import (
	"dex/internal/api/ws"
	"log"
	"net/http"
	"dex/internal/ui"
	"dex/internal/ui/admin"
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
	auth     *AuthHandler
	staking  *StakingHandler
	wsHub    *ws.ShardedHub // Ссылка на шардированный хаб
}

// NewServer создает HTTP сервер с роутингом (Mux).
func NewServer(addr string, handlers *Handlers, wsHub *ws.ShardedHub) *Server {
	authService := NewMockWeb3AuthService()
	authHandler := NewAuthHandler(authService, handlers.Store)
	stakingHandler := NewStakingHandler(handlers.Store)

	return &Server{
		addr:     addr,
		handlers: handlers,
		auth:     authHandler,
		staking:  stakingHandler,
		wsHub:    wsHub,
	}
}

// Start запускает HTTP-сервер. Блокирующий вызов.
func (s *Server) Start() error {
	// прямо в паттернах путей, поэтому сторонние роутеры вроде Chi или Gorilla не нужны.
	mux := http.NewServeMux()

	// REST API маршруты (Торговля и баланс)
	mux.HandleFunc("/api/v1/order", s.handlers.HandlePlaceOrder)
	mux.HandleFunc("/api/v1/order/cancel", s.handlers.HandleCancelOrder)
	mux.HandleFunc("/api/v1/orderbook", s.handlers.HandleGetOrderbook)
	mux.HandleFunc("/api/v1/balance", s.handlers.HandleGetBalance)
	mux.HandleFunc("/api/v1/orders", s.handlers.HandleGetOrders)

	// Auth Web3
	mux.HandleFunc("/api/v1/auth/nonce", s.auth.HandleGetNonce)
	mux.HandleFunc("/api/v1/auth/verify", s.auth.HandleVerify)

	// DeFi (Staking)
	mux.HandleFunc("/api/v1/staking/stake", s.staking.HandleStake)
	mux.HandleFunc("/api/v1/staking/unstake", s.staking.HandleUnstake)
	mux.HandleFunc("/api/v1/staking/list", s.staking.HandleGetStakes)

	// Фьючерсные маршруты
	mux.HandleFunc("/api/v1/futures/order", s.handlers.HandlePlaceFuturesOrder)
	mux.HandleFunc("/api/v1/futures/positions", s.handlers.HandleGetFuturesPositions)
	
	// Admin API маршруты
	mux.HandleFunc("/api/v1/admin/users", s.handlers.HandleAdminUsers)
	mux.HandleFunc("/api/v1/admin/users/block", s.handlers.HandleAdminBlockUser)
	// Админ API
	// Админ API
	mux.HandleFunc("/api/v1/admin/metrics", s.handlers.HandleAdminMetrics)
	mux.HandleFunc("/api/v1/admin/action", s.handlers.HandleAdminAction)

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

	// Wallet & DeFi Earn
	mux.HandleFunc("/wallet", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(ui.RenderWalletPage()))
	})

	// Futures Trading Terminal
	mux.HandleFunc("/futures", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(ui.RenderFuturesPage()))
	})

	// Admin Dashboard
	mux.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		// Простая HTML-авторизация вместо Basic Auth (чтобы избежать ERR_TOO_MANY_REDIRECTS или багов браузера)
		cookie, err := r.Cookie("admin_session")
		if err != nil || cookie.Value != "authenticated" {
			if r.Method == http.MethodPost && r.FormValue("password") == "admin123" {
				http.SetCookie(w, &http.Cookie{
					Name:     "admin_session",
					Value:    "authenticated",
					Path:     "/",
					HttpOnly: true,
					MaxAge:   3600,
				})
				http.Redirect(w, r, "/admin", http.StatusFound)
				return
			}
			
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write([]byte(`
				<!DOCTYPE html>
				<html>
				<head>
					<title>Admin Login</title>
					<style>
						body { background: #0b0e11; color: #fff; font-family: sans-serif; display: flex; justify-content: center; align-items: center; height: 100vh; margin: 0; }
						.login-box { background: #1e2329; padding: 40px; border-radius: 8px; text-align: center; border: 1px solid #333; }
						input { padding: 10px; border-radius: 4px; border: 1px solid #333; background: #2b3139; color: #fff; margin-bottom: 20px; width: 200px; }
						button { padding: 10px 20px; background: #fcd535; color: #000; border: none; border-radius: 4px; cursor: pointer; font-weight: bold; }
						button:hover { background: #f0c828; }
					</style>
				</head>
				<body>
					<div class="login-box">
						<h2>Admin Panel</h2>
						<form method="POST">
							<input type="password" name="password" placeholder="Enter Password" required><br>
							<button type="submit">Login</button>
						</form>
					</div>
				</body>
				</html>
			`))
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(admin.RenderLayout()))
	})
	log.Printf("\n=======================================================")
	log.Printf("HTTP REST сервер успешно запущен на %s", s.addr)
	log.Printf("Главная (Рынки):  http://localhost%s/", s.addr)
	log.Printf("Спот терминал:    http://localhost%s/spot", s.addr)
	log.Printf("Фьючерсы:         http://localhost%s/futures", s.addr)
	log.Printf("Кошелек (Earn):   http://localhost%s/wallet", s.addr)
	log.Printf("-------------------------------------------------------")
	log.Printf("Админ-панель:     http://localhost%s/admin", s.addr)
	log.Printf("Вход в админку:   Логин: admin | Пароль: admin123")
	log.Printf("=======================================================\n")
	
	return http.ListenAndServe(s.addr, mux)
}
