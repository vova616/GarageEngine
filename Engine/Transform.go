package Engine

import (
	"fmt"
	"math"
)

var (
	Zero     = Vector{0, 0, 0}
	Up       = Vector{0, 1, 0}
	Down     = Vector{0, -1, 0}
	Left     = Vector{-1, 0, 0}
	Right    = Vector{1, 0, 0}
	Forward  = Vector{0, 0, 1}
	Backward = Vector{0, 0, -1}
	One      = Vector{1, 1, 1}
	MinusOne = Vector{-1, -1, -1}
)

type Vector struct {
	X, Y, Z float32
}

func (v *Vector) String() string {
	return fmt.Sprintf("(%f,%f,%f)", v.X, v.Y, v.Z)
}

func NewVector2(x, y float32) Vector {
	return Vector{x, y, 1}
}

func NewVector3(x, y, z float32) Vector {
	return Vector{x, y, z}
}

func (v *Vector) Add(vect Vector) Vector {
	return Vector{v.X + vect.X, v.Y + vect.Y, v.Z + vect.Z}
}

func (v *Vector) Sub(vect Vector) Vector {
	return Vector{v.X - vect.X, v.Y - vect.Y, v.Z - vect.Z}
}

func (v *Vector) Mul(vect Vector) Vector {
	return Vector{v.X * vect.X, v.Y * vect.Y, v.Z * vect.Z}
}

func (v *Vector) Mul2(vect float32) Vector {
	return Vector{v.X * vect, v.Y * vect, v.Z * vect}
}

func (v *Vector) Distance(vect Vector) float32 {
	x := v.X - vect.X
	y := v.Y - vect.Y
	return float32(math.Sqrt(float64(x*x + y*y)))
}

func (v *Vector) Div(vect Vector) Vector {
	return Vector{v.X / vect.X, v.Y / vect.Y, v.Z / vect.Z}
}

func (v *Vector) Transform(transform Matrix) Vector {
	return NewVector3(
		(v.X*transform[0])+(v.Y*transform[4])+(v.Z*transform[8])+transform[12],
		(v.X*transform[1])+(v.Y*transform[5])+(v.Z*transform[9])+transform[13],
		(v.X*transform[2])+(v.Y*transform[6])+(v.Z*transform[10])+transform[14])
}

func (v *Vector) fixAngle() {
	for v.X >= 360 {
		v.X -= 360
	}
	for v.X <= -360 {
		v.X += 360
	}

	for v.Y >= 360 {
		v.Y -= 360
	}
	for v.Y <= -360 {
		v.Y += 360
	}

	for v.Z >= 360 {
		v.Z -= 360
	}
	for v.Z <= -360 {
		v.Z += 360
	}
}

func (v *Vector) Length() float32 {
	return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y)))
}

func (v *Vector) Normalize() {
	l := v.Length()
	v.X /= l
	v.Y /= l
	v.Z /= l
}

func (v *Vector) Normalized() Vector {
	l := v.Length()
	if l == 0 {
		return NewVector3(0, 0, 0)
	}
	return NewVector3(v.X/l, v.Y/l, v.Z/l)
}

type Transform struct {
	gameObject *GameObject
	parent     *Transform
	position   Vector
	rotation   Vector
	scale      Vector

	children []*Transform

	worldPosition Vector
	worldRotation Vector
	worldScale    Vector
	matrix        *Matrix
	parentMatrix  *Matrix
	updatedMatrix bool
}

func NewTransform(g *GameObject) *Transform {
	return &Transform{g, nil, Zero, Zero, One, make([]*Transform, 0), Zero, Zero, One, NewIdentity(), NewIdentity(), false}
}

func (t *Transform) Position() Vector {
	return t.position
}

func (t *Transform) Rotation() Vector {
	return t.rotation
}

func (t *Transform) Rotation2D() Vector {
	angle := (t.rotation.Z) * RadianConst
	return NewVector2(float32(math.Cos(float64(angle))), float32(math.Sin(float64(angle))))
}

func (t *Transform) Direction2D(up Vector) Vector {
	angle := float32(RadianConst)

	angle *= (t.rotation.Z + float32(math.Atan2(float64(up.Y), float64(up.X)))*float32(DegreeConst))

	return NewVector2(float32(math.Cos(float64(angle))), float32(math.Sin(float64(angle))))
}

func (t *Transform) Scale() Vector {
	return t.scale
}

func (t *Transform) SetPosition(vect Vector) {
	if vect == t.position {
		return
	}
	t.updatedMatrix = false
	t.position = vect
}

func (t *Transform) SetPositionf(x, y float32) {
	t.SetPosition(NewVector3(x, y, 1))
}

func (t *Transform) SetRotation(vect Vector) {
	if t.rotation == vect {
		return
	}
	t.updatedMatrix = false
	t.rotation = vect
}

