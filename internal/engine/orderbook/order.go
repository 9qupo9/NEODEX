package orderbook

import "dex/internal/domain"

// OrderNode — это узел интрузивного двусвязного списка (Intrusive Doubly-Linked List).
// Оборачивает domain.Order и добавляет указатели на следующий и предыдущий элементы.
// Использование двусвязного списка критически важно, так как позволяет удалять отмененные ордера за время O(1),
// не сдвигая элементы массива в памяти (как было бы со срезами - slice).
type OrderNode struct {
	Order *domain.Order
	Next  *OrderNode
	Prev  *OrderNode
}

// OrderQueue — очередь ордеров на конкретном ценовом уровне.
// Работает по принципу FIFO (First In, First Out).
type OrderQueue struct {
	Head *OrderNode // Первый ордер в очереди (тот, что стоит дольше всех)
	Tail *OrderNode // Последний добавленный ордер
	Size int        // Количество ордеров в очереди для быстрого получения длины O(1)
}

// NewOrderQueue инициализирует пустую очередь.
func NewOrderQueue() *OrderQueue {
	return &OrderQueue{}
}

// Push добавляет новый ордер в конец очереди (Tail).
// Возвращает указатель на созданный узел для быстрого поиска и удаления в будущем.
func (q *OrderQueue) Push(order *domain.Order) *OrderNode {
	node := &OrderNode{Order: order}
	if q.Head == nil {
		q.Head = node
		q.Tail = node
	} else {
		q.Tail.Next = node
		node.Prev = q.Tail
		q.Tail = node
	}
	q.Size++
	return node
}

// Remove удаляет конкретный узел из очереди за O(1), меняя указатели соседних узлов.
func (q *OrderQueue) Remove(node *OrderNode) {
	if node.Prev != nil {
		node.Prev.Next = node.Next
	} else {
		q.Head = node.Next
	}

	if node.Next != nil {
		node.Next.Prev = node.Prev
	} else {
		q.Tail = node.Prev
	}

	// Отвязываем узел для сборщика мусора (GC)
	node.Next = nil
	node.Prev = nil
	q.Size--
}

// Pop извлекает и удаляет первый (Head) ордер из очереди.
// Вызывается движком сведения (Matcher), когда цена достигает этого уровня.
func (q *OrderQueue) Pop() *domain.Order {
	if q.Head == nil {
		return nil
	}
	node := q.Head
	q.Remove(node)
	return node.Order
}
