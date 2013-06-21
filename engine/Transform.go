package engine

import (
	"math"
)

type Transform struct {
	BaseComponent
	parent   *Transform
	position Vector
	rotation Vector
	scale    Vector

	children []*Transform

	worldPosition Vector
	worldRotation Vector
	worldScale    Vector
	matrix        Matrix
	parentMatrix  Matrix
	updatedMatrix bool
	childOfScene  bool

	depth       int
	inDepthList bool

	inverted        Matrix
	updatedInverted bool
}

func NewTransform() *Transform {
	return &Transform{BaseComponent: NewComponent(), scale: One, matrix: Identity()}
}

func (t *Transform) OnComponentAdd() {
	t.gameObject.transform = t
}

func (t *Transform) SetDepth(depth int) {
	t.removeFromDepthMap()
	t.depth = depth
	//If object is in scene add to depth map
	if t.InScene() {
		t.updateDepth()
	}
}

func (t *Transform) SetDepthRecursive(depth int) {
	t.SetDepth(depth)
	for _, c := range t.children {
		c.SetDepthRecursive(depth)
	}
}

func (t *Transform) updateDepth() {
	if !t.inDepthList {
		depthMap.Add(int(t.depth), t.gameObject)
		t.inDepthList = true
	}
}

func (t *Transform) checkDepthRecursive() {
	t.updateDepth()
	for _, c := range t.children {
		c.checkDepthRecursive()
	}
}

func (t *Transform) removeFromDepthMap() {
	if t.inDepthList {
		depthMap.Remove(int(t.depth), t.gameObject)
	}
	t.inDepthList = false
}

func (t *Transform) removeFromDepthMapRecursive() {
	t.removeFromDepthMap()
	for _, c := range t.children {
		c.removeFromDepthMapRecursive()
	}
}

/*
Checking if object is somewhere in scene.
*/
func (t *Transform) InScene() bool {
	return t.childOfScene
}

func (t *Transform) Depth() int {
	return t.depth
}

func (t *Transform) Position() Vector {
	return t.position
}

func (t *Transform) Rotation() Vector {
	return t.rotation
}

func (t *Transform) Angle() float32 {
	return t.WorldRotation().Z
}

func (t *Transform) Direction() Vector {
	angle := t.Angle() * RadianConst
	return NewVector2(float32(math.Cos(float64(angle))), float32(math.Sin(float64(angle))))
}

func (t *Transform) DirectionTransform(up Vector) Vector {
	angle := float32(RadianConst)
	angle *= (t.Angle() + float32(math.Atan2(float64(up.Y), float64(up.X)))*float32(DegreeConst))
	return NewVector2(float32(math.Cos(float64(angle))), float32(math.Sin(float64(angle))))
}

func (t *Transform) Scale() Vector {
	return t.scale
}

func (t *Transform) SetPosition(vect Vector) {
	t.updatedMatrix = false
	t.position = vect
}

func (t *Transform) SetPositionf(x, y float32) {
	t.updatedMatrix = false
	t.position.X, t.position.Y = x, y
}

func (t *Transform) SetRotation(vect Vector) {
	t.updatedMatrix = false
	t.rotation = vect
}

func (t *Transform) SetRotationf(z float32) {
	t.updatedMatrix = false
	t.rotation.Z = z
}

func (t *Transform) SetScale(vect Vector) {
	t.updatedMatrix = false
	t.scale = vect
}

func (t *Transform) SetScalef(x, y float32) {
	t.updatedMatrix = false
	t.scale.X, t.scale.Y = x, y
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
		t.SetPosition(vect)
		return
	}
	t.SetPosition(vect.Transform(t.parent.InvertedMatrix()))
}

func (t *Transform) SetWorldPositionf(x, y float32) {
	t.SetWorldPosition(NewVector3(x, y, t.worldPosition.Z))
}

func (t *Transform) SetWorldRotation(vect Vector) {
	if t.parent == nil {
		t.SetRotation(vect)
	} else {
		t.worldRotation = vect
		t.SetRotation(vect.Sub(t.parent.WorldRotation()))
	}
}

func (t *Transform) SetWorldRotationf(z float32) {
	t.SetWorldRotation(NewVector3(t.worldRotation.X, t.worldRotation.Y, z))
}

