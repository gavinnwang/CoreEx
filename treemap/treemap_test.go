package treemap

import (
	"fmt"
	"strconv"
	"testing"
)

func Test(t *testing.T) {
	treemap := New[int, string]()
	treemap.Add(1, "hello")
	fmt.Printf(treemap.Get(1))
}

func BenchmarkSeqSet(b *testing.B) {
	tr := New[int, string]()
	NumIterations := 100000
	for i := 0; i < b.N; i++ {
		for j := 0; j < NumIterations; j++ {
			tr.Add(j, "")
		}
		tr.Clear()
	}
	b.ReportAllocs()
}

func BenchmarkSeqGet(b *testing.B) {
	tr := New[int, string]()
	NumIterations := 100000
	for i := 0; i < NumIterations; i++ {
		tr.Add(i, strconv.Itoa(i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result, _ := tr.Get(i % NumIterations)
		fmt.Println(result)
	}
	b.ReportAllocs()
}
