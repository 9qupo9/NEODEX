package tcp

import (
	"encoding/binary"
	"io"
)

// Кастомный бинарный протокол для P2P-взаимодействия и HFT-ботов.
// HTTP/JSON слишком медленные для маркет-мейкеров. Бинарный TCP протокол
// экономит такты процессора на парсинг и уменьшает размер пакета в сети.
//
// Формат пакета:
// [4 байта: Длина пакета (uint32, BigEndian)]
// [1 байт: Тип сообщения (MsgType)]
// [N байт: Полезная нагрузка (Payload, например protobuf или raw bytes)]

const (
	MsgTypePing  byte = 0x01 // Пинг для поддержания соединения (Keep-Alive)
	MsgTypeOrder byte = 0x02 // Отправка нового ордера
	MsgTypeTrade byte = 0x03 // Трансляция совершенной сделки
)

// Message — абстрактная структура распарсенного бинарного сообщения.
type Message struct {
	Type    byte
	Payload []byte
}

// ReadMsg читает из TCP-сокета (или любого io.Reader) ровно одно сообщение.
// Блокируется до тех пор, пока не прочитает целиком длину, тип и payload.
func ReadMsg(r io.Reader) (*Message, error) {
	// 1. Читаем длину (4 байта)
	var length uint32
	if err := binary.Read(r, binary.BigEndian, &length); err != nil {
		return nil, err
	}
	
	// 2. Читаем тип сообщения (1 байт)
	msgType := make([]byte, 1)
	if _, err := io.ReadFull(r, msgType); err != nil {
		return nil, err
	}
	
	// 3. Читаем остаток пакета (payload)
	// Длина payload = Общая длина - 1 байт (так как 1 байт ушел на msgType)
	payload := make([]byte, length-1)
	if _, err := io.ReadFull(r, payload); err != nil {
		return nil, err
	}
	
	return &Message{Type: msgType[0], Payload: payload}, nil
}

// WriteMsg упаковывает данные в бинарный формат и отправляет в TCP-сокет.
func WriteMsg(w io.Writer, msgType byte, payload []byte) error {
	// Длина = 1 байт типа + длина самого пейлоада
	length := uint32(1 + len(payload))
	
	// 1. Записываем длину
	if err := binary.Write(w, binary.BigEndian, length); err != nil {
		return err
	}
	
	// 2. Записываем тип
	if _, err := w.Write([]byte{msgType}); err != nil {
		return err
	}
	
	// 3. Записываем пейлоад
	_, err := w.Write(payload)
	return err
}
