package SpaceCookies

import (
	"github.com/vova616/GarageEngine/Engine"
	//"Engine/Components"
	//"github.com/jteeuwen/glfw"
	"github.com/vova616/GarageEngine/Engine/Input"
	//"log"
	"github.com/vova616/chipmunk"
	"github.com/vova616/chipmunk/vect"
	//"fmt"
	"math"
	"math/rand"
	"time"
)

type ShipController struct {
	Engine.BaseComponent
	Speed            float32
	RotationSpeed    float32
	Missle           *Missle
	MisslesPosition  []Engine.Vector
	MisslesDirection [][]Engine.Vector
	MissleLevel      int
	MaxMissleLevel   int
	lastShoot        time.Time
	Destoyable       *Destoyable
	HPBar            *Engine.GameObject

	UseMouse bool

	JetFire         *Engine.GameObject
	JetFireParent   *Engine.GameObject
	JetFirePool     []*ResizeScript
	JetFirePosition []Engine.Vector
}

func NewShipController() *ShipController {

	misslesDirection := [][]Engine.Vector{{{0, 1, 0}, {0, 1, 0}},
		{{-0.2, 1, 0}, {0.2, 1, 0}, {0, 1, 0}},
		{{-0.2, 1, 0}, {0.2, 1, 0}, {0, 1, 0}, {-0.2, 1, 0}, {0.2, 1, 0}}}

	misslePositions := []Engine.Vector{{-28, 10, 0}, {28, 10, 0}, {0, 20, 0}, {-28, 40, 0}, {28, 40, 0}}

	return &ShipController{Engine.NewComponent(), 500000, 250, nil, misslePositions, misslesDirection, 0, len(misslesDirection) - 1,
		time.Now(), nil, nil, false, nil, nil, nil, []Engine.Vector{{-0.1, -0.51, 0}, {0.1, -0.51, 0}}}
}

func (sp *ShipController) OnComponentBind(binded *Engine.GameObject) {
	sp.GameObject().AddComponent(Engine.NewPhysics2(false, chipmunk.NewCircle(vect.Vect{0, 0}, 15)))
}

func (sp *ShipController) Start() {
	ph := sp.GameObject().Physics
	ph.Body.SetMass(50)
	ph.Shape.Group = 1
	sp.Destoyable = sp.GameObject().ComponentTypeOfi(sp.Destoyable).(*Destoyable)
	sp.OnHit(nil, nil)

	sp.JetFireParent = Engine.NewGameObject("JetFireParent")
	sp.JetFireParent.Transform().SetParent2(sp.GameObject())

	uvJet := Engine.IndexUV(atlas, Jet_A)

	if sp.JetFire != nil {
		l := 1
		sp.JetFirePool = make([]*ResizeScript, l*len(sp.JetFirePosition))

		for i := 0; i < len(sp.JetFirePosition); i++ {
			for j := 0; j < l; j++ {

				jfp := Engine.NewGameObject("JetFireParent2")
				jfp.Transform().SetParent2(sp.JetFireParent)
				jfp.Transform().SetPosition(sp.JetFirePosition[i])
				rz := NewResizeScript(0.0, 0.0, 0.0, 0.0, 0.0, 0.0)
				jfp.AddComponent(rz)

				jf := sp.JetFire.Clone()
				jf.Transform().SetParent2(jfp)
				jf.Transform().SetPositionf(0, -((uvJet.Ratio)/2)*jf.Transform().Scale().Y, 0)

				//	jf.Transform().SetWorldScalef(10, 10, 1)

				sp.JetFirePool[(i*l)+j] = rz
			}
		}
	}
	//sp.Physics.Shape.Friction = 0.5
}

func (sp *ShipController) OnHit(enemey *Engine.GameObject, damager *DamageDealer) {
	if sp.HPBar != nil && sp.Destoyable != nil {
		hp := (float32(sp.Destoyable.HP) / float32(sp.Destoyable.FullHP)) * 100
		s := sp.HPBar.Transform().Scale()
		s.X = hp
		sp.HPBar.Transform().SetScale(s)
	}
}

