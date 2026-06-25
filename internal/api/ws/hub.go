package ws

import (
	"crypto/rand"
	"encoding/binary"
	"sync"
)

// Shard выступает в роли "мини-хаба" для группы WebSocket подключений.
// Он управляет только своей частью клиентов (партицией) и безопасно рассылает им данные.
type Shard struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.Mutex
}

func newShard() *Shard {
	return &Shard{
		broadcast:  make(chan []byte, 1024), // Буфер на 1024 пакета для защиты от блокировок
		register:   make(chan *Client, 128),
		unregister: make(chan *Client, 128),
		clients:    make(map[*Client]bool),
	}
}

// Run — бесконечный цикл (Event Loop) мультиплексирования для одного шарда.
func (s *Shard) Run() {
	for {
		select {
		case client := <-s.register:
			s.mu.Lock()
			s.clients[client] = true
			s.mu.Unlock()
			
		case client := <-s.unregister:
			s.mu.Lock()
			if _, ok := s.clients[client]; ok {
				delete(s.clients, client)
				close(client.send)
			}
			s.mu.Unlock()
			
		case message := <-s.broadcast:
			s.mu.Lock()
			for client := range s.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(s.clients, client)
				}
			}
			s.mu.Unlock()
		}
	}
}

// ShardedHub — главный диспетчер, состоящий из множества Шардов.
// Решает проблему единого узкого горлышка при рассылке на миллионы соединений.
type ShardedHub struct {
	shards     []*Shard
	shardCount uint32
}

// NewShardedHub создает хаб, разбитый на shardCount независимых партиций.
// Рекомендуется от 64 до 256 шардов для максимальной утилизации CPU.
func NewShardedHub(shardCount uint32) *ShardedHub {
	sh := &ShardedHub{
		shards:     make([]*Shard, shardCount),
		shardCount: shardCount,
	}

	for i := uint32(0); i < shardCount; i++ {
		sh.shards[i] = newShard()
	}

	return sh
}

// Run запускает Event Loop для каждого шарда в своей выделенной горутине.
func (sh *ShardedHub) Run() {
	for _, shard := range sh.shards {
		go shard.Run()
	}
}

// Broadcast рассылает (fan-out) входящее сообщение сразу во все шарды параллельно.
// Главный Event Loop больше не ждет обхода всех миллионов клиентов.
func (sh *ShardedHub) Broadcast(msg []byte) {
	for _, shard := range sh.shards {
		// Отправляем сообщение в канал шарда.
		// Использование буферизированных каналов защищает от блокировки, 
		// если один из шардов временно занят.
		shard.broadcast <- msg
	}
}

// GetShard возвращает конкретный шард для подключения нового клиента.
// Для распределения используется случайный Round-Robin, чтобы шарды наполнялись равномерно.
func (sh *ShardedHub) GetShard() *Shard {
	var b [4]byte
	rand.Read(b[:])
	rnd := binary.LittleEndian.Uint32(b[:])
	return sh.shards[rnd%sh.shardCount]
}
