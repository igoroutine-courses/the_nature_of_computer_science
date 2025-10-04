package _1_small_objects

import (
	"testing"
)

type SmallObject struct {
	data [128]byte
}

func allocSmallObjects(n int) []*SmallObject {
	out := make([]*SmallObject, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, &SmallObject{})
	}
	return out
}

var result []*SmallObject

func BenchmarkAllocSmall(b *testing.B) {
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			result = allocSmallObjects(1_00)
		}
	})
}
