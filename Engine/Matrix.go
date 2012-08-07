package Engine

import (
"math"
"unsafe"
)

const (
	PI = math.Pi
	PI_180 = math.Pi / 180
	PI_360 = math.Pi / 360
)

type Matrix [16]float32

func (mA *Matrix) Mul(mB Matrix) {
	NewMatrix := Identity()
	for i := 0; i < 4; i++ { //Cycle through each vector of first matrix.
		NewMatrix[i*4] = mA[i*4]*mB[0] + mA[i*4+1]*mB[4] + mA[i*4+2]*mB[8] + mA[i*4+3]*mB[12]
		NewMatrix[i*4+1] = mA[i*4]*mB[1] + mA[i*4+1]*mB[5] + mA[i*4+2]*mB[9] + mA[i*4+3]*mB[13]
		NewMatrix[i*4+2] = mA[i*4]*mB[2] + mA[i*4+1]*mB[6] + mA[i*4+2]*mB[10] + mA[i*4+3]*mB[14]
		NewMatrix[i*4+3] = mA[i*4]*mB[3] + mA[i*4+1]*mB[7] + mA[i*4+2]*mB[11] + mA[i*4+3]*mB[15]
	}
	*mA = NewMatrix
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
    
    det1 := mA[0] * mA[5] - mA[1] * mA[4];
    det2 := mA[0] * mA[6] - mA[2] * mA[4];
    det3 := mA[0] * mA[7] - mA[3] * mA[4];
    det4 := mA[1] * mA[6] - mA[2] * mA[5];
    det5 := mA[1] * mA[7] - mA[3] * mA[5];
    det6 := mA[2] * mA[7] - mA[3] * mA[6];
    det7 := mA[8] * mA[13] - mA[9] * mA[12];
    det8 := mA[8] * mA[14] - mA[10] * mA[12];
    det9 := mA[8] * mA[15] - mA[11] * mA[12];
    det10 := mA[9] * mA[14] - mA[10] * mA[13];
    det11 := mA[9] * mA[15] - mA[11] * mA[13];
    det12 := mA[10] * mA[15] - mA[11] * mA[14];
    
    detmA := (float32)(det1*det12 - det2*det11 + det3*det10 + det4*det9 - det5*det8 + det6*det7);
    
    invDetmA := 1 / detmA;
    
    var ret Matrix; // Allow for mA and result to point to the same structure
    
    ret[0] = (mA[5]*det12 - mA[6]*det11 + mA[7]*det10) * invDetmA;
    ret[1] = (-mA[1]*det12 + mA[2]*det11 - mA[3]*det10) * invDetmA;
    ret[2] = (mA[13]*det6 - mA[14]*det5 + mA[15]*det4) * invDetmA;
    ret[3] = (-mA[9]*det6 + mA[10]*det5 - mA[11]*det4) * invDetmA;
    ret[4] = (-mA[4]*det12 + mA[6]*det9 - mA[7]*det8) * invDetmA;
    ret[5] = (mA[0]*det12 - mA[2]*det9 + mA[3]*det8) * invDetmA;
    ret[6] = (-mA[12]*det6 + mA[14]*det3 - mA[15]*det2) * invDetmA;
    ret[7] = (mA[8]*det6 - mA[10]*det3 + mA[11]*det2) * invDetmA;
    ret[8] = (mA[4]*det11 - mA[5]*det9 + mA[7]*det7) * invDetmA;
    ret[9] = (-mA[0]*det11 + mA[1]*det9 - mA[3]*det7) * invDetmA;
    ret[10] = (mA[12]*det5 - mA[13]*det3 + mA[15]*det1) * invDetmA;
    ret[11] = (-mA[8]*det5 + mA[9]*det3 - mA[11]*det1) * invDetmA;
    ret[12] = (-mA[4]*det10 + mA[5]*det8 - mA[6]*det7) * invDetmA;
    ret[13] = (mA[0]*det10 - mA[1]*det8 + mA[2]*det7) * invDetmA;
    ret[14] = (-mA[12]*det4 + mA[13]*det2 - mA[14]*det1) * invDetmA;
    ret[15] = (mA[8]*det4 - mA[9]*det2 + mA[10]*det1) * invDetmA;
    
    return ret
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
	return 	NewVector3(mA[12],mA[13],mA[14])
}

func NewIdentity() *Matrix {
	return &Matrix{1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1}
}

func (mA *Matrix) Scale(x,y,z float32) {
	m2 := Identity()
	m2[0] = x;
	m2[5] = y;
	m2[10] = z;
	mA.Mul(m2)
}

func (mA *Matrix) Translate(x,y,z float32) {
	m2 := Identity()
	m2[12] = x;
	m2[13] = y;
	m2[14] = z;
	mA.Mul(m2)
}

func (mA *Matrix) Ortho( left,  right,  bottom,  top,  Znear,  Zfar float32) {
	m2 := Identity()
	
	m2[0] = 2/(right-left);
	m2[3] = -((right+left)/(right-left));
	
	m2[5] = 2/(top-bottom);
	m2[7] = -((top+bottom)/(top-bottom));
	
	m2[10] = 2/(Zfar-Znear);
	m2[11] = -((Zfar+Znear)/(Zfar-Znear));

	mA.Mul(m2)
}


func (mA *Matrix) Rotate(a,x,y,z float32) {
	m2 := Identity()
	angle := a*PI_180;
	
	acos := float32(math.Cos(float64(angle)))
	asin := float32(math.Sin(float64(angle)))
	
	m2[0] = 1+(1-acos)*(x*x-1);
	m2[1] = -z*asin+(1-acos)*x*y;
	m2[2] = y*asin+(1-acos)*x*z;
	
	m2[4] = z*asin+(1-acos)*x*y;
	m2[5] = 1+(1-acos)*(y*y-1);
	m2[6] = -x*asin+(1-acos)*y*z;
	
	m2[8] = -y*asin+(1-acos)*x*z;
	m2[9] = x*asin+(1-acos)*y*z;
	m2[10] = 1+(1-acos)*(z*z-1);
	
	mA.Mul(m2)
}