func (sp *ShipController) OnDie(byTimer bool) {
	for i := 0; i < 20; i++ {
		n := Explosion.Clone()
		n.Transform().SetParent2(GameSceneGeneral.Layer1)
		n.Transform().SetWorldPosition(sp.Transform().WorldPosition())
		s := n.Transform().Scale()
		n.Transform().SetScale(s.Mul2(rand.Float32() * 8))
		n.AddComponent(Engine.NewPhysics(false, 1, 1))

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

		a := sp.Transform().Rotation()
		//scale := sp.Transform().Scale()
		for i := 0; i < len(sp.MisslesDirection[sp.MissleLevel]); i++ {
			pos := sp.MisslesPosition[i]
			s := sp.Transform().Direction2D(sp.MisslesDirection[sp.MissleLevel][i])
			//s = s.Mul(sp.MisslesDirection[i])

			p := sp.Transform().WorldPosition()
			_ = a
			m := Engine.Identity()
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
			nfire.Physics.Body.SetMoment(Engine.Inf)
			nfire.Transform().SetRotation(sp.Transform().Rotation())
		}
	}
}

func (sp *ShipController) Update() {

	r2 := sp.Transform().Direction2D(Engine.Up)
	r3 := sp.Transform().Direction2D(Engine.Left)
	ph := sp.GameObject().Physics
	rx, ry := r2.X*Engine.DeltaTime(), r2.Y*Engine.DeltaTime()
	rsx, rsy := r3.X*Engine.DeltaTime(), r3.Y*Engine.DeltaTime()

	jet := false
	back := false

	if Input.KeyDown('W') {
		ph.Body.AddForce(sp.Speed*rx, sp.Speed*ry)
		jet = true
	}

	if Input.KeyDown('S') {
		ph.Body.AddForce(-sp.Speed*rx, -sp.Speed*ry)
		jet = true
		back = true
	}

	if jet {
		for _, resize := range sp.JetFirePool {
			if back {
				resize.SetValues(0.1, 0.1, 0.2, 0.0, 0.2, 0.3)
			} else {
				resize.SetValues(0.2, 0.2, 0.3, 0.0, 0.6, 0.8)
			}
		}
		if !sp.JetFireParent.IsActive() {
			sp.JetFireParent.SetActiveRecursive(true)
			for _, resize := range sp.JetFirePool {
				resize.State = 0
				if back {
					resize.SetValues(0.1, 0.1, 0.2, 0.0, 0.2, 0.3)
				} else {
					resize.SetValues(0.2, 0.2, 0.3, 0.0, 0.6, 0.8)
				}
			}
		}
	} else {
		sp.JetFireParent.SetActiveRecursive(false)
	}

	rotSpeed := sp.RotationSpeed
	if Input.KeyDown(Input.KeyLshift) {
		rotSpeed = 100
	}

	if sp.UseMouse {
		v := Engine.GetScene().SceneBase().Camera.MouseWorldPosition()
		v = v.Sub(sp.Transform().WorldPosition())
		v.Normalize()
		angle := float32(math.Atan2(float64(v.Y), float64(v.X))) * Engine.DegreeConst
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
			sp.Transform().SetRotationf(0, 0, r.Z-rotSpeed*Engine.DeltaTime())
		}
		if Input.KeyDown('A') {
			ph.Body.SetAngularVelocity(0)
			ph.Body.SetTorque(0)
			sp.Transform().SetRotationf(0, 0, r.Z+rotSpeed*Engine.DeltaTime())
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
		Engine.EnablePhysics = !Engine.EnablePhysics
	}
	if Input.KeyPress('T') {
		sp.UseMouse = !sp.UseMouse
	}
}

func (sp *ShipController) LateUpdate() {
	if GameSceneGeneral.SceneData.Camera.GameObject() != nil {
		GameSceneGeneral.SceneData.Camera.Transform().SetPositionf(sp.Transform().Position().X-float32(Engine.Width/2), sp.Transform().Position().Y-float32(Engine.Height/2), 0)
	}
}
