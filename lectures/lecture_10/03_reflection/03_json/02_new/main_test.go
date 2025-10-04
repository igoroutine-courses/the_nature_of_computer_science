//go:build goexperiment.jsonv2

package main

import (
	"encoding/json/v2"
	"testing"

	"github.com/stretchr/testify/require"
)

type person struct {
	ID       int
	Name     string
	Metadata map[int]any
}

// jsontext.Encoder
// jsontext.Decoder

func BenchmarkJson(b *testing.B) {
	b.Run("marshal", func(b *testing.B) {
		b.ReportAllocs()

		p := createPerson()
		b.ResetTimer()

		for b.Loop() {
			_, err := json.Marshal(p)
			require.NoError(b, err)
		}
	})

	b.Run("unmarshal", func(b *testing.B) {
		b.ReportAllocs()

		p := createPerson()
		data, err := json.Marshal(p)
		require.NoError(b, err)

		b.ResetTimer()

		for b.Loop() {
			err = json.Unmarshal(data, p)
			require.NoError(b, err)
		}
	})
}

func createPerson() *person {
	m := make(map[int]any)
	for i := range 100 {
		m[i] = struct {
			A string
			B string
			D []int
		}{
			A: "1231231",
			B: "4313424",
			D: []int{1, 2, 3, 43, 243, 2, 53, 253, 25, 325, 23, 52, 35, 23, 523, 5, 23, 52, 35, 2},
		}
	}

	return &person{
		ID:       1,
		Name:     "test",
		Metadata: m,
	}
}
