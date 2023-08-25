package list

import (
	"fmt"
	"strings"
)

type List[value any] struct {
	root *Node[value]
	len  int
}

type Node[value any] struct {
	next, prev *Node[value]
	list       *List[value]
	Value      value
}

func New[value any]() *List[value] { return new(List[value]).init() }

// Init is a convient helper function that initializes the sentinel root and len
func (l *List[value]) init() *List[value] {
	l.root = &Node[value]{}
	l.root.next = l.root
	l.root.prev = l.root
	l.len = 0
	return l
}

func (n *Node[value]) Next() *Node[value] {
	if p := n.next; n.list != nil && p != n.list.root {
		return p
	}
	return nil
}

func (n *Node[value]) Prev() *Node[value] {
	if p := n.prev; n.list != nil && p != n.list.root {
		return p
	}
	return nil
}

func (l *List[value]) Len() int { return l.len }

func (l *List[value]) Front() *Node[value] {
	if l.len == 0 {
		return nil
	}
	return l.root.next
}

func (l *List[value]) Back() *Node[value] {
	if l.len == 0 {
		return nil
	}
	return l.root.prev
}

// insert inserts n after at increments l and return n
func (l *List[value]) insert(n, at *Node[value]) *Node[value] {
	n.prev = at
	n.next = at.next
	n.prev.next = n
	n.next.prev = n
	n.list = l
	l.len++
	return n
}

// insertValue is a wrapper for calling insert
func (l *List[value]) insertValue(v value, at *Node[value]) *Node[value] {
	return l.insert(&Node[value]{Value: v}, at)
}

func (l *List[value]) remove(n *Node[value]) {
	n.prev.next = n.next
	n.next.prev = n.prev
	n.next = nil // avoid memory leaks
	n.prev = nil // avoid memory leaks
	n.list = nil
	l.len--
}

// Remove removes n from l if n is an element and returns the value e.Value
func (l *List[value]) Remove(n *Node[value]) value {
	if n.list == l {
		l.remove(n)
	}
	return n.Value
}

// move moves
func (l *List[value]) move(n, at *Node[value]) {
	if n == at {
		return
	}
	n.prev.next = n.next
	n.next.prev = n.prev

	n.prev = at
	n.next = at.next
	n.prev.next = n
	n.next.prev = n
}

func (l *List[value]) PushFront(v value) *Node[value] {
	return l.insertValue(v, l.root)
}

func (l *List[value]) PushBack(v value) *Node[value] {
	return l.insertValue(v, l.root.prev)
}

// String implements the stringer interface to print list
func (l *List[value]) String() string {
	var sb strings.Builder
	sb.WriteString("[")

	current := l.root.next
	if current != l.root {
		sb.WriteString(fmt.Sprint(current.Value))
		current = current.next
	}

	for current != l.root {
		sb.WriteString(", ")
		sb.WriteString(fmt.Sprint(current.Value))
		current = current.next
	}

	sb.WriteString("]")

	return sb.String()
}
