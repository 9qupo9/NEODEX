package ws

import (
	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
)

// magicString — глобальная константа из RFC 6455 для WebSocket-рукопожатия (handshake).
const magicString = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"

// Client описывает одного подключенного по WebSocket пользователя (например, вкладка в браузере).
type Client struct {
	conn net.Conn      // Сырое TCP-соединение, перехваченное из HTTP
	hub  *Hub          // Ссылка на центральный хаб, рассылающий сообщения
	send chan []byte   // Буферизированный канал для исходящих сообщений
}

// Upgrade превращает обычный HTTP-запрос (GET) в постоянное двунаправленное WebSocket соединение.
// Выполняется строго по спецификации RFC 6455 без сторонних зависимостей.
func Upgrade(w http.ResponseWriter, r *http.Request, hub *Hub) {
	// Проверяем заголовки Upgrade
	if r.Header.Get("Upgrade") != "websocket" {
		http.Error(w, "Ожидался websocket", http.StatusBadRequest)
		return
	}

	// Читаем ключ клиента
	key := r.Header.Get("Sec-WebSocket-Key")
	
	// Хешируем ключ по стандарту: SHA1(Key + MagicString) -> Base64
	hash := sha1.New()
	hash.Write([]byte(key + magicString))
	acceptKey := base64.StdEncoding.EncodeToString(hash.Sum(nil))

	// Перехватываем управление сокетом у стандартного HTTP-сервера (Hijack)
	hj, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Webserver doesn't support hijacking", http.StatusInternalServerError)
		return
	}
	
	conn, bufrw, err := hj.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Вручную отправляем HTTP 101 Switching Protocols
	response := fmt.Sprintf("HTTP/1.1 101 Switching Protocols\r\n"+
		"Upgrade: websocket\r\n"+
		"Connection: Upgrade\r\n"+
		"Sec-WebSocket-Accept: %s\r\n\r\n", acceptKey)

	bufrw.WriteString(response)
	bufrw.Flush()

	// Создаем клиента
	client := &Client{
		conn: conn,
		hub:  hub,
		send: make(chan []byte, 256), // Буфер на 256 сообщений
	}
	client.hub.register <- client // Регистрируем в хабе

	// Запускаем две горутины на чтение и запись
	go client.writePump()
	go client.readPump(bufrw)
}

// readPump слушает входящие сообщения от клиента (пинги/понги).
// В чистом Go для полноценного WS тут нужен парсинг фреймов с маскировкой (Masking).
func (c *Client) readPump(rw *bufio.ReadWriter) {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	// Упрощенная логика: мы просто следим за тем, не закрыл ли клиент соединение.
	// Для реального production-кода (без gorilla/websocket) здесь должен быть парсер заголовков фреймов.
	for {
		_, err := rw.ReadByte()
		if err != nil {
			break // Клиент отвалился
		}
	}
}

// writePump перекладывает сообщения из внутреннего канала в TCP сокет, оборачивая в WS фреймы.
func (c *Client) writePump() {
	defer c.conn.Close()
	for msg := range c.send {
		// Оборачиваем данные в текстовый фрейм (Text Frame, FIN=1, Opcode=1)
		// 0x81 = 1000 0001
		// Поддержка сообщений < 126 байт. Для больших сообщений нужен расширенный заголовок длины.
		// TODO: Добавить корректный энкодер длины > 125 байт.
		header := []byte{0x81, byte(len(msg))}
		c.conn.Write(append(header, msg...))
	}
}
