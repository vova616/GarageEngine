package engine

import (
	"github.com/vova616/garageEngine/engine/input"
)

type Rect struct {
	Min, Max Vector
}

func (r *Rect) Overlaps(s Rect) bool {
	return r.Min.X < s.Max.X && s.Min.X < r.Max.X &&
		r.Min.Y < s.Max.Y && s.Min.Y < r.Max.Y
}

type Camera struct {
	BaseComponent
	Projection *Matrix

	size        float32
	rect        Rect
	realRect    Rect
	center      Vector
	sizeIsScale bool
}

func NewCamera() *Camera {
	c := &Camera{BaseComponent: NewComponent(), Projection: NewIdentity(), size: 1, sizeIsScale: true}

	/*
		c.rect.Min.X - left
		c.rect.Max.X - right
		c.rect.Min.Y - bottom
		c.rect.Max.Y - top
	*/

	c.center = NewVector2(0, 0)
	c.rect.Min.X = -float32(Width) / 2
	c.rect.Max.X = float32(Width) / 2
	c.rect.Min.Y = -float32(Height) / 2
	c.rect.Max.Y = float32(Height) / 2

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

func (c *Camera) SetSize(size float32) {
	c.size = size
	if c.sizeIsScale {
		c.Transform().SetScalef(size, size)
	}
	c.UpdateResolution()
}

func (c *Camera) Size() float32 {
	return c.size
}

func (c *Camera) InvertedMatrix() Matrix {
	if c.sizeIsScale {
		c.Transform().updateMatrix()
		return c.Transform().matrix.Invert()
	}

	m := Identity()
	pos := c.Transform().WorldPosition()
	r := c.Transform().WorldRotation()
	m.Rotate(r.X, 1, 0, 0)
	m.Rotate(r.Y, 0, 1, 0)
	m.Rotate(r.Z, 0, 0, -1)
	m.Translate(pos.X, pos.Y, pos.Z)

	return m.Invert()

}

func (c *Camera) InsideScreen(ratio float32, position Vector, scale Vector) bool {
	cameraPos := c.Transform().WorldPosition()

	bigScale := scale.X * ratio
	if scale.Y > bigScale {
		bigScale = scale.Y
	}
	//bigScale = -bigScale

	r := c.rect
	r.Min = r.Min.Mul2(c.size)
	r.Max = r.Max.Mul2(c.size)
	r.Min = r.Min.Add(cameraPos)
	r.Max = r.Max.Add(cameraPos)

	r2 := Rect{}
	r2.Min.X = (-float32(bigScale) / 2) + position.X
	r2.Max.X = (float32(bigScale) / 2) + position.X
	r2.Min.Y = (-float32(bigScale) / 2) + position.Y
	r2.Max.Y = (float32(bigScale) / 2) + position.Y

	return r.Overlaps(r2)
}

func (c *Camera) UpdateResolution() {

	c.rect.Min.X = -float32(Width) / 2
	c.rect.Max.X = float32(Width) / 2
	c.rect.Min.Y = -float32(Height) / 2
	c.rect.Max.Y = float32(Height) / 2

	if c.sizeIsScale {
		c.realRect.Min.X = c.center.X - (c.center.X - c.rect.Min.X)
		c.realRect.Max.X = c.center.X - (c.center.X - c.rect.Max.X)
		c.realRect.Min.Y = c.center.Y - (c.center.Y - c.rect.Min.Y)
		c.realRect.Max.Y = c.center.Y - (c.center.Y - c.rect.Max.Y)
	} else {
		c.realRect.Min.X = c.center.X - (c.center.X-c.rect.Min.X)*c.size
		c.realRect.Max.X = c.center.X - (c.center.X-c.rect.Max.X)*c.size
		c.realRect.Min.Y = c.center.Y - (c.center.Y-c.rect.Min.Y)*c.size
		c.realRect.Max.Y = c.center.Y - (c.center.Y-c.rect.Max.Y)*c.size
	}
	c.Projection.Ortho(c.realRect.Min.X, c.realRect.Max.X, c.realRect.Min.Y, c.realRect.Max.Y, -1000, 1000)

}

func (c *Camera) MouseWorldPosition() Vector {
	v := c.MouseLocalPosition()

	return c.ScreenToWorld(v.X, v.Y)
}

func (c *Camera) MouseLocalPosition() Vector {
	xx, yy := input.MousePosition()
	x, y := float32(xx), float32(yy)
	if c.sizeIsScale {
		x, y = x+c.rect.Min.X, c.rect.Max.Y-y
	} else {
		x, y = (x*c.size)+c.realRect.Min.X, c.realRect.Max.Y-(y*c.size)
	}

	return NewVector2(x, y)
}

func (c *Camera) ScreenToWorld(x, y float32) Vector {

	m := Identity()
	if c.sizeIsScale {
		m.Translate(float32(x), float32(y), 0)
		m = Mul(m, c.Transform().Matrix())
	} else {
		m.Translate(float32(x)*c.size, float32(y)*c.size, 0)
		pos := c.Transform().WorldPosition()
		r := c.Transform().WorldRotation()
		m.Rotate(r.X, 1, 0, 0)
		m.Rotate(r.Y, 0, 1, 0)
		m.Rotate(r.Z, 0, 0, -1)
		m.Translate(pos.X, pos.Y, pos.Z)
	}

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
