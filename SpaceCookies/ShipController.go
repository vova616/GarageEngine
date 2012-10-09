package SpaceCookies

import (
	. "github.com/vova616/GarageEngine/Engine"
	//"Engine/Components"
	//"github.com/jteeuwen/glfw"
	"github.com/vova616/GarageEngine/Engine/Input"
	//"log"
	c "github.com/vova616/chipmunk"
	. "github.com/vova616/chipmunk/vect"
	//"fmt"
	"math"
	"math/rand"
	"time"
)

type ShipController struct {
	BaseComponent
	Speed           float32
	RotationSpeed   float32
	Missle          *Missle
	MisslesPosition []Vector
	lastShoot       time.Time
	Destoyable      *Destoyable
	HPBar           *GameObject

	UseMouse bool

	JetFire         *GameObject
	JetFireParent   *GameObject
	JetFirePool     []*GameObject
	JetFirePosition []Vector
}

func NewShipController() *ShipController {
	return &ShipController{NewComponent(), 500000, 250, nil, []Vector{{-28, 10, 0},
		{28, 10, 0}}, time.Now(), nil, nil, false, nil, nil, nil, []Vector{{-0.35, -0.51, 0}, {0.35, -0.51, 0}}}
}

func (sp *ShipController) OnComponentBind(binded *GameObject) {
	sp.GameObject().AddComponent(NewPhysics2(false, c.NewCircle(Vect{0, 0}, 15)))
}

func (sp *ShipController) Start() {
	ph := sp.GameObject().Physics
	ph.Body.SetMass(50)
	ph.Shape.Group = 1
	sp.Destoyable = sp.GameObject().ComponentTypeOfi(sp.Destoyable).(*Destoyable)
	sp.OnHit(nil, nil)

	sp.JetFireParent = NewGameObject("JetFireParent")
	sp.JetFireParent.Transform().SetParent2(sp.GameObject())

	uvJet := IndexUV(atlas, Jet_A)

	if sp.JetFire != nil {
		l := 1
		sp.JetFirePool = make([]*GameObject, l*len(sp.JetFirePosition))

		for i := 0; i < len(sp.JetFirePosition); i++ {
			for j := 0; j < l; j++ {

				jfp := NewGameObject("JetFireParent2")
				jfp.Transform().SetParent2(sp.JetFireParent)
				jfp.Transform().SetPosition(sp.JetFirePosition[i])
				jfp.AddComponent(NewResizeScript(0.2, 0.3, 0.5, 1.0))

				jf := sp.JetFire.Clone()
				jf.Transform().SetParent2(jfp)
				jf.Transform().SetPositionf(0, -((uvJet.Ratio)/2)*jf.Transform().Scale().Y, 0)

				//	jf.Transform().SetWorldScalef(10, 10, 1)

				sp.JetFirePool[(i*l)+j] = jf
			}
		}
	}
	//sp.Physics.Shape.Friction = 0.5
}

func (sp *ShipController) OnHit(enemey *GameObject, damager *DamageDealer) {
	if sp.HPBar != nil && sp.Destoyable != nil {
		hp := (float32(sp.Destoyable.HP) / float32(sp.Destoyable.FullHP)) * 100
		s := sp.HPBar.Transform().Scale()
		s.X = hp
		sp.HPBar.Transform().SetScale(s)
	}
}

func (sp *ShipController) OnDie() {
	for i := 0; i < 20; i++ {
		n := Explosion.Clone()
		n.Transform().SetParent2(GameSceneGeneral.Layer1)
		n.Transform().SetWorldPosition(sp.Transform().WorldPosition())
		s := n.Transform().Scale()
		n.Transform().SetScale(s.Mul2(rand.Float32() * 8))
		n.AddComponent(NewPhysics(false, 1, 1))

		n.Transform().SetRotationf(0, 0, rand.Float32()*360)
		rot := n.Transform().Rotation2D()
		n.Physics.Body.SetVelocity(-rot.X*100, -rot.Y*100)

		n.Physics.Body.SetMass(1)
		n.Physics.Shape.Group = 1
		n.Physics.Shape.IsSensor = true
	}
	sp.GameObject().Destroy()
}

