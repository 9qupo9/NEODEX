package http

import (
	"crypto/rand"
	"dex/internal/storage"
	"dex/pkg/decimal"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"
)

// IAuthService определяет контракт для сервиса авторизации.
// HTTP слой (Handlers) зависит от интерфейса, а не от конкретной реализации JWT или Web3.
type IAuthService interface {
	GenerateNonce(address string) (string, error)
	VerifySignature(address, signature, expectedMessage string) (string, error)
}

// Принцип SOLID: Single Responsibility Principle (SRP)
// MockWeb3AuthService отвечает исключительно за генерацию токенов сессий и проверку Web3 авторизации.
// Пока что мы используем заглушку (Mock) вместо тяжеловесной go-ethereum криптографии, чтобы ускорить работу.
type MockWeb3AuthService struct {
	mu       sync.Mutex
	nonces   map[string]string // Хранит случайные строки для подписи
	sessions map[string]string // Хранит активные сессии (token -> address)
}

func NewMockWeb3AuthService() *MockWeb3AuthService {
	return &MockWeb3AuthService{
		nonces:   make(map[string]string),
		sessions: make(map[string]string),
	}
}

// GenerateNonce генерирует уникальную строку для защиты от Replay Attack.
func (s *MockWeb3AuthService) GenerateNonce(address string) (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	nonce := hex.EncodeToString(bytes)

	s.mu.Lock()
	s.nonces[strings.ToLower(address)] = nonce
	s.mu.Unlock()

	return nonce, nil
}

// VerifySignature проверяет EIP-712 подпись (Mock) и выдает сессионный токен.
func (s *MockWeb3AuthService) VerifySignature(address, signature, _ string) (string, error) {
	// TODO: Интегрировать go-ethereum/crypto для реальной проверки ECRECOVER
	// Временно симулируем успех, если передана любая подпись (так как это тестовая среда)

	bytes := make([]byte, 32)
	rand.Read(bytes)
	token := hex.EncodeToString(bytes)

	s.mu.Lock()
	s.sessions[token] = strings.ToLower(address)
	s.mu.Unlock()

	return token, nil
}

// AuthHandler отвечает исключительно за обработку HTTP-маршрутов аутентификации (SRP).
type AuthHandler struct {
	AuthService IAuthService
	Store       storage.Store
}

func NewAuthHandler(authService IAuthService, store storage.Store) *AuthHandler {
	return &AuthHandler{
		AuthService: authService,
		Store:       store,
	}
}

// HandleGetNonce выдает Nonce для подписи кошельком (Метамаск).
func (h *AuthHandler) HandleGetNonce(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "Требуется адрес кошелька", http.StatusBadRequest)
		return
	}

	nonce, err := h.AuthService.GenerateNonce(address)
	if err != nil {
		http.Error(w, "Ошибка генерации nonce", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"nonce": nonce})
}

// HandleVerify обрабатывает входящую подпись и регистрирует/авторизует пользователя.
func (h *AuthHandler) HandleVerify(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Address   string `json:"address"`
		Signature string `json:"signature"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}

	token, err := h.AuthService.VerifySignature(req.Address, req.Signature, "")
	if err != nil {
		http.Error(w, "Ошибка верификации подписи", http.StatusUnauthorized)
		return
	}

	// Автоматически регистрируем пользователя в хранилище (Storage), если его там еще нет
	_, err = h.Store.GetAccount(req.Address)
	if err != nil { // Аккаунт не найден, создаем
		acc, _ := h.Store.CreateAccount(req.Address)
		// Начисляем тестовые токены новому юзеру для тестнета (Faucet)
		// TODO: В будущем убрать эту раздачу и сделать пополнение смарт-контрактом
		acc.Deposit("USDT", decimal.MustParse("100000"))
		acc.Deposit("BTC", decimal.MustParse("10"))
		acc.Deposit("ETH", decimal.MustParse("50"))
		h.Store.SaveAccount(acc)
	}

	// Отправляем токен (в реальном приложении лучше через HttpOnly Cookie)
	http.SetCookie(w, &http.Cookie{
		Name:     "neodex_session",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(24 * time.Hour),
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token":   token,
		"address": req.Address,
		"message": "Успешная авторизация по Web3 подписи",
	})
}
