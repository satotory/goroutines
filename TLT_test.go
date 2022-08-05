package main

import "testing"

func BenchmarkTLT10(b *testing.B) {
	for i := 0; i < b.N; i++ {
		TLT()
	}
}
