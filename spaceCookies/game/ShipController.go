package game

import (
	"github.com/vova616/GarageEngine/engine"
	//"Engine/components"
	//"github.com/go-gl/glfw"
	"github.com/vova616/GarageEngine/engine/input"
	//"log"
	//"fmt"
	//
	"github.com/vova616/GarageEngine/engine/audio"
	"math"
	"math/rand"
	"time"
)

type ShipController struct {
	engine.BaseComponent `json:"-"`
	Speed                float32
	RotationSpeed        float32
	Missle               *Missle `json:"-"`
	MisslesPosition      []engine.Vector
	MisslesDirection     [][]engine.Vector
	MissleLevel          int
	MaxMissleLevel       int                `json:"-"`
	lastShoot            time.Time          `json:"-"`
	Destoyable           *Destoyable        `json:"-"`
	HPBar                *engine.GameObject `json:"-"`

	UseMouse bool `json:"-"`

	JetFire         *engine.GameObject `json:"-"`
	JetFireParent   *engine.GameObject `json:"-"`
	JetFirePool     []*ResizeScript    `json:"-"`
	JetFirePosition []engine.Vector    `json:"-"`

	FireSource *audio.AudioSource
}

func NewShipController() *ShipController {
	misslesDirection := [][]engine.Vector{{{0, 1, 0}, {0, 1, 0}},
		{{-0.2, 1, 0}, {0.2, 1, 0}, {0, 1, 0}},
		{{-0.2, 1, 0}, {0.2, 1, 0}, {0, 1, 0}, {-0.2, 1, 0}, {0.2, 1, 0}}}

	misslePositions := []engine.Vector{{-28, 10, 0}, {28, 10, 0}, {0, 20, 0}, {-28, 40, 0}, {28, 40, 0}}

	return &ShipController{engine.NewComponent(), 500000, 250, nil, misslePositions, misslesDirection, 0, len(misslesDirection) - 1,
		time.Now(), nil, nil, true, nil, nil, nil, []engine.Vector{{-0.1, -0.51, 0}, {0.1, -0.51, 0}}, nil}
}

func (sp *ShipController) OnComponentAdd() {
	sp.GameObject().AddComponent(engine.NewPhysicsCircle(false))
}

func (sp *ShipController) Start() {
	ph := sp.GameObject().Physics
	ph.Body.SetMass(50)
	ph.Shape.Group = 1
	sp.Destoyable = sp.GameObject().ComponentTypeOf(sp.Destoyable).(*Destoyable)
	sp.OnHit(nil, nil)

	sp.JetFireParent = engine.NewGameObject("JetFireParent")
	sp.JetFireParent.Transform().SetParent2(sp.GameObject())
	sp.JetFireParent.Transform().SetPositionf(0, 0)
	sp.JetFireParent.Transform().SetScalef(1, 1)
	sp.JetFireParent.Transform().SetRotationf(0)
	uvJet := engine.IndexUV(atlas, Jet_A)

	if sp.JetFire != nil {
		l := 1
		sp.JetFirePool = make([]*ResizeScript, l*len(sp.JetFirePosition))

		for i := 0; i < len(sp.JetFirePosition); i++ {
			for j := 0; j < l; j++ {

				jfp := engine.NewGameObject("JetFireParent2")
				jfp.Transform().SetParent2(sp.JetFireParent)
				jfp.Transform().SetPosition(sp.JetFirePosition[i])
				jfp.Transform().SetRotationf(0)
				rz := NewResizeScript(0.0, 0.0, 0.0, 0.0, 0.0, 0.0)
				jfp.AddComponent(rz)

				jf := sp.JetFire.Clone()
				jf.Transform().SetParent2(jfp)
				jf.Transform().SetPositionf(0, -((uvJet.Ratio)/2)*jf.Transform().Scale().Y)
				jf.Transform().SetRotationf(0)
				//	jf.Transform().SetWorldScalef(10, 10, 1)

				sp.JetFirePool[(i*l)+j] = rz
			}
		}
	}

	myPos := engine.Vector{sp.Transform().Position().X - float32(engine.Width/2), sp.Transform().Position().Y - float32(engine.Height/2), 0}
	GameSceneGeneral.SceneData.Camera.Transform().SetPosition(myPos)
	sp.MaxMissleLevel = len(sp.MisslesDirection) - 1
	//sp.Physics.Shape.Friction = 0.5
}

func (sp *ShipController) OnHit(enemey *engine.GameObject, damager *DamageDealer) {
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
		n.AddComponent(engine.NewPhysics(false))

		n.Transform().SetRotationf(rand.Float32() * 360)
		rot := n.Transform().Direction()
		n.Physics.Body.SetVelocity(-rot.X*100, -rot.Y*100)

		n.Physics.Body.SetMass(1)
		n.Physics.Shape.Group = 1
		n.Physics.Shape.IsSensor = true
	}
	sp.GameObject().Destroy()
}

