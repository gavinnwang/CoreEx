package treemap

import "golang.org/x/exp/constraints"

type TreeMap[Key, Value any] struct {
	sentinel   *node[Key, Value]
	minNode    *node[Key, Value]
	maxNode    *node[Key, Value]
	count      int
	keyCompare func(a Key, b Key) bool
}

type node[Key, Value any] struct {
	left    *node[Key, Value]
	right   *node[Key, Value]
	parent  *node[Key, Value]
	isBlack bool
	key     Key
	value   Value
}

// https://go.dev/blog/intro-generics
// the ordered contraint describes the set of all types that can be ordered or compared with the < operator
func New[Key constraints.Ordered, Value any]() *TreeMap[Key, Value] {
	sentinel := &node[Key, Value]{isBlack: true}
	return &TreeMap[Key, Value]{
		minNode:    sentinel,
		maxNode:    sentinel,
		sentinel:   sentinel,
		keyCompare: defaultKeyCompare[Key],
	}
}

func NewWith[Key, Value any](keyCompare func(a, b Key) bool) *TreeMap[Key, Value] {
	sentinel := &node[Key, Value]{isBlack: true}
	return &TreeMap[Key, Value]{
		minNode:    sentinel,
		maxNode:    sentinel,
		sentinel:   sentinel,
		keyCompare: keyCompare,
	}
}

func (t *TreeMap[Key, Value]) Len() int { return t.count }

func defaultKeyCompare[Key constraints.Ordered](a, b Key) bool {
	return a < b
}

// Time complexity: O(log N).
// Sets the value and overrides if the key already exists
func (t *TreeMap[Key, Value]) Add(key Key, value Value) {
	parent := t.sentinel
	current := parent.left
	less := true
	// search through the tree until reaches a leaf node that is either bigger or smaller than the insert key
	for current != nil {
		parent = current
		switch {
		case t.keyCompare(key, current.key):
			current = current.left
			less = true
		case t.keyCompare(current.key, key):
			current = current.right
			less = false
		default:
			current.value = value
			// silently overrides the value
			return
		}
	}
	x := &node[Key, Value]{
		parent: parent,
		value:  value,
		key:    key,
	}
	if less {
		parent.left = x
	} else {
		parent.right = x
	}
	if t.minNode.left != nil {
		t.minNode = t.minNode.left
	}
	if t.maxNode == t.sentinel {
		t.maxNode = t.sentinel.left
	} else if t.maxNode.right != nil {
		t.maxNode = t.maxNode.right
	}
	t.addAndRebalance(x)
	t.count++
}

func (t *TreeMap[Key, Value]) addAndRebalance(x *node[Key, Value]) {
	root := t.sentinel.left
	// if x is root then it must be black else it will be red
	x.isBlack = x == root
	for x != root && !x.parent.isBlack {
		// if x's parent is red then x's parent definitely has a parent that is black
		if x.parent == x.parent.parent.left {
			y := x.parent.parent.right
			if y != nil && !y.isBlack {
				x = x.parent
				x.isBlack = true
				x = x.parent
				x.isBlack = x == root
				y.isBlack = true
			} else {
				if x != x.parent.left {
					x = x.parent
					rotateLeft(x)
				}
				x = x.parent
				x.isBlack = true
				x = x.parent
				x.isBlack = false
				rotateRight(x)
				break
			}
		} else {
			y := x.parent.parent.left
			if y != nil && !y.isBlack {
				x = x.parent
				x.isBlack = true
				x = x.parent
				x.isBlack = x == root
				y.isBlack = true
			} else {
				if x == x.parent.left {
					x = x.parent
					rotateRight(x)
				}
				x = x.parent
				x.isBlack = true
				x = x.parent
				x.isBlack = false
				rotateLeft(x)
				break
			}
		}
	}
}

func rotateLeft[Key, Value any](x *node[Key, Value]) {
	y := x.right
	x.right = y.left
	if x.right != nil {
		x.right.parent = x
	}
	y.parent = x.parent
	if x == x.parent.left {
		x.parent.left = y
	} else {
		x.parent.right = y
	}
	y.left = x
	x.parent = y
}

func rotateRight[Key, Value any](x *node[Key, Value]) {
	y := x.left
	x.right = y.right
	if x.left != nil {
		x.left.parent = x
	}
	y.parent = x.parent
	if x == x.parent.left {
		x.parent.left = y
	} else {
		x.parent.right = y
	}
	x.right = x
	x.parent = y
}

// Time complexity: O(log N)
func (t *TreeMap[Key, Value]) Remove(key Key) bool {
	z := t.findNode(key)
	if z == nil {
		return false
	}
	if t.maxNode == z {
		if z.left != nil {
			t.maxNode = z.left
		} else {
			t.maxNode = z.parent
		}
	}
	if t.minNode == z {
		if z.right != nil {
			t.minNode = z.right
		} else {
			t.minNode = z.parent
		}
	}
	t.count--
	removeAndRebalance(t.sentinel.left, z)
	return true
}

