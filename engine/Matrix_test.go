package engine

import "testing"

func BenchmarkMatrix_Mult(b *testing.B) {
	b.StopTimer()
	m := Identity()
	m2 := Identity()
	m2.Translate(10, 20, 30)
	m2.Rotate(10, 1, 1, 1)
	m2.Scale(2, 3, 4)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		m.MulPtr(&m2)
	}
}

func BenchmarkMatrix_Scale(b *testing.B) {
	b.StopTimer()
	m := Identity()
	m.Translate(10, 20, 30)
	m.Rotate(10, 1, 1, 1)
	m.Scale(2, 3, 4)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		m.Scale(10, 20, 30)
	}
}

func BenchmarkMatrix_Translate(b *testing.B) {
	b.StopTimer()
	m := Identity()
	m.Translate(10, 20, 30)
	m.Rotate(10, 1, 1, 1)
	m.Scale(2, 3, 4)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		m.Translate(10, 20, 30)
	}
}

func BenchmarkMatrix_RotateX(b *testing.B) {
	b.StopTimer()
	m := Identity()
	m.Translate(10, 20, 30)
	m.Rotate(10, 1, 1, 1)
	m.Scale(2, 3, 4)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		m.RotateX(10, 20)
	}
}

func BenchmarkMatrix_RotateY(b *testing.B) {
	b.StopTimer()
	m := Identity()
	m.Translate(10, 20, 30)
	m.Rotate(10, 1, 1, 1)
	m.Scale(2, 3, 4)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		m.RotateY(10, 20)
	}
}

func BenchmarkMatrix_RotateZ(b *testing.B) {
	b.StopTimer()
	m := Identity()
	m.Translate(10, 20, 30)
	m.Rotate(10, 1, 1, 1)
	m.Scale(2, 3, 4)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		m.RotateZ(10, 20)
	}
}

func BenchmarkMatrix_Invert(b *testing.B) {
	b.StopTimer()
	m := Identity()
	m.Translate(10, 20, 30)
	m.Rotate(10, 1, 1, 1)
	m.Scale(2, 3, 4)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		m = m.Invert()
	}
}
