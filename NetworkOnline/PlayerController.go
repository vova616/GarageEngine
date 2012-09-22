package NetworkOnline

import (
	. "github.com/vova616/GarageEngine/Engine"
	//"Engine/Components"
	"github.com/jteeuwen/glfw"
	"github.com/vova616/GarageEngine/Engine/Input"
	//"log"
	c "github.com/vova616/chipmunk"
	. "github.com/vova616/chipmunk/vect"
	//"fmt"
	//"time"
)

type PlayerController struct {
	BaseComponent
	Speed     float32
	JumpSpeed float32
	Physics   *Physics
	state     int
	Fire      *GameObject

	Floor *GameObject
	Fires []*GameObject
}

func NewPlayerController() *PlayerController {
	return &PlayerController{NewComponent(), 10, 20000, nil, -1, nil, nil, make([]*GameObject, 0)}
}

func (sp *PlayerController) Start() {
	if sp.GameObject().Physics == nil {
		return
	}
	sp.Physics = sp.GameObject().Physics
	sp.Physics.Body.SetMass(1)
	sp.Physics.Shape.Group = 1
	//sp.Physics.Shape.Friction = 0.5

	sp.TestCoroutines()
}

func (sp *PlayerController) TestCoroutines() {
	autoShoot := func() {
		for i := 0; i < 3; i++ {
			Wait(3)
			sp.Shoot()
		}
	}

	as := StartCoroutine(autoShoot)

	fastShoot := func() {
		Wait(3)
		YieldCoroutine(as)
		for i := 0; i < 10; i++ {
			YieldSkip()
			YieldSkip()
			YieldSkip()
			sp.Shoot()
		}
		sp.TestCoroutines()
	}

	StartCoroutine(fastShoot)
}

func (sp *PlayerController) Shoot() {
	if sp.Fire != nil {
		nfire := sp.Fire.Clone()
		sp.Fires = append(sp.Fires, nfire)
		nfire.Transform().SetParent2(GameSceneGeneral.Layer1)
		nfire.Transform().SetWorldPosition(sp.Transform().WorldPosition())
		nfire.AddComponent(NewPhysics2(false, c.NewCircle(Vect{0, 0}, Float(20))))
		nfire.Physics.Body.IgnoreGravity = true
		nfire.Physics.Body.SetMass(200)
		s := sp.Transform().Rotation()
		s2 := nfire.Transform().Rotation()
		s2.Y = s.Y
		if s2.Y == 180 {
			s2.Z = 90
			nfire.Transform().Translatef(-20, 0, 0)
			nfire.Physics.Body.SetVelocity(-550, 0)
		} else {
			s2.Z = -90
			nfire.Transform().Translatef(20, 0, 0)
			nfire.Physics.Body.SetVelocity(550, 0)
		}
		nfire.Physics.Shape.Group = 1
		nfire.Physics.Body.SetMoment(Inf)
		nfire.Transform().SetRotation(s2)
	}
}

func (sp *PlayerController) Update() {
	if Input.KeyPress(glfw.KeyUp) {
		if sp.Floor != nil {
			sp.Physics.Body.AddForce(0, sp.JumpSpeed)
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

		EnablePhysics = !EnablePhysics
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
		if fire.Transform().Rotation().Z <= -80 && fire.Physics.Body.Velocity().X <= 1 {
			//fire.Destory()
			//HACK
			fire.Transform().SetWorldPosition(NewVector3(-10000, -1000, -1000))
			sp.Fires = append(sp.Fires[:i], sp.Fires[i+1:]...)
			i--
		} else if fire.Transform().Rotation().Z >= 80 && fire.Physics.Body.Velocity().X >= -1 {
			fire.Transform().SetWorldPosition(NewVector3(-10000, -1000, -1000))
			sp.Fires = append(sp.Fires[:i], sp.Fires[i+1:]...)
			i--
		}
	}
}

func (sp *PlayerController) LateUpdate() {
	//GameSceneGeneral.SceneData.Camera.Transform().SetPosition(NewVector3(300-sp.Transform().Position().X, 0, 0))
}

func (sp *PlayerController) OnCollisionEnter(collision Collision) {
	//sp.IsOnFloor = true
	cons := collision.Data.Contacts
	for _, con := range cons {
		if Dot(con.Normal(), Vect{0, 1}) < 0 {
			sp.Floor = collision.ColliderA
			return
		}
	}
}

func (sp *PlayerController) OnCollisionExit(collision Collision) {
	if collision.ColliderA == sp.Floor {
		sp.Floor = nil
	}
}
