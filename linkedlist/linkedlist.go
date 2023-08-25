package list

import (
	"fmt"
	"strings"
)

type List[Value any] struct {
	root *Node[Value]
	len  int
}

type Node[Value any] struct {
	next, prev *Node[Value]
	list       *List[Value]
	Value      Value
}

func New[Value any]() *List[Value] { return new(List[Value]).init() }

// Init is a convient helper function that initializes the sentinel root and len
func (l *List[Value]) init() *List[Value] {
	l.root = &Node[Value]{}
	l.root.next = l.root
	l.root.prev = l.root
	l.len = 0
	return l
}

func (n *Node[Value]) Next() *Node[Value] {
	if p := n.next; n.list != nil && p != n.list.root {
		return p
	}
	return nil
}

func (n *Node[Value]) Prev() *Node[Value] {
	if p := n.prev; n.list != nil && p != n.list.root {
		return p
	}
	return nil
}

func (l *List[Value]) Len() int { return l.len }

func (l *List[Value]) Front() *Node[Value] {
	if l.len == 0 {
		return nil
	}
	return l.root.next
}

func (l *List[Value]) Back() *Node[Value] {
	if l.len == 0 {
		return nil
	}
	return l.root.prev
}

// insert inserts n after at increments l and return n
func (l *List[Value]) insert(n, at *Node[Value]) *Node[Value] {
	n.prev = at
	n.next = at.next
	n.prev.next = n
	n.next.prev = n
	n.list = l
	l.len++
	return n
}

// insertValue is a wrapper for calling insert
func (l *List[Value]) insertValue(v Value, at *Node[Value]) *Node[Value] {
	return l.insert(&Node[Value]{Value: v}, at)
}

func (l *List[Value]) remove(n *Node[Value]) {
	n.prev.next = n.next
	n.next.prev = n.prev
	n.next = nil // avoid memory leaks
	n.prev = nil // avoid memory leaks
	n.list = nil
	l.len--
}

// Remove removes n from l if n is an element and returns the value e.Value
func (l *List[Value]) Remove(n *Node[Value]) Value {
	if n.list == l {
		l.remove(n)
	}
	return n.Value
}

// move moves
func (l *List[Value]) move(n, at *Node[Value]) {
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

func (l *List[Value]) PushFront(v Value) *Node[Value] {
	return l.insertValue(v, l.root)
}

func (l *List[Value]) PushBack(v Value) *Node[Value] {
	return l.insertValue(v, l.root.prev)
}

// String implements the stringer interface to print list
func (l *List[Value]) String() string {
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