func removeAndRebalance[Key, Value any](
	root, z *node[Key, Value],
) {
	var y *node[Key, Value]
	if z.left == nil || z.right == nil {
		y = z
	} else {
		y = successor(z)
	}
	var x *node[Key, Value]
	if y.left != nil {
		x = y.left
	} else {
		x = y.right
	}
	var w *node[Key, Value]
	if x != nil {
		x.parent = y.parent
	}
	if y == y.parent.left {
		y.parent.left = x
		if y != root {
			w = y.parent.right
		} else {
			root = x // w == nil
		}
	} else {
		y.parent.right = x
		w = y.parent.left
	}
	removedBlack := y.isBlack
	if y != z {
		y.parent = z.parent
		if z == z.parent.left {
			y.parent.left = y
		} else {
			y.parent.right = y
		}
		y.left = z.left
		y.left.parent = y
		y.right = z.right
		if y.right != nil {
			y.right.parent = y
		}
		y.isBlack = z.isBlack
		if root == z {
			root = y
		}
	}
	if removedBlack && root != nil {
		if x != nil {
			x.isBlack = true
		} else {
			for {
				if w != w.parent.left {
					if !w.isBlack {
						w.isBlack = true
						w.parent.isBlack = false
						rotateLeft(w.parent)
						if root == w.left {
							root = w
						}
						w = w.left.right
					}
					if (w.left == nil || w.left.isBlack) && (w.right == nil || w.right.isBlack) {
						w.isBlack = false
						x = w.parent
						if x == root || !x.isBlack {
							x.isBlack = true
							break
						}
						if x == x.parent.left {
							w = x.parent.right
						} else {
							w = x.parent.left
						}
					} else {
						if w.right == nil || w.right.isBlack {
							w.left.isBlack = true
							w.isBlack = false
							rotateRight(w)
							w = w.parent
						}
						w.isBlack = w.parent.isBlack
						w.parent.isBlack = true
						w.right.isBlack = true
						rotateLeft(w.parent)
						break
					}
				} else {
					if !w.isBlack {
						w.isBlack = true
						w.parent.isBlack = false
						rotateRight(w.parent)
						if root == w.right {
							root = w
						}
						w = w.right.left
					}
					if (w.left == nil || w.left.isBlack) && (w.right == nil || w.right.isBlack) {
						w.isBlack = false
						x = w.parent
						if !x.isBlack || x == root {
							x.isBlack = true
							break
						}
						if x == x.parent.left {
							w = x.parent.right
						} else {
							w = x.parent.left
						}
					} else {
						if w.left == nil || w.left.isBlack {
							w.right.isBlack = true
							w.isBlack = false
							rotateLeft(w)
							w = w.parent
						}
						w.isBlack = w.parent.isBlack
						w.parent.isBlack = true
						w.left.isBlack = true
						rotateRight(w.parent)
						break
					}
				}
			}
		}
	}
}

func successor[Key, Value any](x *node[Key, Value]) *node[Key, Value] {
	if x.right != nil {
		return mostLeft(x.right)
	}
	for x != x.parent.left {
		x = x.parent
	}
	return x.parent
}

func mostLeft[Key, Value any](
	x *node[Key, Value],
) *node[Key, Value] {
	for x.left != nil {
		x = x.left
	}
	return x
}

func (t *TreeMap[Key, Value]) findNode(key Key) *node[Key, Value] {
	current := t.sentinel.left
	for current != nil {
		switch {
		case t.keyCompare(key, current.key):
			current = current.left
		case t.keyCompare(current.key, key):
			current = current.right
		default:
			return current
		}
	}
	return nil
}

// Get retrieves a value from a map for specified key and reports if it exists.
// Complexity: O(log N).
func (t *TreeMap[Key, Value]) Get(key Key) (Value, bool) {
	node := t.findNode(key)
	if node == nil {
		node = t.sentinel
	}
	return node.value, node != t.sentinel
}

// Clear clears the map.
// Complexity: O(1).
func (t *TreeMap[Key, Value]) Clear() {
	t.count = 0
	t.minNode = t.sentinel
	t.maxNode = t.sentinel
	t.sentinel.left = nil
}

func (t *TreeMap[Key, Value]) GetMin() (Value, bool) {
	if t.minNode == t.sentinel {
		return *new(Value), false
	}
	return t.minNode.value, true
}

func (t *TreeMap[Key, Value]) GetMax() (Value, bool) {
	if t.maxNode == t.sentinel {
		return *new(Value), false
	}
	return t.maxNode.value, true
}

type ForwardIterator[Key, Value any] struct {
	tree *TreeMap[Key, Value]
	node *node[Key, Value]
}

func (i ForwardIterator[Key, Value]) Key() Key { return i.node.key }

func (i ForwardIterator[Key, Value]) Value() Value { return i.node.value }

// Iterator returns an iterator for tree map.
// It starts at the first element and goes to the one-past-the-end position.
// You can iterate a map at O(N) complexity.

// Method complexity: O(1)
func (t *TreeMap[Key, Value]) Iterator() ForwardIterator[Key, Value] {
	return ForwardIterator[Key, Value]{tree: t, node: t.minNode}
}

// Valid reports if the iterator position is valid.
// In other words it returns true if an iterator is not at the one-before-the-start position.
func (i ForwardIterator[Key, Value]) Valid() bool { return i.node != i.tree.sentinel }

func (i *ForwardIterator[Key, Value]) Next() {
	if i.node == i.tree.sentinel {
		panic("out of bound iteration")
	}
	i.node = successor(i.node)
}
