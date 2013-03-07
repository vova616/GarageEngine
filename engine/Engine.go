package engine

import (
	"github.com/vova616/gl"
	//"log"
	"github.com/vova616/garageEngine/engine/input"
	//"os"
	"fmt"
	"github.com/go-gl/glfw"
	"github.com/vova616/chipmunk"
	"github.com/vova616/chipmunk/vect"
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

	steps          = float64(1)
	stepTime       = float64(1) / float64(60) / steps
	maxPhysicsTime = float64(1) / float64(40)

	lastTime time.Time = time.Now()

	CorrectWrongPhysics = true
	EnablePhysics       = true
	Debug               = false
	InternalFPS         = float64(100)

	BehaviorTicks = 5

	windowTitle = "Engine Test"
	Width       = 1280
	Height      = 720

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

	input.ClearInput()

	sn := scene.New()
	mainScene = sn
	sn.Load()

	internalFPS := NewGameObject("InternalFPS")
	internalFPS.AddComponent(NewFPS())
	sn.SceneBase().AddGameObject(internalFPS)
}

func GetScene() Scene {
	return mainScene
}

func CurrentCamera() *Camera {
	return mainScene.SceneBase().Camera
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

func SetTitle(title string) {
	glfw.SetWindowTitle(title)
	windowTitle = title
}

func Title() string {
	return windowTitle
}

func StartEngine() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	runtime.LockOSThread()
	fmt.Println("Enginge started!")
	var err error
	if err = glfw.Init(); err != nil {
		panic(err)
	}
	fmt.Println("GLFW Initialized!")

	glfw.OpenWindowHint(glfw.Accelerated, 1)

	if err = glfw.OpenWindow(Width, Height, 8, 8, 8, 8, 8, 8, glfw.Windowed); err != nil {
		panic(err)
	}

	glfw.SetSwapInterval(1) //0 to disable vsync, 1 to enable it
	glfw.SetWindowTitle(windowTitle)
	glfw.SetWindowSizeCallback(onResize)
	glfw.SetKeyCallback(input.OnKey)
	glfw.SetCharCallback(input.OnChar)
	glfw.SetMouseButtonCallback(input.ButtonPress)

	input.MousePosition = glfw.MousePos

	if err = initGL(); err != nil {
		panic(err)
	}
	fmt.Println("Opengl Initialized!")

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

	internalMaterial = NewBasicMaterial(spriteVertexShader, spriteFragmentShader)
	err = internalMaterial.Load()
	if err != nil {
		fmt.Println(err)
	}

	initDefaultPlane()
	glfw.SwapBuffers()

	gameTime = time.Time{}
	lastTime = time.Now()
	dl = glfw.Time()
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

var dd = float64(0)
var dl = float64(0)

func Run() {
	timeNow := time.Now()
	gameTime = gameTime.Add(timeNow.Sub(lastTime))
	//deltaTime = float64(timeNow.Sub(lastTime).Nanoseconds()) / float64(time.Second)
	lastTime = timeNow
	before := timeNow

	dn := glfw.Time()
	dd = dn - dl
	dl = dn
	deltaTime = dd

	timer := NewTimer()
	timer.Start()

	//time.Sleep(time.Second)

	gl.ClearColor(0, 0, 0, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.LoadIdentity()

	var destroyDelta,
		startDelta,
		fixedUpdateDelta,
		physicsDelta,
		updateDelta,
		lateUpdateDelta,
		drawDelta,
		coroutinesDelta,
		stepDelta,
		behaviorDelta,
		startPhysicsDelta,
		endPhysicsDelta time.Duration

	if mainScene != nil {
		d := deltaTime
		if d > maxPhysicsTime {
			d = maxPhysicsTime
		}
		fixedTime += d
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

				timer.StartCustom("PreStep Physics Delta")
				Iter(arr, preStepGameObject)
				startPhysicsDelta = timer.StopCustom("PreStep Physics Delta")

				Space.Step(vect.Float(stepTime))
				fixedTime -= stepTime

				_ = timer.StopCustom("Physics step time")

				timer.StartCustom("PostStep Physics Delta")
				Iter(arr, postStepGameObject)
				endPhysicsDelta = timer.StopCustom("PostStep Physics Delta")
			}
			if fixedTime > 0 && fixedTime < stepTime {
				Iter(arr, interpolateGameObject)
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
		for i := int8(-127); ; i++ {
			drawArr, exists := depthMap[i]
			if exists && len(drawArr) > 0 {
				IterNoChildren(drawArr, drawGameObject)
			}
			if i == 127 {
				break
			}
		}
		drawDelta = timer.StopCustom("Draw routines")

		timer.StartCustom("coroutines")
		RunCoroutines()
		coroutinesDelta = timer.StopCustom("coroutines")

		timer.StartCustom("BehaviorTree")
		RunBT(BehaviorTicks)
		behaviorDelta = timer.StopCustom("BehaviorTree")

		input.UpdateInput()

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
		fmt.Println("PreStepDelta time", startPhysicsDelta)
		fmt.Println("PostStepDelta time", endPhysicsDelta)
		fmt.Println("StepTime time", Space.StepTime)
		fmt.Println("ApplyImpulse time", Space.ApplyImpulsesTime)
		fmt.Println("ReindexQueryTime time", Space.ReindexQueryTime)
		fmt.Println("Arbiters", len(Space.Arbiters))
		fmt.Println("##################")
		fmt.Println()
	}
}

func Iter(objs []*GameObject, f func(*GameObject)) {
	for i := len(objs) - 1; i >= 0; i-- {
		obj := objs[i]
		if obj != nil {
			f(obj)
			//Checks if the objs array has been changed
			if obj != objs[i] {
				i++
			} else {
				Iter2(obj.Transform().children, f)
			}
		}
	}
}

func Iter2(objs []*Transform, f func(*GameObject)) {
	for i := len(objs) - 1; i >= 0; i-- {
		if objs[i] != nil {
			obja := objs[i]
			obj := obja.GameObject()
			f(obj)
			//Checks if the objs array has been changed
			if obja != objs[i] {
				i++
			} else {
				Iter2(obj.Transform().children, f)
			}
		}
	}
}

func IterNoChildren(objs []*GameObject, f func(*GameObject)) {
	for i := len(objs) - 1; i >= 0; i-- {
		obj := objs[i]
		if obj != nil {
			f(obj)
			//Checks if the objs array has been changed
			if obj != objs[i] {
				i++
			}
		}
	}
}

func preStepGameObject(g *GameObject) {
	if g.Physics != nil && g.active && !g.Physics.Body.IsStatic() && g.Physics.started() {
		pos := g.Transform().WorldPosition()
		angle := g.Transform().Angle() * RadianConst

		if g.Physics.Interpolate {
			//Interpolation check: if position/angle has been changed directly and not by the physics engine, change g.Physics.lastPosition/lastAngle
			if vect.Float(pos.X) != g.Physics.interpolatedPosition.X || vect.Float(pos.Y) != g.Physics.interpolatedPosition.Y {
				g.Physics.interpolatedPosition = vect.Vect{vect.Float(pos.X), vect.Float(pos.Y)}
				g.Physics.Body.SetPosition(g.Physics.interpolatedPosition)
			}
			if vect.Float(angle) != g.Physics.interpolatedAngle {
				g.Physics.interpolatedAngle = vect.Float(angle)
				g.Physics.Body.SetAngle(g.Physics.interpolatedAngle)
			}
		} else {
			var pPos vect.Vect
			pPos.X, pPos.Y = vect.Float(pos.X), vect.Float(pos.Y)

			g.Physics.Body.SetAngle(vect.Float(angle))
			g.Physics.Body.SetPosition(pPos)
		}
		g.Physics.lastPosition = g.Physics.Body.Position()
		g.Physics.lastAngle = g.Physics.Body.Angle()
	}
}

func postStepGameObject(g *GameObject) {
	if g.Physics != nil && g.active && !g.Physics.Body.IsStatic() && g.Physics.started() {
		/*
			When parent changes his position/rotation it changes his children position/rotation too but the physics engine thinks its in different position
			so we need to check how much it changed and apply to the new position/rotation so we wont fuck up things too much.

			Note:If position/angle is changed in between preStep and postStep it will be overrided.
		*/
		if CorrectWrongPhysics {
			b := g.Physics.Body
			angle := float32(b.Angle())
			lAngle := float32(g.Physics.lastAngle)
			lAngle += angle - lAngle

			pos := b.Position()
			lPos := g.Physics.lastPosition
			lPos.X += (pos.X - lPos.X)
			lPos.Y += (pos.Y - lPos.Y)

			if g.Physics.Interpolate {
				g.Physics.interpolatedAngle = vect.Float(lAngle)
				g.Physics.interpolatedPosition = lPos
			}

			b.SetPosition(lPos)
			b.SetAngle(g.Physics.interpolatedAngle)

			g.Transform().SetWorldRotationf(lAngle * DegreeConst)
			g.Transform().SetWorldPositionf(float32(lPos.X), float32(lPos.Y))
		} else {
			b := g.Physics.Body
			angle := b.Angle()
			pos := b.Position()

			if g.Physics.Interpolate {
				g.Physics.interpolatedAngle = angle
				g.Physics.interpolatedPosition = pos
			}

			g.Transform().SetWorldRotationf(float32(angle) * DegreeConst)
			g.Transform().SetWorldPositionf(float32(pos.X), float32(pos.Y))
		}
	}
}

func interpolateGameObject(g *GameObject) {
	if g.Physics != nil && g.Physics.Interpolate && g.active && !g.Physics.Body.IsStatic() && g.Physics.started() {
		nextPos := g.Physics.Body.Position()
		currPos := g.Physics.lastPosition

		nextAngle := g.Physics.Body.Angle()
		currAngle := g.Physics.lastAngle

		alpha := vect.Float(fixedTime / stepTime)
		x := currPos.X + ((nextPos.X - currPos.X) * alpha)
		y := currPos.Y + ((nextPos.Y - currPos.Y) * alpha)
		a := currAngle + ((nextAngle - currAngle) * alpha)
		g.Transform().SetWorldPositionf(float32(x), float32(y))
		g.Transform().SetWorldRotationf(float32(a) * DegreeConst)

		g.Physics.interpolatedAngle = a
		g.Physics.interpolatedPosition.X, g.Physics.interpolatedPosition.Y = x, y
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
	for i := len(objs) - 1; i >= 0; i-- {
		obj := objs[i]
		if obj != except {
			f(obj)
		}
		//Checks if the objs array has been changed
		if obj != objs[i] {
			i++
		} else {
			Iter2Except(obj.Transform().children, f, except)
		}
	}
}

func Iter2Except(objs []*Transform, f func(*GameObject), except *GameObject) {
	for i := len(objs) - 1; i >= 0; i-- {
		obj := objs[i].GameObject()
		if obj == nil {
			continue
		}
		obja := objs[i]
		if obj != except {
			f(obj)
		}
		//Checks if the objs array has been changed
		if obja != objs[i] {
			i++
		} else {
			Iter2Except(obj.Transform().children, f, except)
		}
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
	if !gameObject.active || gameObject.Physics == nil {
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
