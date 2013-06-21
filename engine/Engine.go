package engine

import (
	"github.com/go-gl/gl"
	//"log"
	"github.com/vova616/GarageEngine/engine/bt"
	"github.com/vova616/GarageEngine/engine/cr"
	"github.com/vova616/GarageEngine/engine/input"
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
	internalFPSObject   *FPS

	BehaviorTicks = 5

	windowTitle = "Engine Test"
	Width       = 1280
	Height      = 720

	terminated chan bool
)

func LoadScene(scene Scene) {
	if insideGameloop {
		nextScene = scene
		return
	}

	ResourceManager.Release()
	bt.Clear()
	cr.Clear()

	if Space != nil {
		for _, g := range mainScene.SceneBase().gameObjects {
			if g != nil {
				g.Destroy()
			}
		}
		iter(&mainScene.SceneBase().gameObjects, destoyGameObject)
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

	internalFPSGObject := NewGameObject("InternalFPS")
	internalFPSObject = NewFPS()
	internalFPSGObject.AddComponent(internalFPSObject)
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

func Terminate() {
	glfw.Terminate()
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
	glfw.SetMouseWheel(0)
	glfw.SetMouseWheelCallback(input.MouseWheelCallback)

	input.MouseWheelPosition = glfw.MouseWheel
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

	CurrentCamera().Clear()

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

		arr := &sd.gameObjects

		timer.StartCustom("Destory routines")
		iter(arr, destoyGameObject)
		destroyDelta = timer.StopCustom("Destory routines")

		//Better to not do it every frame.
		mainScene.SceneBase().cleanNil()

		timer.StartCustom("Start routines")
		iter(arr, startGameObject)
		startDelta = timer.StopCustom("Start routines")

		//

		timer.StartCustom("Physics time")
		if EnablePhysics {
			timer.StartCustom("Physics step time")

			for fixedTime >= stepTime {
				timer.StartCustom("FixedUpdate routines")
				iter(arr, fixedUdpateGameObject)
				fixedUpdateDelta = timer.StopCustom("FixedUpdate routines")

				timer.StartCustom("PreStep Physics Delta")
				iter(arr, preStepGameObject)
				startPhysicsDelta = timer.StopCustom("PreStep Physics Delta")

				Space.Step(vect.Float(stepTime))
				fixedTime -= stepTime

				_ = timer.StopCustom("Physics step time")

				timer.StartCustom("PostStep Physics Delta")
				iter(arr, postStepGameObject)
				endPhysicsDelta = timer.StopCustom("PostStep Physics Delta")
			}
			if fixedTime > 0 && fixedTime < stepTime {
				iter(arr, interpolateGameObject)
			}
		}
		physicsDelta = timer.StopCustom("Physics time")

		timer.StartCustom("Update routines")
		internalFPSObject.Update()
		iter(arr, udpateGameObject)
		updateDelta = timer.StopCustom("Update routines")

		timer.StartCustom("LateUpdate routines")
		iter(arr, lateudpateGameObject)
		lateUpdateDelta = timer.StopCustom("LateUpdate routines")

		timer.StartCustom("Draw routines")
		depthMap.Iter(drawGameObject)
		drawDelta = timer.StopCustom("Draw routines")

		timer.StartCustom("coroutines")
		cr.Run()
		coroutinesDelta = timer.StopCustom("coroutines")

		timer.StartCustom("BehaviorTree")
		bt.Run(BehaviorTicks)
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

	if GetScene() != nil && CurrentCamera() != nil {
		CurrentCamera().UpdateResolution()
	}
}

func PanicPath() string {
	fullPath := ""
	skip := 3
	for i := skip; ; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		if i > skip {
			fullPath += ", "
		}
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		fullPath += fmt.Sprintf("%s:%d", file, line)
	}
	return fullPath
}
