package Engine

import (
	"github.com/banthar/gl"
	//"log"
	. "../Engine/Input"
	//"os"
	"fmt"
	"github.com/jteeuwen/glfw"
	c "chipmunk"
	. "chipmunk/vect"
	"math"
	"runtime"
	"time"
)

func init() {
	fmt.Print()
}

const (
	RadianConst = math.Pi / 180
	DegreeConst = 180 / math.Pi
	
)



var (
	
	Inf = Float(math.Inf(1))
	
	scenes       []Scene = make([]Scene, 0)
	activeScenes []Scene = make([]Scene, 0)
	mainScene    Scene
	running                       = false
	Space        *c.Space = c.NewSpace()
	deltaTime    float32
	fixedTime    float32
	stepTime     = float32(1) / float32(50)

	Title  = "Engine Test"
	Width  = 640
	Height = 480

	terminated chan bool
)

func init() {
	Space.Iterations = 100
	terminated = make(chan bool)
}

func LoadScene(scene Scene) {
	sn := scene.New()
	sn.Load()
	mainScene = sn
}

func AddScene(scene Scene) {
	scenes = append(scenes, scene)
}

func ShutdownRecived() {
	terminated <- true
}

func Terminated() {
	<-terminated
}

func Terminate() {
	glfw.Terminate()
	ShutdownRecived()
}

func DeltaTime() float32 {
	return deltaTime
}

func StartEngine() {
	runtime.LockOSThread()
	println("Enginge started!")
	var err error
	if err = glfw.Init(); err != nil {
		panic(err)
	}

	if err = glfw.OpenWindow(Width, Height, 8, 8, 8, 8, 8, 8, glfw.Windowed); err != nil {
		panic(err)
	}

	glfw.SetSwapInterval(1) //0 to make FPS Maximum
	glfw.SetWindowTitle(Title)
	glfw.SetWindowSizeCallback(onResize)
	glfw.SetKeyCallback(OnKey)
	glfw.SetMouseButtonCallback(ButtonPress)

	if err = initGL(); err != nil {
		panic(err)
	}
}

func MainLoop() bool {
	running = true
	if running && glfw.WindowParam(glfw.Opened) == 1 {
		Run()
	} else {
		return false
	}
	return true
}

func Run() {
	before := time.Now()
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.LoadIdentity()

	if mainScene != nil {
		fixedTime += deltaTime
		sd := mainScene.SceneBase()

		arr := sd.gameObjects

		Iter(arr, startGameObject)

		for fixedTime > stepTime {
			Iter(arr, fixedUdpateGameObject)
			for _, b := range Space.AllBodies {
				g, ok := b.UserData.(*Physics)
				if ok && g != nil {
					pos := g.Transform().WorldPosition()
					b.SetAngle(Float(180-g.Transform().WorldRotation().Z)*RadianConst)
					//fmt.Println(g.Transform().WorldRotation().Z, b.Transform.Angle())
					b.SetPosition(Vect{Float(pos.X), Float(pos.Y)})
					*g.lastCollision = *g.currentCollision
					g.currentCollision.ShapeA = nil 
					g.currentCollision.ShapeB = nil
				}
			}
			Space.Step(Float(stepTime)) 
			//Space.Step(Float(0.1)) 

			fixedTime -= stepTime

			updatePosition := func(g *GameObject) {
				if g.Physics != nil {

					b := g.Physics.Body
					r := g.Transform().WorldRotation() 
					g.Transform().SetWorldRotation(NewVector3(r.X, r.Y, 180-(float32(b.Angle())*DegreeConst)))
					pos := b.Position() 
					g.Transform().SetWorldPosition(NewVector2(float32(pos.X), float32(pos.Y)))

					//fmt.Println(r.Z, b.Transform.Angle())

				}
			}

			Iter(arr, updatePosition)

			for _,i := range Space.Arbiters {
				if i.NumContacts == 0 {
					continue
				}
				a, _ := i.ShapeA.Body.UserData.(*Physics)
				b, _ := i.ShapeB.Body.UserData.(*Physics)
				if a != nil && b != nil {
					*a.currentCollision = *i
					*b.currentCollision = *i
					if (a.lastCollision.ShapeA == a.currentCollision.ShapeA && a.lastCollision.ShapeB == a.currentCollision.ShapeB) ||
						(a.lastCollision.ShapeA == a.currentCollision.ShapeB && a.lastCollision.ShapeB == a.currentCollision.ShapeA) {
						onCollisionGameObject(a.GameObject(), a.currentCollision)
					} else {
						if a.lastCollision.ShapeA != nil && a.lastCollision.ShapeB != nil {
							onCollisionExitGameObject(a.GameObject(), a.lastCollision)
						}
						onCollisionEnterGameObject(a.GameObject(), a.currentCollision)
						onCollisionGameObject(a.GameObject(), a.currentCollision)
					}
					if (b.lastCollision.ShapeA == b.currentCollision.ShapeA && b.lastCollision.ShapeB == b.currentCollision.ShapeB) ||
						(b.lastCollision.ShapeA == b.currentCollision.ShapeB && b.lastCollision.ShapeB == b.currentCollision.ShapeA) {
						onCollisionGameObject(b.GameObject(), b.currentCollision)
					} else {
						if b.lastCollision.ShapeA != nil && b.lastCollision.ShapeB != nil {
							onCollisionExitGameObject(b.GameObject(), b.lastCollision)
						}
						onCollisionEnterGameObject(b.GameObject(), b.currentCollision)
						onCollisionGameObject(b.GameObject(), b.currentCollision)
					}

				}
			}

			for _, b := range Space.AllBodies {
				g, ok := b.UserData.(*Physics)
				if ok && g != nil {
					if g.lastCollision.ShapeA != nil && g.lastCollision.ShapeB != nil &&
						g.currentCollision.ShapeA == nil && g.currentCollision.ShapeB == nil {
						onCollisionExitGameObject(g.GameObject(), g.lastCollision)
					}
				}
			}
		}
		Iter(arr, udpateGameObject)
		Iter(arr, drawGameObject)

		UpdateInput()
	}
	glfw.SwapBuffers()
	now := time.Now()
	deltaTime = float32(now.Sub(before).Nanoseconds()/int64(time.Millisecond)) / 1000
}