func (sp *ShipController) Shoot() {
	if sp.Missle != nil {
		s := sp.Transform().Direction2D(Up)
		a := sp.Transform().Rotation()
		//scale := sp.Transform().Scale()
		for _, pos := range sp.MisslesPosition {
			p := sp.Transform().WorldPosition()
			_ = a
			m := Identity()
			//m.Scale(scale.X, scale.Y, scale.Z)
			m.Translate(pos.X, pos.Y, pos.Z)
			m.Rotate(a.Z, 0, 0, -1)
			m.Translate(p.X, p.Y, p.Z)
			p = m.Translation()

			nfire := sp.Missle.GameObject().Clone()
			nfire.Transform().SetParent2(GameSceneGeneral.Layer3)
			nfire.Transform().SetWorldPosition(p)
			nfire.Physics.Body.IgnoreGravity = true
			nfire.Physics.Body.SetMass(0.1)
			nfire.Tag = MissleTag

			v := sp.GameObject().Physics.Body.Velocity()

			nfire.Physics.Body.SetVelocity(float32(v.X), float32(v.Y))
			nfire.Physics.Body.AddForce(s.X*3000, s.Y*3000)

			nfire.Physics.Shape.Group = 1
			nfire.Physics.Body.SetMoment(Inf)
			nfire.Transform().SetRotation(sp.Transform().Rotation())
		}
	}
}

func (sp *ShipController) Update() {

	r2 := sp.Transform().Direction2D(Up)
	r3 := sp.Transform().Direction2D(Left)
	ph := sp.GameObject().Physics
	rx, ry := r2.X*DeltaTime(), r2.Y*DeltaTime()
	rsx, rsy := r3.X*DeltaTime(), r3.Y*DeltaTime()

	jet := false

	if Input.KeyDown('W') {
		ph.Body.AddForce(sp.Speed*rx, sp.Speed*ry)
		jet = true
	}

	if Input.KeyDown('S') {
		ph.Body.AddForce(-sp.Speed*rx, -sp.Speed*ry)
		jet = true
	}

	if jet {
		sp.JetFireParent.SetActiveRecursive(true)
	} else {
		sp.JetFireParent.SetActiveRecursive(false)
	}

	rotSpeed := sp.RotationSpeed
	if Input.KeyDown(Input.KeyLshift) {
		rotSpeed = 100
	}

	if sp.UseMouse {
		v := GetScene().SceneBase().Camera.MouseWorldPosition()
		v = v.Sub(sp.Transform().WorldPosition())
		v.Normalize()
		angle := float32(math.Atan2(float64(v.Y), float64(v.X))) * DegreeConst
		sp.Transform().SetRotationf(0, 0, float32(int(angle-90)))

		if Input.KeyDown('D') || Input.KeyDown('E') {
			ph.Body.SetAngularVelocity(0)
			ph.Body.SetTorque(0)
			ph.Body.AddForce(-sp.Speed*rsx, -sp.Speed*rsy)
		}
		if Input.KeyDown('A') || Input.KeyDown('Q') {
			ph.Body.SetAngularVelocity(0)
			ph.Body.SetTorque(0)
			ph.Body.AddForce(sp.Speed*rsx, sp.Speed*rsy)
		}
	} else {
		r := sp.Transform().Rotation()
		if Input.KeyDown('D') {
			ph.Body.SetAngularVelocity(0)
			ph.Body.SetTorque(0)
			sp.Transform().SetRotationf(0, 0, r.Z-rotSpeed*DeltaTime())
		}
		if Input.KeyDown('A') {
			ph.Body.SetAngularVelocity(0)
			ph.Body.SetTorque(0)
			sp.Transform().SetRotationf(0, 0, r.Z+rotSpeed*DeltaTime())
		}

		if Input.KeyDown('E') {
			ph.Body.SetAngularVelocity(0)
			ph.Body.SetTorque(0)
			ph.Body.AddForce(-sp.Speed*rsx, -sp.Speed*rsy)
		}
		if Input.KeyDown('Q') {
			ph.Body.SetAngularVelocity(0)
			ph.Body.SetTorque(0)
			ph.Body.AddForce(sp.Speed*rsx, sp.Speed*rsy)
		}
	}

	if Input.MouseDown(Input.MouseLeft) {
		if time.Now().After(sp.lastShoot) {
			sp.Shoot()
			sp.lastShoot = time.Now().Add(time.Millisecond * 200)
		}
	}

	if Input.KeyPress('P') {
		EnablePhysics = !EnablePhysics
	}
	if Input.KeyPress('T') {
		sp.UseMouse = !sp.UseMouse
	}
}

func (sp *ShipController) LateUpdate() {
	if GameSceneGeneral.SceneData.Camera.GameObject() != nil {
		GameSceneGeneral.SceneData.Camera.Transform().SetPosition(NewVector3(sp.Transform().Position().X-float32(Width/2), sp.Transform().Position().Y-float32(Height/2), 0))
	}
}
