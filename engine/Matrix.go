package engine

import (
	"math"
	"unsafe"
)

const (
	PI     = math.Pi
	PI_180 = math.Pi / 180
	PI_360 = math.Pi / 360
)

type Matrix [16]float32

func Mul(m1, m2 Matrix) Matrix {
	return Matrix{
		m1[0]*m2[0] + m1[1]*m2[4] + m1[2]*m2[8] + m1[3]*m2[12],
		m1[0]*m2[1] + m1[1]*m2[5] + m1[2]*m2[9] + m1[3]*m2[13],
		m1[0]*m2[2] + m1[1]*m2[6] + m1[2]*m2[10] + m1[3]*m2[14],
		m1[0]*m2[3] + m1[1]*m2[7] + m1[2]*m2[11] + m1[3]*m2[15],

		m1[4]*m2[0] + m1[5]*m2[4] + m1[6]*m2[8] + m1[7]*m2[12],
		m1[4]*m2[1] + m1[5]*m2[5] + m1[6]*m2[9] + m1[7]*m2[13],
		m1[4]*m2[2] + m1[5]*m2[6] + m1[6]*m2[10] + m1[7]*m2[14],
		m1[4]*m2[3] + m1[5]*m2[7] + m1[6]*m2[11] + m1[7]*m2[15],

		m1[8]*m2[0] + m1[9]*m2[4] + m1[10]*m2[8] + m1[11]*m2[12],
		m1[8]*m2[1] + m1[9]*m2[5] + m1[10]*m2[9] + m1[11]*m2[13],
		m1[8]*m2[2] + m1[9]*m2[6] + m1[10]*m2[10] + m1[11]*m2[14],
		m1[8]*m2[3] + m1[9]*m2[7] + m1[10]*m2[11] + m1[11]*m2[15],

		m1[12]*m2[0] + m1[13]*m2[4] + m1[14]*m2[8] + m1[15]*m2[12],
		m1[12]*m2[1] + m1[13]*m2[5] + m1[14]*m2[9] + m1[15]*m2[13],
		m1[12]*m2[2] + m1[13]*m2[6] + m1[14]*m2[10] + m1[15]*m2[14],
		m1[12]*m2[3] + m1[13]*m2[7] + m1[14]*m2[11] + m1[15]*m2[15],
	}
}

func (m1 *Matrix) Mul(m2 Matrix) {
	*m1 = Matrix{
		m1[0]*m2[0] + m1[1]*m2[4] + m1[2]*m2[8] + m1[3]*m2[12],
		m1[0]*m2[1] + m1[1]*m2[5] + m1[2]*m2[9] + m1[3]*m2[13],
		m1[0]*m2[2] + m1[1]*m2[6] + m1[2]*m2[10] + m1[3]*m2[14],
		m1[0]*m2[3] + m1[1]*m2[7] + m1[2]*m2[11] + m1[3]*m2[15],

		m1[4]*m2[0] + m1[5]*m2[4] + m1[6]*m2[8] + m1[7]*m2[12],
		m1[4]*m2[1] + m1[5]*m2[5] + m1[6]*m2[9] + m1[7]*m2[13],
		m1[4]*m2[2] + m1[5]*m2[6] + m1[6]*m2[10] + m1[7]*m2[14],
		m1[4]*m2[3] + m1[5]*m2[7] + m1[6]*m2[11] + m1[7]*m2[15],

		m1[8]*m2[0] + m1[9]*m2[4] + m1[10]*m2[8] + m1[11]*m2[12],
		m1[8]*m2[1] + m1[9]*m2[5] + m1[10]*m2[9] + m1[11]*m2[13],
		m1[8]*m2[2] + m1[9]*m2[6] + m1[10]*m2[10] + m1[11]*m2[14],
		m1[8]*m2[3] + m1[9]*m2[7] + m1[10]*m2[11] + m1[11]*m2[15],

		m1[12]*m2[0] + m1[13]*m2[4] + m1[14]*m2[8] + m1[15]*m2[12],
		m1[12]*m2[1] + m1[13]*m2[5] + m1[14]*m2[9] + m1[15]*m2[13],
		m1[12]*m2[2] + m1[13]*m2[6] + m1[14]*m2[10] + m1[15]*m2[14],
		m1[12]*m2[3] + m1[13]*m2[7] + m1[14]*m2[11] + m1[15]*m2[15],
	}
}

