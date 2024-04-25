package in_memory_cache

import (
	"grpc/internal/pkg/repository"
)

type Lru struct {
	LL *DoublyLinkedList
}

func NewLRU() *Lru {
	return &Lru{
		LL: NewDoublyList(),
	}
}

func (l *Lru) evict() int64 {
	key := l.LL.Front().key
	l.LL.RemoveFromFront()
	return key
}

func (l *Lru) get(node *Node) {
	l.LL.MoveNodeToEnd(node)
}

func (l *Lru) set(node *Node) {
	l.LL.AddToEnd(node)
}

func (l *Lru) set_overwrite(node *Node, value repository.PickupPoint) {
	node.value = value
	l.LL.MoveNodeToEnd(node)
}
