package ws

import "sync"

// Hub выступает в роли диспетчера WebSocket подключений.
// Он управляет списком всех активных клиентов и безопасно рассылает (broadcasting) им рыночные данные.
type Hub struct {
	clients    map[*Client]bool // Набор всех активных клиентов
	broadcast  chan []byte      // Входящий канал с сообщениями, которые нужно разослать всем
	register   chan *Client     // Канал для регистрации новых подключений
	unregister chan *Client     // Канал для отключения клиентов
	mu         sync.Mutex       // Мьютекс для защиты карты clients
}

// NewHub инициализирует каналы хаба.
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

// Broadcast отправляет сырые байты (JSON) в канал рассылки.
func (h *Hub) Broadcast(msg []byte) {
	h.broadcast <- msg
}

// Run — бесконечный цикл (Event Loop) мультиплексирования.
// Обрабатывает регистрацию, отключение и рассылку сообщений.
// Запускается как отдельная горутина при старте приложения (main.go).
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			// Подключился новый клиент
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			
		case client := <-h.unregister:
			// Клиент отвалился (закрыл вкладку браузера или пропал интернет)
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client) // Удаляем из пула
				close(client.send)        // Закрываем его персональный канал, чтобы writePump завершился
			}
			h.mu.Unlock()
			
		case message := <-h.broadcast:
			// Пришло новое сообщение (например, новый стакан из Broadcaster)
			h.mu.Lock()
			for client := range h.clients {
				select {
				case client.send <- message:
					// Сообщение успешно отправлено в буфер клиента
				default:
					// Если буфер клиента (256 элементов) переполнен, значит клиент "виснет".
					// Отключаем его насильно, чтобы он не тормозил всю рассылку.
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.Unlock()
		}
	}
}
