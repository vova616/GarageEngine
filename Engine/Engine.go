package Engine

import (
	"github.com/vova616/gl"
	//"log"
	"github.com/vova616/GarageEngine/Engine/Input"
	//"os"
	"fmt"
	"github.com/go-gl/glfw"
	"github.com/vova616/chipmunk"
	"github.com/vova616/chipmunk/vect"
	"math"
	"runtime"
	"time"
)

type Arbiter struct {
	*chipmunk.Arbiter
	gameObjectA *GameObject //Always the owner
	gameObjectB *GameObject //Always the other target
	Swapped     bool
}

func newArbiter(arb *chipmunk.Arbiter, owner *GameObject) Arbiter {
	newArb := Arbiter{Arbiter: arb}
	pa, ba := arb.BodyA.CallbackHandler.(*Physics)
	if ba {
		newArb.gameObjectA = pa.GameObject()
	}
	pb, bb := arb.BodyB.CallbackHandler.(*Physics)
	if bb {
		newArb.gameObjectB = pb.GameObject()
	}

	//Find owner
	if newArb.gameObjectA != owner {
		newArb.gameObjectA, newArb.gameObjectB = newArb.gameObjectB, newArb.gameObjectA
		newArb.Swapped = true
	}
	return newArb
}

func (arbiter *Arbiter) GameObjectA() *GameObject {
	return arbiter.gameObjectA
}

func (arbiter *Arbiter) GameObjectB() *GameObject {
	return arbiter.gameObjectB
}

func (arbiter *Arbiter) Normal(contact *chipmunk.Contact) Vector {
	normal := contact.Normal()
	if arbiter.Swapped {
		return NewVector2(float32(-normal.X), float32(-normal.Y))
	}
	return NewVector2(float32(normal.X), float32(normal.Y))
}

func init() {
	fmt.Print()
}

const (
	RadianConst = math.Pi / 180
	DegreeConst = 180 / math.Pi
	MouseTag    = "Mouse"
)

var (
	Inf = vect.Float(math.Inf(1))

	scenes       []Scene = make([]Scene, 0)
	activeScenes []Scene = make([]Scene, 0)
	mainScene    Scene

	nextScene Scene = nil

	running        = false
	insideGameloop = false

	Space     *chipmunk.Space = nil
	deltaTime float64
	fixedTime float64
	gameTime  time.Time

	steps    = float64(1)
	stepTime = float64(1) / float64(60) / steps

	lastTime time.Time = time.Now()

	EnablePhysics = true
	Debug         = false
	InternalFPS   = float64(100)

	BehaviorTicks = 5

	Title  = "Engine Test"
	Width  = 1280
	Height = 720

	terminated chan bool
)

func init() {
	terminated = make(chan bool)

}

func LoadScene(scene Scene) {
	if insideGameloop {
		nextScene = scene
		return
	}

	ResourceManager.Release()
	coroutines = coroutines[:0]
	Routines = Routines[:0]

	if Space != nil {
		for _, g := range mainScene.SceneBase().gameObjects {
			if g != nil {
				g.Destroy()
			}
		}
		Iter(mainScene.SceneBase().gameObjects, destoyGameObject)
		mainScene.SceneBase().gameObjects = nil
		Space.Destory()
		runtime.GC()
		Space = chipmunk.NewSpace()
	} else {
		Space = chipmunk.NewSpace()
	}

	Input.ClearInput()

	sn := scene.New()
	sn.Load()

	internalFPS := NewGameObject("InternalFPS")
	internalFPS.AddComponent(NewFPS())
	sn.SceneBase().AddGameObject(internalFPS)

	mainScene = sn
}

func GetScene() Scene {
	return mainScene
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

func DeltaTime() float64 {
	return deltaTime
}

func GameTime() time.Time {
	return gameTime
}

func StartEngine() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	runtime.LockOSThread()
	println("Enginge started!")
	var err error
	if err = glfw.Init(); err != nil {
		panic(err)
	}
	println("GLFW Initialized!")

	glfw.OpenWindowHint(glfw.Accelerated, 1)

	if err = glfw.OpenWindow(Width, Height, 8, 8, 8, 8, 8, 8, glfw.Windowed); err != nil {
		panic(err)
	}

	glfw.SetSwapInterval(1) //0 to make FPS Maximum
	glfw.SetWindowTitle(Title)
	glfw.SetWindowSizeCallback(onResize)
	glfw.SetKeyCallback(Input.OnKey)
	glfw.SetCharCallback(Input.OnChar)
	glfw.SetMouseButtonCallback(Input.ButtonPress)
	Input.MousePosition = glfw.MousePos

	if err = initGL(); err != nil {
		panic(err)
	}
	println("Opengl Initialized!")

	TextureMaterial = NewBasicMaterial(spriteVertexShader, spriteFragmentShader)
	err = TextureMaterial.Load()
	if err != nil {
		fmt.Println(err)
	}

	SDFMaterial = NewBasicMaterial(sdfVertexShader, sdfFragmentShader)
	err = SDFMaterial.Load()
	if err != nil {
		fmt.Println(err)
	}

	gameTime = time.Time{}
	lastTime = time.Now()
}

