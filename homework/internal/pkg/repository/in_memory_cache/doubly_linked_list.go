package in_memory_cache

import (
	"fmt"
	"homework/internal/pkg/repository"
)

type Node struct {
	key   int64
	value repository.PickupPoint
	prev  *Node
	next  *Node
}

type DoublyLinkedList struct {
	len  int
	tail *Node
	head *Node
}

func NewDoublyList() *DoublyLinkedList {
	return &DoublyLinkedList{}
}

func (d *DoublyLinkedList) AddToFront(key int64, value repository.PickupPoint) {
	newNode := &Node{
		key:   key,
		value: value,
	}
	if d.head == nil {
		d.head = newNode
		d.tail = newNode
	} else {
		newNode.next = d.head
		d.head.prev = newNode
		d.head = newNode
	}
	d.len++
	return
}

func (d *DoublyLinkedList) RemoveFromFront() {
	if d.head == nil {
		return
	} else if d.head == d.tail {
		d.head = nil
		d.tail = nil
	} else {
		d.head = d.head.next
	}
	d.len--
}

func (d *DoublyLinkedList) AddToEnd(node *Node) {
	newNode := node
	if d.head == nil {
		d.head = newNode
		d.tail = newNode
	} else {
		currentNode := d.head
		for currentNode.next != nil {
			currentNode = currentNode.next
		}
		newNode.prev = currentNode
		currentNode.next = newNode
		d.tail = newNode
	}
	d.len++
}
func (d *DoublyLinkedList) Front() *Node {
	return d.head
}

func (d *DoublyLinkedList) MoveNodeToEnd(node *Node) {
	prev := node.prev
	next := node.next

	if prev != nil {
		prev.next = next
	}

	if next != nil {
		next.prev = prev
	}
	if d.tail == node {
		d.tail = prev
	}
	if d.head == node {
		d.head = next
	}
	node.next = nil
	node.prev = nil
	d.len--
	d.AddToEnd(node)
}

func (d *DoublyLinkedList) TraverseForward() error {
	if d.head == nil {
		return fmt.Errorf("TraverseError: List is empty")
	}
	temp := d.head
	for temp != nil {
		fmt.Printf("key = %v, value = %v, prev = %v, next = %v\n", temp.key, temp.value, temp.prev, temp.next)
		temp = temp.next
	}
	fmt.Println()
	return nil
}

func (d *DoublyLinkedList) Size() int {
	return d.len
}