func (t *Transform) SetRotationf(z float32) {
	t.SetRotation(NewVector3(0, 0, z))
}

func (t *Transform) SetScale(vect Vector) {
	if t.scale == vect {
		return
	}
	t.updatedMatrix = false
	t.scale = vect
}

func (t *Transform) SetScalef(x, y float32) {
	t.SetScale(NewVector3(x, y, 1))
}

func (t *Transform) WorldPosition() Vector {
	if t.parent == nil {
		return t.position
	}
	t.updateMatrix()
	return t.worldPosition
}

func (t *Transform) WorldRotation() Vector {
	if t.parent == nil {
		return t.rotation
	}
	t.updateMatrix()
	return t.worldRotation
}

func (t *Transform) WorldScale() Vector {
	if t.parent == nil {
		return t.scale
	}
	t.updateMatrix()
	return t.worldScale
}

func (t *Transform) SetWorldPosition(vect Vector) {
	if t.parent == nil {
		t.SetPosition(vect.Sub(t.position))
	} else {
		x := t.parent.Matrix()
		t.SetPosition(vect.Transform(x.Invert()))
	}
}

func (t *Transform) SetWorldPositionf(x, y float32) {
	t.SetWorldPosition(NewVector3(x, y, 1))
}

func (t *Transform) SetWorldRotation(vect Vector) {
	var p Vector
	if t.parent == nil {
		p = t.rotation
	} else {
		p = t.parent.WorldRotation()
	}

	t.SetRotation(vect.Sub(p))
}

func (t *Transform) SetWorldRotationf(z float32) {
	t.SetWorldRotation(NewVector3(0, 0, z))
}

func (t *Transform) SetWorldScale(vect Vector) {
	var p Vector
	if t.parent == nil {
		p = t.scale
	} else {
		p = t.parent.WorldScale()
	}

	t.SetScale(p.Div(t.scale))
}

func (t *Transform) SetWorldScalef(x, y float32) {
	t.SetWorldScale(NewVector3(x, y, 1))
}

func (t *Transform) Parent() *Transform {
	return t.parent
}

func (t *Transform) GameObject() *GameObject {
	return t.gameObject
}

func (t *Transform) Child(index int) *Transform {
	if index < len(t.children) {
		return t.children[index]
	}
	return nil
}

func (t *Transform) Children() []*Transform {
	arr := make([]*Transform, len(t.children))
	copy(arr, t.children)
	return arr
}

func (t *Transform) Translate(v Vector) {
	a := t.Position()
	t.SetPosition(a.Add(v))
}

func (t *Transform) Translatef(x, y float32) {
	t.Translate(NewVector3(x, y, 1))
}

func (t *Transform) SetParent(parent *Transform) {
	if t.parent != nil {
		for i, c := range t.parent.children {
			if t == c {
				t.parent.children = append(t.parent.children[:i], t.parent.children[i+1:]...)
				break
			}
		}
	}
	t.parent = parent
	t.updatedMatrix = false
	if parent != nil {
		parent.children = append(parent.children, t)
	}
}

func (t *Transform) SetParent2(g *GameObject) {
	if g == nil {
		t.SetParent(nil)
	} else {
		t.SetParent(g.transform)
	}
}

func (t *Transform) updateMatrix() {
	if t.updatedMatrix && ((t.parent != nil && t.parent.Matrix() == *t.parentMatrix) || t.parent == nil) {
		return
	}

	trans := t

	s, r, p := trans.scale, trans.rotation, trans.position

	trans.matrix.Reset()
	mat := trans.matrix

	mat.Scale(s.X, s.Y, s.Z)
	mat.Rotate(r.X, 1, 0, 0)
	mat.Rotate(r.Y, 0, 1, 0)
	mat.Rotate(r.Z, 0, 0, -1)

	mat.Translate(p.X, p.Y, p.Z)

	if trans.parent != nil {
		*t.parentMatrix = trans.parent.Matrix()
		mat.Mul(*t.parentMatrix)
		t.worldScale = trans.parent.worldScale.Mul(trans.scale)
		t.worldRotation = trans.parent.worldRotation.Add(trans.rotation)
	} else {
		t.worldScale = trans.scale
		t.worldRotation = trans.rotation
	}

	t.worldPosition = mat.Translation()

	//fmt.Println(t.GameObject().name)

	t.updatedMatrix = true
}

func (t *Transform) Matrix() Matrix {
	t.updateMatrix()
	return *t.matrix
}

func (t *Transform) clone(parent *GameObject) *Transform {
	tn := NewTransform(parent)
	tn.position = t.position
	tn.rotation = t.rotation
	tn.scale = t.scale
	for _, c := range t.children {
		c.gameObject.Clone().transform.SetParent(tn)
	}
	return tn
}
