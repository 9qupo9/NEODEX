package main

import (
	"encoding/binary"
	"log"
	"net"
	"time"
)

// MsgType определяет тип бинарного сообщения
const (
	MsgTypePing  byte = 0x01
	MsgTypeOrder byte = 0x02
)

// MarketMakerBot — это простой HFT бот, который подключается к бирже по TCP
// и периодически выставляет ордера для поддержания ликвидности (Market Making).
func main() {
	log.Println("Запуск Market Maker бота...")

	// Подключаемся к локальному HFT TCP серверу
	conn, err := net.Dial("tcp", "127.0.0.1:9000")
	if err != nil {
		log.Fatalf("Не удалось подключиться к TCP серверу: %v", err)
	}
	defer conn.Close()
	log.Println("Успешно подключено к бирже. Начинаем маркет-мейкинг.")

	accountID := "mm-bot-1000"
	base := "BTC"
	quote := "USDT"

	// Бесконечный цикл: каждые 2 секунды ставим Bid и Ask вокруг некой цены
	centerPrice := float64(64000.0)

	for {
		// 1. Формируем BID ордер (покупка ниже текущей цены)
		bidPrice := uint64((centerPrice - 10.0) * 100000000) // Цена * 10^8
		bidQty := uint64(0.01 * 100000000)
		sendOrder(conn, accountID, base, quote, 0, 0, bidPrice, bidQty) // 0 = Buy, 0 = Limit

		// 2. Формируем ASK ордер (продажа выше текущей цены)
		askPrice := uint64((centerPrice + 10.0) * 100000000) // Цена * 10^8
		askQty := uint64(0.01 * 100000000)
		sendOrder(conn, accountID, base, quote, 1, 0, askPrice, askQty) // 1 = Sell, 0 = Limit

		// Двигаем цену случайным образом для симуляции рынка
		centerPrice += float64(time.Now().UnixNano()%20) - 10.0

		// Читаем ответ (ACK) от сервера
		buf := make([]byte, 6)
		conn.Read(buf)

		log.Printf("Ордера выставлены вокруг %.2f USDT", centerPrice)
		time.Sleep(2 * time.Second)
	}
}

// sendOrder упаковывает ордер в бинарный формат (66 байт) и отправляет на сервер
func sendOrder(conn net.Conn, accountID, base, quote string, side, orderType byte, price, qty uint64) {
	payload := make([]byte, 66)

	copy(payload[0:32], []byte(accountID))
	copy(payload[32:40], []byte(base))
	copy(payload[40:48], []byte(quote))

	payload[48] = side
	payload[49] = orderType

	binary.BigEndian.PutUint64(payload[50:58], price)
	binary.BigEndian.PutUint64(payload[58:66], qty)

	// Формат фрейма: [Len: 4 bytes][Type: 1 byte][Payload]
	msgLen := uint32(1 + len(payload))
	header := make([]byte, 5)
	binary.BigEndian.PutUint32(header[0:4], msgLen)
	header[4] = MsgTypeOrder

	conn.Write(append(header, payload...))
}