func (t *Transform) SetWorldScale(vect Vector) {
	if t.parent == nil {
		t.SetScale(vect)
	} else {
		t.worldScale = vect
		t.SetScale(vect.Div(t.parent.WorldScale()))
	}
}

func (t *Transform) SetWorldScalef(x, y float32) {
	t.SetWorldScale(NewVector3(x, y, t.worldScale.Z))
}

func (t *Transform) Parent() *Transform {
	return t.parent
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
	t.Translate(NewVector3(x, y, 0))
}

func (t *Transform) removeFromParent() {
	parent := t.parent
	if parent != nil {
		for i, c := range parent.children {
			if t == c {
				parent.children[len(parent.children)-1], parent.children[i], parent.children = nil, parent.children[len(parent.children)-1], parent.children[:len(parent.children)-1]
				break
			}
		}
	}
}

func (t *Transform) SetParent(newParent *Transform) {
	//if current parent is the requested parent, do nothing.
	if t.parent == newParent {
		if newParent == nil {
			if t.childOfScene {
				return
			}
		} else {
			return
		}
	}

	//Remove this transform from his parent
	t.removeFromParent()

	//Update depth
	t.updateDepth()

	//check if object was outside of scene and now it is
	wasOutsideScene := !t.InScene() && (newParent == nil || newParent.InScene())

	//Keep the position after changing parents
	if wasOutsideScene {
		scale := t.Scale()
		position := t.Position()
		rotation := t.Rotation()
		t.parent = newParent
		t.updatedMatrix = false
		t.SetPosition(position)
		t.SetRotation(rotation)
		t.SetScale(scale)
	} else {
		scale := t.WorldScale()
		position := t.WorldPosition()
		rotation := t.WorldRotation()
		t.parent = newParent
		t.updatedMatrix = false
		t.SetWorldPosition(position)
		t.SetWorldRotation(rotation)
		t.SetWorldScale(scale)
	}

	//Add this transform to newParent
	if newParent != nil {
		newParent.children = append(newParent.children, t)
	}

	//If object was outside of scene active it silencly and check if depth needs to be updated
	if wasOutsideScene {
		GetScene().SceneBase().addGameObject(t.gameObject)
		if t.gameObject.active {
			t.gameObject.silentActive(true)
		}
		t.checkDepthRecursive()
	}
}

func (t *Transform) SetParent2(g *GameObject) {
	if g == nil {
		t.SetParent(nil)
	} else {
		t.SetParent(g.transform)
	}
}

/*
Faster option will be to use Stamps, each transform will have stamp and parterStamp
*/
func (t *Transform) updateMatrix() bool {
	if t.updatedMatrix {
		if t.parent != nil {
			t.parent.updateMatrix()
			if t.parent.matrix == t.parentMatrix {
				return false
			}
		} else {
			return false
		}
	}

	trans := t

	s, r, p := trans.scale, trans.rotation, trans.position

	trans.matrix.Reset()
	mat := &trans.matrix

	mat.Scale(s.X, s.Y, s.Z)
	mat.RotateXYZ(r.X, r.Y, r.Z)
	mat.Translate(p.X, p.Y, p.Z)

	if trans.parent != nil {
		trans.parent.updateMatrix()
		t.parentMatrix = trans.parent.matrix
		mat.MulPtr(&t.parentMatrix)
		t.worldScale = trans.parent.worldScale.Mul(trans.scale)
		t.worldRotation = trans.parent.worldRotation.Add(trans.rotation)
	} else {
		t.worldScale = trans.scale
		t.worldRotation = trans.rotation
	}

	t.worldPosition = mat.Translation()

	//fmt.Println(t.GameObject().name)

	t.updatedMatrix = true
	t.updatedInverted = false
	return true
}

func (t *Transform) Matrix() Matrix {
	t.updateMatrix()
	return t.matrix
}

func (t *Transform) InvertedMatrix() Matrix {
	t.updateMatrix()
	if !t.updatedInverted {
		t.inverted = t.matrix.Invert()
		t.updatedInverted = true
	}
	return t.inverted
}

func (t *Transform) Clone() {
	chld := t.children
	t.parent = nil
	t.childOfScene = false
	t.inDepthList = false
	t.children = make([]*Transform, len(t.children))
	t.updatedInverted = false
	t.updatedMatrix = false
	for _, c := range chld {
		g := c.GameObject()
		if g != nil {
			g.Clone().transform.SetParent(t)
		}
	}
}
