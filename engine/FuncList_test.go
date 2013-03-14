package engine

import "testing"

func BenchmarkFuncList(b *testing.B) {
	x := 0
	f := func() { x++ }
	fnc := NewFuncList(nil)
	for i := 0; i < 100; i++ {
		fnc.Add(f)
	}
	for i := 0; i < b.N; i++ {
		fnc.Run()
	}

	if x != 100*b.N {
		b.Errorf("Got %d need %d\n", x, 100*b.N)
	}
}

func TestFuncList(t *testing.T) {
	x := 0
	fnc := NewFuncList(func() { x++ })
	f := fnc.Add(func() { x-- })
	f1 := fnc.Add(func() { x++ })
	fnc.Remove(f)
	fnc.Remove(f1)
	fnc.Add(func() { x += 2 })
	fnc.Add(func() { x-- })
	fnc.Add(func() { x++ })
	fnc.Run()
	if x != 3 {
		t.Errorf("Got %d need %d\n", x, 3)
	}
}
