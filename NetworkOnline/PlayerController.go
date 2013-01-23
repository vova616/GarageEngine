package NetworkOnline

import (
	"github.com/vova616/GarageEngine/Engine"
	//"Engine/Components"
	"github.com/go-gl/glfw"
	"github.com/vova616/GarageEngine/Engine/Input"
	//"log"
	"github.com/vova616/chipmunk"
	"github.com/vova616/chipmunk/vect"
	//"fmt"
	//"time"
)

type PlayerController struct {
	Engine.BaseComponent
	Speed     float32
	JumpSpeed float32
	Physics   *Engine.Physics
	state     int
	Fire      *Engine.GameObject

	Floor *Engine.GameObject
	Fires []*Engine.GameObject
}

func (sp *PlayerController) OnCollisionEnter(arbiter Engine.Arbiter) bool {
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
func (sp *PlayerController) OnCollisionExit(arbiter Engine.Arbiter) {
	println("Exit " + arbiter.GameObjectB().Name())
	if arbiter.GameObjectB() == sp.Floor {
		sp.Floor = nil
	}
}

func NewPlayerController() *PlayerController {
	return &PlayerController{Engine.NewComponent(), 10, 20000, nil, -1, nil, nil, make([]*Engine.GameObject, 0)}
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
			Engine.CoSleep(3)
			sp.Shoot()
		}
	}

	as := Engine.StartCoroutine(autoShoot)

	fastShoot := func() {
		Engine.CoSleep(3)
		Engine.CoYieldCoroutine(as)
		for i := 0; i < 10; i++ {
			Engine.CoYieldSkip()
			Engine.CoYieldSkip()
			Engine.CoYieldSkip()
			sp.Shoot()
		}
		sp.TestCoroutines()
	}

	Engine.StartCoroutine(fastShoot)
}

func (sp *PlayerController) Shoot() {
	if sp.Fire != nil {
		nfire := sp.Fire.Clone()
		sp.Fires = append(sp.Fires, nfire)
		nfire.Transform().SetParent2(GameSceneGeneral.Layer1)
		nfire.Transform().SetWorldPosition(sp.Transform().WorldPosition())
		nfire.AddComponent(Engine.NewPhysics2(false, chipmunk.NewCircle(vect.Vect{0, 0}, 20)))
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
		nfire.Physics.Body.SetMoment(Engine.Inf)

	}
}

func (sp *PlayerController) Update() {
	if Input.KeyPress(glfw.KeyUp) {
		if sp.Floor != nil {
			sp.Physics.Body.AddForce(0, sp.JumpSpeed)
		} else {
			sp.Shoot()
		}
	}

	tState := 0

	if Input.KeyDown(glfw.KeyLeft) {
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
	if Input.KeyDown(glfw.KeyRight) {
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

	if Input.KeyPress(glfw.KeySpace) {
		sp.Shoot()
	}

	if Input.KeyPress('P') {

		Engine.EnablePhysics = !Engine.EnablePhysics
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
	//GameSceneGeneral.SceneData.Camera.Transform().SetPosition(NewVector3(300-sp.Transform().Position().X, 0, 0))
}
