package NetworkOnline

import (
	. "../Engine"
    //"Engine/Components"
	"../Engine/Input"
	//"log"
)

type PlayerController struct {
	BaseComponent
	Speed float64
	JumpSpeed float64
	Physics *Physics
	state	int
	Fire *GameObject
}

func NewPlayerController() *PlayerController {
	return &PlayerController{NewComponent(),10,20000,nil, -1,nil}
}

func (sp *PlayerController) Start() {
	if sp.GameObject().Physics == nil {
		return
	}
	sp.Physics = sp.GameObject().Physics
	sp.Physics.Body.SetMass(1)
	//sp.Physics.Shape.Friction = 0.5
}

func (sp *PlayerController) Update() {
	if Input.KeyPress('W') {
		//sp.Physics.Body.Force.Y += sp.JumpSpeed
	}
	
	tState := 0
	
	if Input.KeyDown('A') {		 
		//sp.Physics.Body.Velocity.X -= sp.Speed
		if sp.state != 1 {
			sp.GameObject().Sprite.SetAnimation("walk")
		}
		s := sp.Transform().Rotation()
		s.Y = 180
		sp.Transform().SetRotation(s)
		sp.state = 1
		tState = 1
		//sp.Physics.Shape.Friction = 0
	}
	if Input.KeyDown('D') {
		//sp.Physics.Body.Velocity.X += sp.Speed
		if sp.state != 1 {
			sp.GameObject().Sprite.SetAnimation("walk")
		}
		s := sp.Transform().Rotation()
		s.Y = 0
		sp.Transform().SetRotation(s)
		sp.state = 1
		tState = 1
		//sp.Physics.Shape.Friction = 0
	} 
	
	if Input.KeyDown('E') {
		//sp.Physics.Body.Torque += 1
	}
	if Input.KeyDown('Q') {
		//sp.Physics.Body.Torque -= 1
	}
	
	if Input.KeyPress('F') {
		if sp.Fire != nil {
			nfire := sp.Fire.Clone()
			nfire.Transform().SetParent2(GameSceneGeneral.Layer1)
			nfire.Transform().SetWorldPosition(sp.Transform().WorldPosition())
			nfire.AddComponent(NewPhysics(false))
			nfire.Physics.Body.IgnoreGravity = true
			nfire.Physics.Body.SetMass(20)
			s := sp.Transform().Rotation()
			s2 := nfire.Transform().Rotation()
			s2.Y = s.Y
			if s2.Y == 180 {
				s2.Z = 90
				nfire.Transform().Translate2(-60,0,0)
				//nfire.Physics.Body.Velocity.X = -350
			} else {
				s2.Z = -90
				nfire.Transform().Translate2(60,0,0)
				//nfire.Physics.Body.Velocity.X = 350
			}
			nfire.Physics.Body.SetMoment(0)
			nfire.Transform().SetRotation(s2)
		}
	} 
	
	if tState != 1 {
		if sp.state != 0 {
			sp.GameObject().Sprite.SetAnimation("stand")
		}
		//sp.Physics.Shape.Friction = 0.5
		sp.state = 0
	}
	
	
	//if sp.Physics.Body.Velocity.X > 200 {
	//	sp.Physics.Body.Velocity.X = 200
	//} else if sp.Physics.Body.Velocity.X < -200 {
	//	sp.Physics.Body.Velocity.X = -200
	//}
	
	GameSceneGeneral.Camera.Transform().SetPosition(NewVector3(200-sp.Transform().Position().X,0,0))
}
