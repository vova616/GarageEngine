package engine

import "testing"

const benchLoops = 1000000

func BenchmarkFuncList(b *testing.B) {
	x := 0
	f := func() { x++ }
	fnc := NewFuncList(nil)
	for i := 0; i < benchLoops; i++ {
		fnc.Add(f)
	}
	fnc.Run()
	if x != benchLoops {
		b.Fatalf("Got %d need %d\n", x, benchLoops)
	}
}

func TestFuncList(t *testing.T) {
	x := 0
	fnc := NewFuncList(func() { x++ })
	fnc.Add(func() { x-- })
	f := fnc.Add(func() { x++ })
	fnc.Remove(f)
	fnc.Add(func() { x-- })
	fnc.Add(func() { x++ })
	fnc.Run()
	if x != 0 {
		t.Fatalf("Got %d need %d\n", x, 0)
	}
}