func (sp *ShipController) Shoot() {
	if sp.Missle != nil {
		sp.FireSource.Stop()
		sp.FireSource.Play()

		a := sp.Transform().Rotation()
		//scale := sp.Transform().Scale()
		for i := 0; i < len(sp.MisslesDirection[sp.MissleLevel]); i++ {
			pos := sp.MisslesPosition[i]
			s := sp.Transform().DirectionTransform(sp.MisslesDirection[sp.MissleLevel][i])
			//s = s.Mul(sp.MisslesDirection[i])

			p := sp.Transform().WorldPosition()
			_ = a
			m := engine.Identity()
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
			angle := float32(math.Atan2(float64(s.X), float64(s.Y))) * engine.DegreeConst

			nfire.Physics.Body.SetVelocity(float32(v.X), float32(v.Y))
			nfire.Physics.Body.AddForce(s.X*3000, s.Y*3000)

			nfire.Physics.Shape.Group = 1
			nfire.Physics.Body.SetMoment(engine.Inf)
			nfire.Transform().SetRotationf(180 - angle)
		}
	}
}

func (sp *ShipController) Update() {
	delta := float32(engine.DeltaTime())
	r2 := sp.Transform().DirectionTransform(engine.Up)
	r3 := sp.Transform().DirectionTransform(engine.Left)
	ph := sp.GameObject().Physics
	rx, ry := r2.X*delta, r2.Y*delta
	rsx, rsy := r3.X*delta, r3.Y*delta

	jet := false
	back := false

	if input.KeyDown('W') {
		ph.Body.AddForce(sp.Speed*rx, sp.Speed*ry)
		jet = true
	}

	if input.KeyDown('S') {
		ph.Body.AddForce(-sp.Speed*rx, -sp.Speed*ry)
		jet = true
		back = true
	}

	rotSpeed := sp.RotationSpeed
	if input.KeyDown(input.KeyLshift) {
		rotSpeed = 100
	}

	if sp.UseMouse {
		v := engine.GetScene().SceneBase().Camera.MouseWorldPosition()
		v = v.Sub(sp.Transform().WorldPosition())
		v.Normalize()
		angle := float32(math.Atan2(float64(v.Y), float64(v.X))) * engine.DegreeConst

		angle = engine.LerpAngle(sp.Transform().Rotation().Z, float32(int((angle - 90))), delta*rotSpeed/50)
		sp.Transform().SetRotationf(angle)

		ph.Body.SetAngularVelocity(0)
		ph.Body.SetTorque(0)

		if input.KeyDown('D') || input.KeyDown('E') {
			ph.Body.SetAngularVelocity(0)
			ph.Body.SetTorque(0)
			ph.Body.AddForce(-sp.Speed*rsx, -sp.Speed*rsy)
			jet = true
			back = true
		}
		if input.KeyDown('A') || input.KeyDown('Q') {
			ph.Body.SetAngularVelocity(0)
			ph.Body.SetTorque(0)
			ph.Body.AddForce(sp.Speed*rsx, sp.Speed*rsy)
			jet = true
			back = true
		}
	} else {
		r := sp.Transform().Rotation()
		if input.KeyDown('D') {
			ph.Body.SetAngularVelocity(0)
			ph.Body.SetTorque(0)
			sp.Transform().SetRotationf(r.Z - rotSpeed*delta)
			jet = true
			back = true
		}
		if input.KeyDown('A') {
			ph.Body.SetAngularVelocity(0)
			ph.Body.SetTorque(0)
			sp.Transform().SetRotationf(r.Z + rotSpeed*delta)
			jet = true
			back = true
		}

		if input.KeyDown('E') {
			ph.Body.SetAngularVelocity(0)
			ph.Body.SetTorque(0)
			ph.Body.AddForce(-sp.Speed*rsx, -sp.Speed*rsy)
			jet = true
			back = true
		}
		if input.KeyDown('Q') {
			ph.Body.SetAngularVelocity(0)
			ph.Body.SetTorque(0)
			ph.Body.AddForce(sp.Speed*rsx, sp.Speed*rsy)
			jet = true
			back = true
		}
	}

	if input.MouseDown(input.MouseLeft) {
		if time.Now().After(sp.lastShoot) {
			sp.Shoot()
			sp.lastShoot = time.Now().Add(time.Millisecond * 200)
		}
	}

	if input.KeyPress('P') {
		engine.EnablePhysics = !engine.EnablePhysics
	}
	if input.KeyPress('T') {
		sp.UseMouse = !sp.UseMouse
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
}
