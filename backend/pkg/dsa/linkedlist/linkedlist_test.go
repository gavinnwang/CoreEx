package list

import (
	"testing"
)

func Test(t *testing.T) {
	list := New[int]()
	list.PushBack(1)
	list.PushBack(2)
	for i := 0; i < 9999; i++ {
		if i%2 == 0 {
			list.PushBack(i)
		} else {
			list.PushFront(i)
		}
	}
	for list.Len() != 0 {
		list.Remove(list.Front())
	}
}

