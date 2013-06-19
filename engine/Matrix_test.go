package engine

import "testing"

func TestMul(t *testing.T) {
	tests := []struct {
		A, B Matrix
		Out  Matrix
	}{
		{Identity(), Identity(), Identity()},
		{
			Matrix{
				1, 0, 0, 0,
				0, 1, 0, 0,
				0, 0, 1, 0,
				4, 5, 6, 1,
			},
			Matrix{
				1, 0, 0, 0,
				0, 1, 0, 0,
				0, 0, 1, 0,
				1, 2, 3, 1,
			},
			Matrix{
				1, 0, 0, 0,
				0, 1, 0, 0,
				0, 0, 1, 0,
				5, 7, 9, 1,
			},
		},
		{
			Matrix{
				1, 0, 0, 0,
				0, 0, 1, 0,
				0, -1, 0, 0,
				0, 0, 0, 1,
			},
			Matrix{
				0, 0, -1, 0,
				0, 1, 0, 0,
				1, 0, 0, 0,
				0, 0, 0, 1,
			},
			Matrix{
				0, 0, -1, 0,
				1, 0, 0, 0,
				0, -1, 0, 0,
				0, 0, 0, 1,
			},
		},
	}
	for _, test := range tests {
		out := test.A
		out.Mul(test.B)
		if !checkMatrix(out, test.Out, 0.01) {
			t.Errorf("Mul(\n%v,\n%v) =\n%v; want\n%v", test.A, test.B, out, test.Out)
		}
		out = Mul(test.A, test.B)
		if !checkMatrix(out, test.Out, 0.01) {
			t.Errorf("Mul(\n%v,\n%v) =\n%v; want\n%v", test.A, test.B, out, test.Out)
		}
		out = test.A
		out.MulPtr(&test.B)
		if !checkMatrix(out, test.Out, 0.01) {
			t.Errorf("Mul(\n%v,\n%v) =\n%v; want\n%v", test.A, test.B, out, test.Out)
		}
	}
}

// checkMatrix returns whether m1 ~ m2, given a tolerance.
func checkMatrix(m1, m2 Matrix, tol float32) bool {
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if m2[i*4+j] > m1[i*4+j]+tol || m2[i*4+j] < m1[i*4+j]-tol {
				return false
			}
		}
	}
	return true
}

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

func BenchmarkMatrix_RotateXYZ(b *testing.B) {
	b.StopTimer()
	m := Identity()
	m.Translate(10, 20, 30)
	m.Rotate(10, 1, 1, 1)
	m.Scale(2, 3, 4)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		m.RotateXYZ(10, 20, 30)
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