func MainLoop() bool {
	running = true

	if nextScene != nil {
		s := nextScene
		nextScene = nil
		LoadScene(s)
	}

	insideGameloop = true
	if running && glfw.WindowParam(glfw.Opened) == 1 {
		Run()
	} else {
		return false
	}
	insideGameloop = false

	return true
}

func Run() {
	gl.ClearColor(0, 0, 0, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.LoadIdentity()

	gameTime = gameTime.Add(time.Since(lastTime))
	deltaTime = float64(time.Since(lastTime).Nanoseconds()) / float64(time.Second)
	lastTime = time.Now()
	before := time.Now()

	timer := NewTimer()
	timer.Start()

	var destroyDelta,
		startDelta,
		fixedUpdateDelta,
		physicsDelta,
		updateDelta,
		lateUpdateDelta,
		drawDelta,
		coroutinesDelta,
		stepDelta,
		behaviorDelta time.Duration

	if mainScene != nil {
		fixedTime += deltaTime
		sd := mainScene.SceneBase()

		arr := sd.gameObjects

		timer.StartCustom("Destory routines")
		Iter(arr, destoyGameObject)
		destroyDelta = timer.StopCustom("Destory routines")

		timer.StartCustom("Start routines")
		Iter(arr, startGameObject)
		startDelta = timer.StopCustom("Start routines")

		//

		timer.StartCustom("Physics time")
		if EnablePhysics {
			timer.StartCustom("Physics step time")
			for fixedTime >= stepTime {
				timer.StartCustom("FixedUpdate routines")
				Iter(arr, fixedUdpateGameObject)
				fixedUpdateDelta = timer.StopCustom("FixedUpdate routines")

				timer.StartCustom("Physics time")

				setPosition := func(g *GameObject) {
					if g.Physics != nil && g.Physics.started() {
						pos := g.Transform().WorldPosition()

						//Interpolation check: if position/angle has been changed directly and not by the physics engine, change g.Physics.lastPosition/lastAngle
						var pAngle vect.Float
						var pPos vect.Vect
						if vect.Float(pos.X) != g.Physics.lastPosition.X || vect.Float(pos.Y) != g.Physics.lastPosition.Y {
							g.Physics.lastPosition.X, g.Physics.lastPosition.Y = vect.Float(pos.X), vect.Float(pos.Y)
						}
						if vect.Float(g.Transform().WorldRotation().Z) != g.Physics.lastAngle {
							g.Physics.lastAngle = pAngle
						}
						pPos = g.Physics.lastPosition
						pAngle = g.Physics.lastAngle * RadianConst
						g.Physics.lastPosition = pPos
						g.Physics.lastAngle = g.Physics.lastAngle

						//Set physics data
						g.Physics.Body.SetAngle(pAngle)
						g.Physics.Body.SetPosition(pPos)

					}
				}
				Iter(arr, setPosition)

				Space.Step(vect.Float(stepTime))
				fixedTime -= stepTime

				updatePosition := func(g *GameObject) {
					if g.Physics != nil && g.Physics.started() {

						/*
							When parent changes his position/rotation it changes his children position/rotation too but the physics engine thinks its in different position
							so we need to check how much it changed and apply to the new position/rotation so we wont fuck up things too much.
						*/
						b := g.Physics.Body
						angle := float32(b.Angle()) * DegreeConst
						lAngle := float32(g.Physics.lastAngle)
						a := g.Transform().Angle()
						a += angle - lAngle
						g.Transform().SetWorldRotationf(a)

						pos := b.Position()
						lPos := g.Physics.lastPosition
						objPos := g.Transform().WorldPosition()
						objPos.X += float32(pos.X - lPos.X)
						objPos.Y += float32(pos.Y - lPos.Y)
						g.Transform().SetWorldPosition(objPos)

						//Interpolation 
						g.Physics.lastPosition = vect.Vect{vect.Float(objPos.X), vect.Float(objPos.Y)}
						g.Physics.lastAngle = vect.Float(a)
						if fixedTime > 0 && fixedTime < stepTime {
							alpha := vect.Float(fixedTime / stepTime)

							objPos.X = float32((vect.Float(objPos.X) * alpha) + (g.Physics.lastPosition.X * (1 - alpha)))
							objPos.Y = float32((vect.Float(objPos.Y) * alpha) + (g.Physics.lastPosition.Y * (1 - alpha)))
							a = float32((vect.Float(a) * alpha) + (g.Physics.lastAngle * (1 - alpha)))

							g.Transform().SetWorldRotationf(a)
							g.Transform().SetWorldPosition(objPos)
						}
					}
				}

				Iter(arr, updatePosition)

				physicsStepDelta := timer.StopCustom("Physics step time")

				//break if its taking too much time
				if float64(physicsStepDelta.Nanoseconds())/float64(time.Second) > stepTime*0.5 {
					break
				}
			}
		}
		physicsDelta = timer.StopCustom("Physics time")

		timer.StartCustom("Update routines")
		Iter(arr, udpateGameObject)
		updateDelta = timer.StopCustom("Update routines")

		timer.StartCustom("LateUpdate routines")
		Iter(arr, lateudpateGameObject)
		lateUpdateDelta = timer.StopCustom("LateUpdate routines")

		timer.StartCustom("Draw routines")
		Iter(arr, drawGameObject)
		drawDelta = timer.StopCustom("Draw routines")

		timer.StartCustom("coroutines")
		RunCoroutines()
		coroutinesDelta = timer.StopCustom("coroutines")

		timer.StartCustom("BehaviorTree")
		RunBT(BehaviorTicks)
		behaviorDelta = timer.StopCustom("BehaviorTree")

		Input.UpdateInput()

		stepDelta = timer.Stop()
	}

	timer.StartCustom("SwapBuffers")
	glfw.SwapBuffers()
	swapBuffersDelta := timer.StopCustom("SwapBuffers")

	now := time.Now()
	deltaDur := now.Sub(before)

	if Debug {
		fmt.Println()
		fmt.Println("##################")
		if InternalFPS < 20 {
			fmt.Println("FPS is lower than 20. FPS:", InternalFPS)
		} else if InternalFPS < 30 {
			fmt.Println("FPS is lower than 30. FPS:", InternalFPS)
		} else if InternalFPS < 40 {
			fmt.Println("FPS is lower than 40. FPS:", InternalFPS)
		}
		if stepDelta > 17*time.Millisecond {
			fmt.Println("StepDelta time is lower than normal")
		}
		fmt.Println("Debugging Times:")
		if (deltaDur.Nanoseconds() / int64(time.Millisecond)) != 0 {
			fmt.Println("Expected FPS", 1000/(deltaDur.Nanoseconds()/int64(time.Millisecond)))
		}
		fmt.Println("Step time", stepDelta)
		fmt.Println("Destroy time", destroyDelta)
		fmt.Println("Start time", startDelta)
		fmt.Println("FixedUpdate time", fixedUpdateDelta)
		fmt.Println("Update time", updateDelta)
		fmt.Println("LateUpdate time", lateUpdateDelta)
		fmt.Println("Draw time", drawDelta)
		fmt.Println("Delta time", deltaDur, deltaTime)
		fmt.Println("SwapBuffers time", swapBuffersDelta)
		fmt.Println("Coroutines time", coroutinesDelta)
		fmt.Println("BehaviorTree time", behaviorDelta)
		fmt.Println("------------------")
		fmt.Println("Physics time:", physicsDelta)
		fmt.Println("StepTime time", Space.StepTime)
		fmt.Println("ApplyImpulse time", Space.ApplyImpulsesTime)
		fmt.Println("ReindexQueryTime time", Space.ReindexQueryTime)
		fmt.Println("##################")
		fmt.Println()
	}
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
	if !gameObject.active {
		return
	}
	//mat := gameObject.Transform().Matrix()

	//gl.LoadMatrixf(mat.Ptr())

	l := len(gameObject.components)
	comps := gameObject.components

	for i := l - 1; i >= 0; i-- {
		if comps[i].started() {
			comps[i].Draw()
		}
	}
}

func IterExcept(objs []*GameObject, f func(*GameObject), except *GameObject) {
	l := len(objs)
	for i := l - 1; i >= 0; i-- {
		if objs[i] != except {
			f(objs[i])
		}
		arr2 := objs[i].Transform().Children()
		Iter2Except(arr2, f, except)
	}
}

func Iter2Except(objs []*Transform, f func(*GameObject), except *GameObject) {
	l := len(objs)
	for i := l - 1; i >= 0; i-- {
		if objs[i].GameObject() != except {
			f(objs[i].GameObject())
		}
		arr2 := objs[i].Children()
		Iter2Except(arr2, f, except)
	}
}

func startGameObject(gameObject *GameObject) {
	if !gameObject.active {
		return
	}
	l := len(gameObject.components)
	comps := gameObject.components

	for i := l - 1; i >= 0; i-- {
		if !comps[i].started() {
			comps[i].setStarted(true)
			comps[i].Start()
		}
	}
}

func destoyGameObject(gameObject *GameObject) {
	if gameObject.destoryMark {
		gameObject.destroy()
		mainScene.SceneBase().RemoveGameObject(gameObject)
	}
}

func onCollisionPreSolveGameObject(gameObject *GameObject, arb Arbiter) bool {
	if !gameObject.active {
		return true
	}
	l := len(gameObject.components)
	comps := gameObject.components

	b := true
	for i := l - 1; i >= 0; i-- {
		b = b && comps[i].OnCollisionPreSolve(arb)
	}
	return b
}

func onCollisionPostSolveGameObject(gameObject *GameObject, arb Arbiter) {
	if !gameObject.active {
		return
	}
	l := len(gameObject.components)
	comps := gameObject.components

	for i := l - 1; i >= 0; i-- {
		if comps[i].started() {
			comps[i].OnCollisionPostSolve(arb)
		}
	}
}

func onCollisionEnterGameObject(gameObject *GameObject, arb Arbiter) bool {
	if gameObject == nil || !gameObject.active {
		return true
	}
	l := len(gameObject.components)
	comps := gameObject.components

	b := true
	for i := l - 1; i >= 0; i-- {
		if comps[i].started() {
			b = b && comps[i].OnCollisionEnter(arb)
		}
	}
	return b
}

func onCollisionExitGameObject(gameObject *GameObject, arb Arbiter) {
	if gameObject == nil || !gameObject.active {
		return
	}
	l := len(gameObject.components)
	comps := gameObject.components

	for i := l - 1; i >= 0; i-- {
		if comps[i].started() {
			comps[i].OnCollisionExit(arb)
		}
	}
}

func onMouseEnterGameObject(gameObject *GameObject, arb Arbiter) bool {
	if gameObject == nil || !gameObject.active {
		return true
	}
	l := len(gameObject.components)
	comps := gameObject.components

	b := true
	for i := l - 1; i >= 0; i-- {
		if comps[i].started() {
			b = b && comps[i].OnMouseEnter(arb)
		}
	}
	return b
}

func onMouseExitGameObject(gameObject *GameObject, arb Arbiter) {
	if gameObject == nil || !gameObject.active {
		return
	}
	l := len(gameObject.components)
	comps := gameObject.components

	for i := l - 1; i >= 0; i-- {
		if comps[i].started() {
			comps[i].OnMouseExit(arb)
		}
	}
}

func udpateGameObject(gameObject *GameObject) {
	if !gameObject.active {
		return
	}

	l := len(gameObject.components)
	comps := gameObject.components

	for i := l - 1; i >= 0; i-- {
		if comps[i].started() {
			comps[i].Update()
		}
	}
}

func lateudpateGameObject(gameObject *GameObject) {
	if !gameObject.active {
		return
	}

	l := len(gameObject.components)
	comps := gameObject.components

	for i := l - 1; i >= 0; i-- {
		if comps[i].started() {
			comps[i].LateUpdate()
		}
	}
}

func fixedUdpateGameObject(gameObject *GameObject) {
	if !gameObject.active {
		return
	}

	l := len(gameObject.components)
	comps := gameObject.components

	for i := l - 1; i >= 0; i-- {
		if comps[i].started() {
			comps[i].FixedUpdate()
		}
	}
}

func initGL() (err error) {
	gl.Init()
	gl.ShadeModel(gl.SMOOTH)
	gl.ClearColor(0, 0, 0, 0)
	gl.ClearDepth(1)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Hint(gl.PERSPECTIVE_CORRECTION_HINT, gl.NICEST)
	gl.DepthFunc(gl.NEVER)
	gl.Enable(gl.BLEND)
	gl.DepthMask(true)

	//loadShader()

	return
}

func drawQuad(srcwidth, destwidth, srcheight, destheight float32) {
	gl.Begin(gl.QUADS)
	gl.TexCoord2i(0, 0)
	gl.Vertex2f(-1, -1)
	gl.TexCoord2i(int(srcwidth), 0)
	gl.Vertex2f(-1+destwidth, -1)
	gl.TexCoord2i(int(srcwidth), int(srcheight))
	gl.Vertex2f(-1+destwidth, -1+destheight)
	gl.TexCoord2i(0, int(srcheight))
	gl.Vertex2f(-1, -1+destheight)
	gl.End()
}

func onResize(w, h int) {
	if h <= 0 {
		h = 1
	}
	if w <= 0 {
		w = 1
	}

	Height = h
	Width = w

	gl.Viewport(0, 0, w, h)

	if GetScene() != nil && GetScene().SceneBase().Camera != nil {
		GetScene().SceneBase().Camera.UpdateResolution()
	}
}
