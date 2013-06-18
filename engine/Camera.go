package engine

import (
	"errors"
	"github.com/go-gl/gl"
	"github.com/vova616/GarageEngine/engine/input"
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
	realRect    Rect   //
	center      Vector //Center of scale
	sizeIsScale bool   //Size is also the scale
	autoScale   bool   //Scale together with window size

	clearColor Color
}

func NewCamera() *Camera {
	c := &Camera{BaseComponent: NewComponent(), Projection: NewIdentity(), size: 1, sizeIsScale: true, autoScale: true}

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
	c.LateUpdate()
}

func (c *Camera) Clear() {
	gl.ClearColor(gl.GLclampf(c.clearColor.R), gl.GLclampf(c.clearColor.G), gl.GLclampf(c.clearColor.B), gl.GLclampf(c.clearColor.A))
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

func (c *Camera) LateUpdate() {
	if c.sizeIsScale {
		s := c.Transform().Scale()
		if s.X != c.size || s.Y != c.size {
			c.Transform().SetScalef(c.size, c.size)
		}
	}
}

//Setting the size and also scale if sizeIsScale is true
func (c *Camera) SetSize(size float32) error {
	if size <= 0 {
		return errors.New("Camera size cannot be 0 or less")
	}
	c.size = size
	if c.sizeIsScale {
		c.Transform().SetScalef(size, size)
	}
	c.UpdateResolution()
	return nil
}

//Size of the camera screen
func (c *Camera) Size() float32 {
	return c.size
}

//InvertedMatrix of the camera, this is needed because we will optimize it someday
func (c *Camera) InvertedMatrix() Matrix {
	if c.sizeIsScale {
		return c.Transform().InvertedMatrix()
	} else {
		x := c.Matrix()
		return x.Invert()
	}
}

//Returns Screen Size
func (c *Camera) ScreenSize() (Width, Height float32) {
	if c.sizeIsScale {
		return (c.realRect.Max.X - c.realRect.Min.X) * c.size, (c.realRect.Max.Y - c.realRect.Min.Y) * c.size
	}
	return c.realRect.Max.X - c.realRect.Min.X, c.realRect.Max.Y - c.realRect.Min.Y
}

//Screen resolution
func (c *Camera) ScreenResolution() (Width, Height float32) {
	return c.rect.Max.X - c.rect.Min.X, c.rect.Max.Y - c.rect.Min.Y
}

//Matrix of the camera, this is needed because sometimes we control the matrix
func (c *Camera) Matrix() Matrix {
	if c.sizeIsScale {
		c.Transform().updateMatrix()
		return c.Transform().matrix
	}

	m := Identity()
	pos := c.Transform().WorldPosition()
	r := c.Transform().WorldRotation()
	m.Rotate(r.X, 1, 0, 0)
	m.Rotate(r.Y, 0, 1, 0)
	m.Rotate(r.Z, 0, 0, -1)
	m.Translate(pos.X, pos.Y, pos.Z)
	return m
}

//Checks if box is in the screen
func (c *Camera) InsideScreen(ratio float32, position Vector, scale Vector) bool {
	cameraPos := c.Transform().WorldPosition()

	bigScale := scale.X * ratio
	if scale.Y > bigScale {
		bigScale = scale.Y
	}
	//bigScale = -bigScale

	r := c.realRect
	if c.sizeIsScale {
		r.Min = r.Min.Mul2(c.size)
		r.Max = r.Max.Mul2(c.size)
	}
	r.Min = r.Min.Add(cameraPos)
	r.Max = r.Max.Add(cameraPos)

	r2 := Rect{}
	r2.Min.X = (-float32(bigScale) / 2) + position.X
	r2.Max.X = (float32(bigScale) / 2) + position.X
	r2.Min.Y = (-float32(bigScale) / 2) + position.Y
	r2.Max.Y = (float32(bigScale) / 2) + position.Y

	return r.Overlaps(r2)
}

//Updates the Projection
func (c *Camera) UpdateResolution() {
	if c.autoScale {
		c.rect.Min.X = -float32(Width) / 2
		c.rect.Max.X = float32(Width) / 2
		c.rect.Min.Y = -float32(Height) / 2
		c.rect.Max.Y = float32(Height) / 2
	}

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
	c.Projection.Ortho(c.realRect.Min.X, c.realRect.Max.X, c.realRect.Min.Y, c.realRect.Max.Y, -100000, 100000)

}

//Mouse world position
func (c *Camera) MouseWorldPosition() Vector {
	v := c.MouseLocalPosition()
	return c.ScreenToWorld(v.X, v.Y)
}

//Mouse local position
func (c *Camera) MouseLocalPosition() Vector {
	xx, yy := input.MousePosition()
	x, y := float32(xx), float32(yy)
	if c.sizeIsScale {
		x, y = x+c.realRect.Min.X, c.realRect.Max.Y-y
	} else {
		x, y = (x*c.size)+c.realRect.Min.X, c.realRect.Max.Y-(y*c.size)
	}

	return NewVector2(x, y)
}

//Takes a point on the screen and turns it into point on world
func (c *Camera) ScreenToWorld(x, y float32) Vector {
	m := Identity()
	if c.sizeIsScale {
		m.Translate(float32(x), float32(y), 0)
		m.Mul(c.Transform().Matrix())
	} else {
		m.Translate(float32(x)*c.size, float32(y)*c.size, 0)
		pos := c.Transform().WorldPosition()
		r := c.Transform().WorldRotation()
		m.RotateXYZ(r.X, r.Y, r.Z)
		m.Translate(pos.X, pos.Y, pos.Z)
	}
	return m.Translation()
}

//Forcing all the objects to render
func (c *Camera) Render() {
	s := GetScene()
	if s != nil {
		tcam := s.SceneBase().Camera
		s.SceneBase().Camera = c
		depthMap.Iter(drawGameObject)
		s.SceneBase().Camera = tcam
	}
}
