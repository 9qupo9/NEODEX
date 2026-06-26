package tcp

import (
	"dex/internal/domain"
	"dex/internal/engine"
	"dex/pkg/decimal"
	"encoding/binary"
	"log"
	"net"
)

// Handler отвечает за бизнес-логику обработки входящих TCP соединений.
type Handler struct {
	Router *engine.EngineRouter
}

// NewHandler создает новый обработчик TCP.
func NewHandler(r *engine.EngineRouter) *Handler {
	return &Handler{Router: r}
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
			// Десериализация бинарного payload
			if len(msg.Payload) != 66 {
				log.Printf("TCP: Неверная длина payload для ордера: %d", len(msg.Payload))
				_ = WriteMsg(conn, MsgTypeOrder, []byte{0x00}) // 0x00 - Error
				continue
			}

			// Парсинг AccountID (32 bytes)
			accBytes := msg.Payload[0:32]
			accLen := 0
			for i, b := range accBytes {
				if b == 0 {
					accLen = i
					break
				}
				accLen = 32
			}
			accountID := string(accBytes[:accLen])

			// Парсинг активов (по 8 байт)
			baseBytes := msg.Payload[32:40]
			baseLen := 0
			for i, b := range baseBytes {
				if b == 0 {
					baseLen = i
					break
				}
				baseLen = 8
			}
			base := string(baseBytes[:baseLen])

			quoteBytes := msg.Payload[40:48]
			quoteLen := 0
			for i, b := range quoteBytes {
				if b == 0 {
					quoteLen = i
					break
				}
				quoteLen = 8
			}
			quote := string(quoteBytes[:quoteLen])

			side := domain.Buy
			if msg.Payload[48] == 1 {
				side = domain.Sell
			}

			orderType := domain.Limit
			if msg.Payload[49] == 1 {
				orderType = domain.Market
			}

			// Парсинг цены и количества (uint64)
			priceRaw := binary.BigEndian.Uint64(msg.Payload[50:58])
			qtyRaw := binary.BigEndian.Uint64(msg.Payload[58:66])

			// Конвертация в decimal (делим на 10^8)
			priceDec := decimal.NewFromInt(int64(priceRaw)).Div(decimal.NewFromInt(100000000))
			qtyDec := decimal.NewFromInt(int64(qtyRaw)).Div(decimal.NewFromInt(100000000))

			pair := domain.Pair{BaseAsset: base, QuoteAsset: quote}
			order := domain.NewOrder("tcp-"+accountID[:4], accountID, pair, side, orderType, priceDec, qtyDec)

			// Передаем в роутер
			trades, err := h.Router.PlaceOrder(order)
			if err != nil {
				_ = WriteMsg(conn, MsgTypeOrder, []byte{0x00}) // 0x00 - Error
			} else {
				// 0x01 - Success, возвращаем количество сделок как 1 байт
				ack := []byte{0x01, byte(len(trades))}
				_ = WriteMsg(conn, MsgTypeOrder, ack)
			}
		}
	}
}