func (m1 *Matrix) MulPtr(m2 *Matrix) {
	*m1 = Matrix{
		m1[0]*m2[0] + m1[1]*m2[4] + m1[2]*m2[8] + m1[3]*m2[12],
		m1[0]*m2[1] + m1[1]*m2[5] + m1[2]*m2[9] + m1[3]*m2[13],
		m1[0]*m2[2] + m1[1]*m2[6] + m1[2]*m2[10] + m1[3]*m2[14],
		m1[0]*m2[3] + m1[1]*m2[7] + m1[2]*m2[11] + m1[3]*m2[15],

		m1[4]*m2[0] + m1[5]*m2[4] + m1[6]*m2[8] + m1[7]*m2[12],
		m1[4]*m2[1] + m1[5]*m2[5] + m1[6]*m2[9] + m1[7]*m2[13],
		m1[4]*m2[2] + m1[5]*m2[6] + m1[6]*m2[10] + m1[7]*m2[14],
		m1[4]*m2[3] + m1[5]*m2[7] + m1[6]*m2[11] + m1[7]*m2[15],

		m1[8]*m2[0] + m1[9]*m2[4] + m1[10]*m2[8] + m1[11]*m2[12],
		m1[8]*m2[1] + m1[9]*m2[5] + m1[10]*m2[9] + m1[11]*m2[13],
		m1[8]*m2[2] + m1[9]*m2[6] + m1[10]*m2[10] + m1[11]*m2[14],
		m1[8]*m2[3] + m1[9]*m2[7] + m1[10]*m2[11] + m1[11]*m2[15],

		m1[12]*m2[0] + m1[13]*m2[4] + m1[14]*m2[8] + m1[15]*m2[12],
		m1[12]*m2[1] + m1[13]*m2[5] + m1[14]*m2[9] + m1[15]*m2[13],
		m1[12]*m2[2] + m1[13]*m2[6] + m1[14]*m2[10] + m1[15]*m2[14],
		m1[12]*m2[3] + m1[13]*m2[7] + m1[14]*m2[11] + m1[15]*m2[15],
	}
}

func (mA *Matrix) Ptr() *float32 {
	return (*float32)(unsafe.Pointer(&mA[0]))
}

func (mA *Matrix) Invert() Matrix {
	//
	// Use Laplace expansion theorem to calculate the inverse of a 4x4 Matrix
	//
	// 1. Calculate the 2x2 determinants needed and the 4x4 determinant based on the 2x2 determinants
	// 2. Create the adjugate Matrix, which satisfies: A * adj(A) = det(A) * I
	// 3. Divide adjugate Matrix with the determinant to find the inverse

	det1 := mA[0]*mA[5] - mA[1]*mA[4]
	det2 := mA[0]*mA[6] - mA[2]*mA[4]
	det3 := mA[0]*mA[7] - mA[3]*mA[4]
	det4 := mA[1]*mA[6] - mA[2]*mA[5]
	det5 := mA[1]*mA[7] - mA[3]*mA[5]
	det6 := mA[2]*mA[7] - mA[3]*mA[6]
	det7 := mA[8]*mA[13] - mA[9]*mA[12]
	det8 := mA[8]*mA[14] - mA[10]*mA[12]
	det9 := mA[8]*mA[15] - mA[11]*mA[12]
	det10 := mA[9]*mA[14] - mA[10]*mA[13]
	det11 := mA[9]*mA[15] - mA[11]*mA[13]
	det12 := mA[10]*mA[15] - mA[11]*mA[14]

	invDetmA := 1 / (det1*det12 - det2*det11 + det3*det10 + det4*det9 - det5*det8 + det6*det7)

	return Matrix{(mA[5]*det12 - mA[6]*det11 + mA[7]*det10) * invDetmA,
		(-mA[1]*det12 + mA[2]*det11 - mA[3]*det10) * invDetmA,
		(mA[13]*det6 - mA[14]*det5 + mA[15]*det4) * invDetmA,
		(-mA[9]*det6 + mA[10]*det5 - mA[11]*det4) * invDetmA,
		(-mA[4]*det12 + mA[6]*det9 - mA[7]*det8) * invDetmA,
		(mA[0]*det12 - mA[2]*det9 + mA[3]*det8) * invDetmA,
		(-mA[12]*det6 + mA[14]*det3 - mA[15]*det2) * invDetmA,
		(mA[8]*det6 - mA[10]*det3 + mA[11]*det2) * invDetmA,
		(mA[4]*det11 - mA[5]*det9 + mA[7]*det7) * invDetmA,
		(-mA[0]*det11 + mA[1]*det9 - mA[3]*det7) * invDetmA,
		(mA[12]*det5 - mA[13]*det3 + mA[15]*det1) * invDetmA,
		(-mA[8]*det5 + mA[9]*det3 - mA[11]*det1) * invDetmA,
		(-mA[4]*det10 + mA[5]*det8 - mA[6]*det7) * invDetmA,
		(mA[0]*det10 - mA[1]*det8 + mA[2]*det7) * invDetmA,
		(-mA[12]*det4 + mA[13]*det2 - mA[14]*det1) * invDetmA,
		(mA[8]*det4 - mA[9]*det2 + mA[10]*det1) * invDetmA}
}

