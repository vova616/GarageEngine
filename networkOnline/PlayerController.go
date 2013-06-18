package networkOnline

import (
	"github.com/vova616/GarageEngine/engine"
	"github.com/vova616/GarageEngine/engine/cr"
	//"Engine/components"
	"github.com/go-gl/glfw"
	"github.com/vova616/GarageEngine/engine/input"
	//"log"
	"github.com/vova616/chipmunk"
	"github.com/vova616/chipmunk/vect"
	//"fmt"
	//"time"
)

type PlayerController struct {
	engine.BaseComponent
	Speed     float32
	JumpSpeed float32
	Physics   *engine.Physics
	state     int
	Fire      *engine.GameObject

	Floor *engine.GameObject
	Fires []*engine.GameObject
}

func (sp *PlayerController) OnCollisionEnter(arbiter engine.Arbiter) bool {
	//println("Enter " + arbiter.GameObjectB().Name())
	count := 0
	for _, con := range arbiter.Contacts {
		//println(arbiter.Normal(con).Y)
		if -arbiter.Normal(con).Y > 0.9 {
			count++
		}
	}
	if count >= 2 {
		sp.Floor = arbiter.GameObjectB()
	}

	return true
}
func (sp *PlayerController) OnCollisionExit(arbiter engine.Arbiter) {
	if arbiter.GameObjectB() == sp.Floor {
		sp.Floor = nil
	}
}

func NewPlayerController() *PlayerController {
	return &PlayerController{engine.NewComponent(), 10, 20000, nil, -1, nil, nil, make([]*engine.GameObject, 0)}
}

func (sp *PlayerController) Start() {
	if sp.GameObject().Physics == nil {
		return
	}
	sp.Physics = sp.GameObject().Physics
	sp.Physics.Body.SetMass(1)
	sp.Physics.Shape.Group = 1
	//sp.Physics.Shape.Friction = 0.5

	//sp.TestCoroutines()
}

func (sp *PlayerController) TestCoroutines() {
	autoShoot := func() {
		for i := 0; i < 3; i++ {
			cr.Sleep(3)
			sp.Shoot()
		}
	}

	as := cr.Start(autoShoot)

	fastShoot := func() {
		cr.Sleep(3)
		cr.YieldCoroutine(as)
		for i := 0; i < 10; i++ {
			cr.Skip()
			cr.Skip()
			cr.Skip()
			sp.Shoot()
		}
		sp.TestCoroutines()
	}

	cr.Start(fastShoot)
}

func (sp *PlayerController) Shoot() {
	if sp.Fire != nil {
		nfire := sp.Fire.Clone()
		sp.Fires = append(sp.Fires, nfire)
		nfire.Transform().SetParent2(GameSceneGeneral.Layer1)
		nfire.Transform().SetWorldPosition(sp.Transform().WorldPosition())
		nfire.AddComponent(engine.NewPhysicsShape(false, chipmunk.NewCircle(vect.Vect{0, 0}, 20)))
		nfire.Physics.Body.IgnoreGravity = true
		nfire.Physics.Body.SetMass(200)
		s := sp.Transform().Rotation()
		if s.Y == 180 {
			nfire.Transform().SetRotationf(90)
			nfire.Transform().Translatef(-20, 0)
			nfire.Physics.Body.SetVelocity(-550, 0)
		} else {
			nfire.Transform().SetRotationf(90)
			nfire.Transform().Translatef(20, 0)
			nfire.Physics.Body.SetVelocity(550, 0)
		}
		nfire.Physics.Shape.Group = 1
		nfire.Physics.Body.SetMoment(engine.Inf)

	}
}

func (sp *PlayerController) Update() {
	if input.KeyPress(glfw.KeyUp) {
		if sp.Floor != nil {
			sp.Physics.Body.AddForce(0, sp.JumpSpeed)
		} else {
			sp.Shoot()
		}
	}

	tState := 0

	if input.KeyDown(glfw.KeyLeft) {
		sp.Physics.Body.AddVelocity(-sp.Speed, 0)
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
	if input.KeyDown(glfw.KeyRight) {
		sp.Physics.Body.AddVelocity(sp.Speed, 0)
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

	if input.KeyPress(glfw.KeySpace) {
		sp.Shoot()
	}

	if input.KeyPress('P') {

		engine.EnablePhysics = !engine.EnablePhysics
	}

	if tState != 1 {
		if sp.state != 0 {
			sp.GameObject().Sprite.SetAnimation("stand")
		}
		//sp.Physics.Shape.Friction = 0.5
		sp.state = 0
	}

	v := sp.Physics.Body.Velocity()

	if v.X > 200 {
		sp.Physics.Body.SetVelocity(200, float32(v.Y))
	} else if v.X < -200 {
		sp.Physics.Body.SetVelocity(-200, float32(v.Y))
	}

	for i := 0; i < len(sp.Fires); i++ {
		fire := sp.Fires[i]
		if fire.Transform().Rotation().Z <= -80 && fire.Physics.Body.Velocity().X >= -1 {
			fire.Destroy()
			sp.Fires = append(sp.Fires[:i], sp.Fires[i+1:]...)
			i--
		} else if fire.Transform().Rotation().Z >= 80 && fire.Physics.Body.Velocity().X <= 1 {
			fire.Destroy()
			sp.Fires = append(sp.Fires[:i], sp.Fires[i+1:]...)
			i--
		}
	}
}

func (sp *PlayerController) LateUpdate() {
	//gameSceneGeneral.SceneData.Camera.Transform().SetPosition(NewVector3(300-sp.Transform().Position().X, 0, 0))
}
