package main

import "testing"

func BenchmarkRun(b *testing.B) {
	b.SetBytes(2)
	run()
}

func TestRun(t *testing.T) {
	run()
}
