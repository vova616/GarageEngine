package engine

import (
	"github.com/vova616/garageEngine/engine/input"
)

type Camera struct {
	BaseComponent
	Projection *Matrix
	Size       float32
}

func NewCamera() *Camera {
	c := &Camera{NewComponent(), NewIdentity(), 1}
	c.UpdateResolution()
	return c
}

func (c *Camera) Update() {
	/*
		w := float32(Width)/2
		h := float32(Height)/2
		proj := NewIdentity()
		proj.Ortho(-w, w, -h, h, -1000, 1000) 
		c.Projection = proj
	*/
}

func (c *Camera) InvertedMatrix() Matrix {
	c.Transform().updateMatrix()
	invert := c.Transform().matrix.Invert()
	return invert
}

func (c *Camera) UpdateResolution() {
	//w := float32(Width) * c.Size * 0.5
	//h := float32(Height) * c.Size * 0.5
	//c.Projection.Ortho(-w, w, -h, h, -1000, 1000) 
	c.Projection.Ortho(0, float32(Width)*c.Size, 0, float32(Height)*c.Size, -1000, 1000)
}

func (c *Camera) MouseWorldPosition() Vector {
	x, y := input.MousePosition()
	x, y = x, (Height)-y

	return c.ScreenToWorld(x, y)
}

func (c *Camera) MouseLocalPosition() Vector {
	x, y := input.MousePosition()
	x, y = x, (Height)-y

	return NewVector2(float32(x), float32(y))
}

func (c *Camera) ScreenToWorld(x, y int) Vector {

	m := Identity()
	m.Translate(float32(x), float32(y), 0)
	m = Mul(m, c.Transform().Matrix())

	return m.Translation()
}

func (c *Camera) Render() {
	s := GetScene()
	if s != nil {
		tcam := s.SceneBase().Camera
		s.SceneBase().Camera = c
		arr := s.SceneBase().gameObjects
		if arr == nil {
			println("arr")
		}
		if c.GameObject() == nil {
			println("c.GameObject()")
		}

		IterExcept(arr, drawGameObject, c.GameObject())
		s.SceneBase().Camera = tcam
	}
}
