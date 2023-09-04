package treemap

import (
	"fmt"
	"strconv"
	"testing"
)

func Test(t *testing.T) {
	treemap := New[int, string]()
	treemap.Put(1, "hello")
	fmt.Printf(treemap.Get(1))
}

func BenchmarkSeqSet(b *testing.B) {
	tr := New[int, string]()
	NumIterations := 100000
	for i := 0; i < b.N; i++ {
		for j := 0; j < NumIterations; j++ {
			tr.Put(j, "")
		}
	}
	b.ReportAllocs()
}

func BenchmarkSeqGet(b *testing.B) {
	tr := New[int, string]()
	NumIterations := 100
	for i := 0; i < NumIterations; i++ {
		tr.Put(i, strconv.Itoa(i))
	}
	b.ResetTimer()
	for i := 0; i < NumIterations; i++ {
		result, _ := tr.Get(i % NumIterations)
		fmt.Println(result)
	}
	b.ReportAllocs()
}

func TestGetMax(t *testing.T) {
	tr := New[int, string]()
	tr.Put(1, "1")
	tr.Put(2, "2")
	max, found := tr.GetMax()
	if !found {
		t.Error("max not found")
	}
	fmt.Printf("max: %v\n", max)
	tr.Clear()
	max, found = tr.GetMax()
	if !found {
		fmt.Println("max not found")
	}
	fmt.Printf("max: %v\n", max)
	tr.Put(1, "1")
	tr.Remove(1)
	tr.Put(4, "4")
	tr.Put(2, "2")
	max, found = tr.GetMax()
	if !found {
		fmt.Println("max not found")
	}
	fmt.Printf("max: %v\n", max)
	

}
