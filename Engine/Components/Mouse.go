package Components

import (
	. "../../Engine"
	. "../../Engine/Input"
	//"fmt"
	c "chipmunk"
	. "chipmunk/vect"
	//"glfw"
	)

type Mouse struct {
	BaseComponent
}

func NewMouse() *Mouse {
	return &Mouse{NewComponent()}
} 

func (m *Mouse)OnComponentBind(binded *GameObject) {
	p := NewPhysics2(false, c.NewCircle(Vect{0.5,0},Float(1)))
	p.Body.SetMass(10)
	p.Body.SetMoment(Inf)
	p.Shape.IsSensor = false 
	p.Body.IgnoreGravity = true
	binded.AddComponent(p)
}  

//	Height = 480
//	Width  = 640 


func (m *Mouse) Update() {
	x,y := MousePosition()
	x,y = x,480-y
	m.Transform().SetPosition(NewVector2(float32(x),float32(y)))
	//m.GameObject().Physics.Body.Velocity = Vect{}
	//m.GameObject().Physics.Body.Force = Vect{}
}


func (m *Mouse) OnCollision(c Collision) {
	a := c.ColliderA.Components()
	b := c.ColliderB.Components()
	for _,c := range a {
		if f,ok := c.(onMouseHoverComponent); ok {
			f.OnMouseHover()
		}
	}
	for _,c := range b {
		if f,ok := c.(onMouseHoverComponent); ok {
			f.OnMouseHover()
		}
	}
}



func (m *Mouse) OnCollisionEnter(c Collision) {
	//fmt.Println("Enter",c.ColliderA.Name(), c.ColliderB.Name())
	a := c.ColliderA.Components()
	b := c.ColliderB.Components()
	for _,c := range a {
		if f,ok := c.(onMouseEnterComponent); ok {
			f.OnMouseEnter()
		}
	}
	for _,c := range b {
		if f,ok := c.(onMouseEnterComponent); ok {
			f.OnMouseEnter()
		}
	}
}



func (m *Mouse) OnCollisionExit(c Collision) {
	//fmt.Println("Exit",c.ColliderA.Name(), c.ColliderB.Name())
	a := c.ColliderA.Components()
	b := c.ColliderB.Components()
	for _,c := range a {
		if f,ok := c.(onMouseExitComponent); ok {
			f.OnMouseExit()
		}
	}
	for _,c := range b {
		if f,ok := c.(onMouseExitComponent); ok {
			f.OnMouseExit()
		}
	}
}

type onMouseHoverComponent interface {
	 OnMouseHover()
}

type onMouseEnterComponent interface {
	 OnMouseEnter()
}

type onMouseExitComponent interface {
	 OnMouseExit()
}

