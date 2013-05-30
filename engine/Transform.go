package engine

import (
	"math"
)

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
	matrix        Matrix
	parentMatrix  Matrix
	updatedMatrix bool
	childOfScene  bool

	depth       int
	inDepthList bool
}

func NewTransform(g *GameObject) *Transform {
	return &Transform{g, nil, Zero, Zero, One, make([]*Transform, 0), Zero, Zero, One, Identity(), Identity(), false, false, 0, false}
}

func (t *Transform) SetDepth(depth int) {
	t.removeFromDepthMap()
	t.depth = depth
	//If object is in scene add to depth map
	if t.InScene() {
		t.checkDepth()
	}
}

func (t *Transform) SetDepthRecursive(depth int) {
	t.SetDepth(depth)
	for _, c := range t.children {
		c.SetDepthRecursive(depth)
	}
}

func (t *Transform) checkDepth() {
	if !t.inDepthList {
		depthMap.Add(int(t.depth), t.gameObject)
		t.inDepthList = true
	}
}

func (t *Transform) checkDepthRecursive() {
	t.checkDepth()
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
	return t.childOfScene || (t.parent != nil && t.parent.InScene())
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
	t.parent.updateMatrix()
	t.SetPosition(vect.Transform(t.parent.matrix.Invert()))
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

func (t *Transform) GameObject() *GameObject {
	return t.gameObject.GameObject()
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

func (t *Transform) SetParent(parent *Transform) {
	if t.parent == parent {
		if parent == nil {
			if t.childOfScene {
				return
			}
		} else {
			return
		}
	}

	//Remove from scene if needed
	if t.childOfScene && parent != nil {
		GetScene().SceneBase().removeGameObject(t.gameObject)
	}

	//if parent is not nil, remove this transform from parent
	if t.parent != nil {
		for i, c := range t.parent.children {
			if t == c {
				t.parent.children = append(t.parent.children[:i], t.parent.children[i+1:]...)
				break
			}
		}
	}

	//if depth is not updated, update it
	t.checkDepth()

	//check if object was outside of scene
	wasOutsideScene := !t.InScene()

	//Keep the position after changing parents
	//
	if wasOutsideScene {
		scale := t.Scale()
		position := t.Position()
		rotation := t.Rotation()
		t.parent = parent
		t.updatedMatrix = false
		t.SetPosition(position)
		t.SetRotation(rotation)
		t.SetScale(scale)
	} else {
		scale := t.WorldScale()
		position := t.WorldPosition()
		rotation := t.WorldRotation()
		t.parent = parent
		t.updatedMatrix = false
		t.SetWorldPosition(position)
		t.SetWorldRotation(rotation)
		t.SetWorldScale(scale)
	}

	//If parent is nil add to scene else add to transform children
	if parent != nil {
		parent.children = append(parent.children, t)
	} else {
		GetScene().SceneBase().addGameObject(t.gameObject)
	}

	//If object was outside of scene active it silencly and check if depth needs to be updated
	if wasOutsideScene {
		t.gameObject.setActiveRecursiveSilent(true)
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
	mat.RotateX(r.X, 1)
	mat.RotateY(r.Y, 1)
	mat.RotateZ(r.Z, -1)
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
	return true
}

func (t *Transform) Matrix() Matrix {
	t.updateMatrix()
	return t.matrix
}

func (t *Transform) clone(parent *GameObject) *Transform {
	tn := NewTransform(parent)
	tn.position = t.position
	tn.rotation = t.rotation
	tn.scale = t.scale
	for _, c := range t.children {
		c.gameObject.Clone().transform.SetParent(tn)
	}
	tn.depth = t.depth
	return tn
}