func (mA *Matrix) Reset() {
	*mA = Identity()
}

func Identity() Matrix {
	return Matrix{1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1}
}

func (mA *Matrix) Translation() Vector {
	return NewVector3(mA[12], mA[13], mA[14])
}

func NewIdentity() *Matrix {
	return &Matrix{1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1}
}

func (mA *Matrix) Scale(x, y, z float32) {
	m2 := Identity()
	m2[0] = x
	m2[5] = y
	m2[10] = z
	mA.Mul(m2)
}

func (mA *Matrix) Translate(x, y, z float32) {
	m2 := Identity()
	m2[12] = x
	m2[13] = y
	m2[14] = z
	mA.Mul(m2)
}

func (mA *Matrix) Ortho(left, right, bottom, top, Znear, Zfar float32) {
	m2 := Identity()
	/*
		m2[0] = 2/(right-left);
		m2[5] = 2/(top-bottom);
		m2[10] = -2/(Zfar-Znear);

		m2[3] = -((right+left)/(right-left));
		m2[7] = -((top+bottom)/(top-bottom));
		m2[14] = -((Zfar+Znear)/(Zfar-Znear));
	*/

	m2[0] = 2 / (right - left)
	m2[5] = 2 / (top - bottom)
	m2[10] = -2 / (Zfar - Znear)

	m2[12] = -((right + left) / (right - left))
	m2[13] = -((top + bottom) / (top - bottom))
	m2[14] = -((Zfar + Znear) / (Zfar - Znear))

	*mA = m2
}

func (mA *Matrix) Rotate(a, x, y, z float32) {
	m2 := Identity()
	angle := a * PI_180

	acos := float32(math.Cos(float64(angle)))
	asin := float32(math.Sin(float64(angle)))

	m2[0] = 1 + (1-acos)*(x*x-1)
	m2[1] = -z*asin + (1-acos)*x*y
	m2[2] = y*asin + (1-acos)*x*z

	m2[4] = z*asin + (1-acos)*x*y
	m2[5] = 1 + (1-acos)*(y*y-1)
	m2[6] = -x*asin + (1-acos)*y*z

	m2[8] = -y*asin + (1-acos)*x*z
	m2[9] = x*asin + (1-acos)*y*z
	m2[10] = 1 + (1-acos)*(z*z-1)

	mA.Mul(m2)
}

func (mA *Matrix) RotateX(a, x float32) {
	m2 := Identity()
	angle := a * PI_180

	acos := float32(math.Cos(float64(angle)))
	asin := float32(math.Sin(float64(angle)))

	m2[0] = 1 + (1-acos)*(x*x-1)
	m2[5] = 1 - (1 - acos)
	m2[6] = -x * asin
	m2[9] = x * asin
	m2[10] = 1 - (1 - acos)

	mA.Mul(m2)
}

func (mA *Matrix) RotateY(a, y float32) {
	m2 := Identity()
	angle := a * PI_180

	acos := float32(math.Cos(float64(angle)))
	asin := float32(math.Sin(float64(angle)))

	m2[0] = 1 - (1 - acos)
	m2[2] = y * asin
	m2[5] = 1 + (1-acos)*(y*y-1)
	m2[8] = -y * asin
	m2[10] = 1 - (1 - acos)

	mA.Mul(m2)
}

func (mA *Matrix) RotateZ(a, z float32) {
	m2 := Identity()
	angle := a * PI_180

	acos := float32(math.Cos(float64(angle)))
	asin := float32(math.Sin(float64(angle)))

	m2[0] = 1 - (1 - acos)
	m2[1] = -z * asin
	m2[4] = z * asin
	m2[5] = 1 - (1 - acos)
	m2[10] = 1 + (1-acos)*(z*z-1)

	mA.Mul(m2)
}