func Iter(objs []*GameObject, f func(*GameObject)) {
	l := len(objs)
	for i := l - 1; i >= 0; i-- {
		f(objs[i])
		arr2 := objs[i].Transform().Children()
		Iter2(arr2, f)
	}
}

func Iter2(objs []*Transform, f func(*GameObject)) {
	l := len(objs)
	for i := l - 1; i >= 0; i-- {
		f(objs[i].GameObject())
		arr2 := objs[i].Children()
		Iter2(arr2, f)
	}
}

func drawGameObject(gameObject *GameObject) {

	mat := gameObject.Transform().Matrix()

	gl.LoadMatrixf(mat.Ptr())

	l := len(gameObject.components)
	comps := gameObject.components

	for i := l - 1; i >= 0; i-- {
		if comps[i].Component().drawableComponent != nil {
			comps[i].Component().drawableComponent.Draw()
		}
	}
}

func startGameObject(gameObject *GameObject) {
	l := len(gameObject.components)
	comps := gameObject.components

	for i := l - 1; i >= 0; i-- {
		if comps[i].Component().startableComponent != nil && !comps[i].Component().started {
			comps[i].Component().started = true
			comps[i].Component().startableComponent.Start()
		}
	}
}

func onCollisionGameObject(gameObject *GameObject, arb *c.Arbiter) {
	l := len(gameObject.components)
	comps := gameObject.components

	for i := l - 1; i >= 0; i-- {
		if comps[i].Component().onCollisionComponent != nil {
			comps[i].Component().onCollisionComponent.OnCollision(NewCollision(arb))
		}
	}
}

func onCollisionEnterGameObject(gameObject *GameObject, arb *c.Arbiter) {
	l := len(gameObject.components)
	comps := gameObject.components

	for i := l - 1; i >= 0; i-- {
		if comps[i].Component().onCollisionEnterComponent != nil {
			comps[i].Component().onCollisionEnterComponent.OnCollisionEnter(NewCollision(arb))
		}
	}
}

func onCollisionExitGameObject(gameObject *GameObject, arb *c.Arbiter) {
	l := len(gameObject.components)
	comps := gameObject.components

	for i := l - 1; i >= 0; i-- {
		if comps[i].Component().onCollisionExitComponent != nil {
			comps[i].Component().onCollisionExitComponent.OnCollisionExit(NewCollision(arb))
		}
	}
}

func udpateGameObject(gameObject *GameObject) {
	l := len(gameObject.components)
	comps := gameObject.components

	for i := l - 1; i >= 0; i-- {
		if comps[i].Component().updateableComponent != nil {
			comps[i].Component().updateableComponent.Update()
		}
	}
}

func fixedUdpateGameObject(gameObject *GameObject) {
	l := len(gameObject.components)
	comps := gameObject.components

	for i := l - 1; i >= 0; i-- {
		if comps[i].Component().fUpdateableComponent != nil {
			comps[i].Component().fUpdateableComponent.FixedUpdate()
		}
	}
}

func initGL() (err error) {
	gl.Init()
	gl.ShadeModel(gl.SMOOTH)
	gl.ClearColor(0, 0, 0, 0)
	gl.ClearDepth(1)
	gl.DepthFunc(gl.LEQUAL)
	//gl.BlendFunc(gl.DST_ALPHA, gl.ZERO);
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Hint(gl.PERSPECTIVE_CORRECTION_HINT, gl.NICEST)
	//gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.NEVER)
	gl.Enable(gl.TEXTURE_2D)
	gl.Enable(gl.BLEND)
	gl.DepthMask(true)
	return
}

func onResize(w, h int) {
	if h == 0 {
		h = 1
	}
	if w == 0 {
		w = 1
	}

	Height = h
	Width = w

	gl.Viewport(0, 0, w, h)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	//glu.Perspective(45.0, float64(w)/float64(h), 0.1, 100.0)
	//glu.Ortho2D(0,float64(w),0,float64(h))
	gl.Ortho(0, float64(w), 0, float64(h), -1000000, 1000000)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
}
