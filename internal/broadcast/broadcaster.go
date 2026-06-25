package broadcast

import "sync"

// Event представляет абстрактное сообщение, которое рассылается по всей бирже.
// Может содержать Trade (Сделка), OrderbookUpdate (Изменение стакана) и т.д.
type Event struct {
	Topic   string      // Название топика (канала), например: "TRADE_BTC_USDT"
	Payload interface{} // Сами данные (JSON, структура)
}

// Broadcaster — это шина событий (Event Bus) основанная на паттерне Publish-Subscribe.
// Выступает в роли "нервной системы" биржи. Движок публикует сюда изменения,
// а WebSocket-серверы подписываются и отправляют изменения клиентам в браузеры.
// Заменяет сторонние брокеры сообщений (Redis PubSub, Kafka) в рамках чистой Go архитектуры.
type Broadcaster struct {
	mu          sync.RWMutex
	subscribers map[string][]chan Event
}

// NewBroadcaster создает новую шину событий.
func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		subscribers: make(map[string][]chan Event),
	}
}

// Subscribe подписывает канал (channel) на определенный топик.
// Вызывается, когда клиент подключается к WebSocket и просит рыночные данные.
func (b *Broadcaster) Subscribe(topic string, ch chan Event) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subscribers[topic] = append(b.subscribers[topic], ch)
}

// Unsubscribe отписывает канал от топика.
// Вызывается при дисконнекте WebSocket клиента для избежания утечек памяти.
func (b *Broadcaster) Unsubscribe(topic string, ch chan Event) {
	b.mu.Lock()
	defer b.mu.Unlock()
	subs := b.subscribers[topic]
	for i, sub := range subs {
		if sub == ch {
			b.subscribers[topic] = append(subs[:i], subs[i+1:]...)
			break
		}
	}
}

// Publish отправляет событие всем подписчикам топика асинхронно.
// Использует неблокирующий select, чтобы один "зависший" клиент 
// не подвесил всю шину и весь торговый движок.
func (b *Broadcaster) Publish(topic string, payload interface{}) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	subs := b.subscribers[topic]
	for _, sub := range subs {
		// Неблокирующая отправка
		select {
		case sub <- Event{Topic: topic, Payload: payload}:
		default:
			// Если буфер канала клиента переполнен (клиент не успевает читать), 
			// мы просто дропаем пакет для этого конкретного клиента, 
			// чтобы не заблокировать Publisher'а (саму биржу).
		}
	}
}
