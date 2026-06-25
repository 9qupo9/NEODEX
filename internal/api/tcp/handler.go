package tcp

import (
	"dex/internal/engine"
	"log"
	"net"
)

// Handler отвечает за бизнес-логику обработки входящих TCP соединений.
type Handler struct {
	Engine *engine.Engine // Ссылка на торговое ядро (очевидно, в реальной бирже тут будет роутер на несколько Engine для разных пар)
}

// NewHandler создает новый обработчик TCP.
func NewHandler(e *engine.Engine) *Handler {
	return &Handler{Engine: e}
}

// HandleConnection вызывается для каждого нового подключившегося TCP-клиента (например, HFT-бота).
// Работает в отдельной горутине для каждого клиента.
func (h *Handler) HandleConnection(conn net.Conn) {
	defer conn.Close()
	log.Printf("TCP Клиент подключился: %s", conn.RemoteAddr())

	// Бесконечный цикл чтения (Event Loop клиента)
	for {
		msg, err := ReadMsg(conn)
		if err != nil {
			log.Printf("TCP Клиент отключился: %s", conn.RemoteAddr())
			return
		}

		// Маршрутизация по типу бинарного сообщения
		switch msg.Type {
		case MsgTypePing:
			// Отвечаем Pong (эхо) для поддержания соединения
			_ = WriteMsg(conn, MsgTypePing, []byte("pong"))
			
		case MsgTypeOrder:
			// TODO: Распарсить msg.Payload (например, из Protobuf в domain.Order).
			// В рамках прототипа мы просто шлем подтверждение получения (ACK - Acknowledgment).
			
			// trade, err := h.Engine.PlaceOrder(parsedOrder)
			// if err != nil { ... }
			
			_ = WriteMsg(conn, MsgTypeOrder, []byte("ack"))
		}
	}
}
