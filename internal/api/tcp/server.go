package tcp

import (
	"log"
	"net"
)

// Server — это низкоуровневый TCP-сервер.
// Он необходим для работы с маркет-мейкерами и P2P-взаимодействия других нод (если мы строим распределенный DEX).
type Server struct {
	addr    string   // Порт и хост для прослушивания (например, ":9000")
	handler *Handler // Ссылка на обработчик бизнес-логики
}

// NewServer инициализирует структуру TCP сервера.
func NewServer(addr string, handler *Handler) *Server {
	return &Server{
		addr:    addr,
		handler: handler,
	}
}

// Start запускает прослушивание порта. Блокирующий вызов.
func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	log.Printf("TCP (P2P/HFT) сервер запущен на %s", s.addr)

	// Цикл приема новых подключений
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Ошибка при принятии TCP соединения: %v", err)
			continue
		}
		
		// Передаем подключение в обработчик, запуская отдельную горутину на каждого клиента.
		go s.handler.HandleConnection(conn)
	}
}
