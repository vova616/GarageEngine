## GarageEngine
This is an educational project, I'm learning as I go, I cannot promise backwards compatibility at this point.<br/>
the name will be probably changed.<br/> 
	
## Install:
You need to download glfw&glew libraries.

Windows:
Use mingw and compile glew.

Linux:
<br/>
sudo apt-get update 
<br/>
sudo apt-get install build-essential binutils-gold freeglut3 freeglut3-dev libglew-dev libglfw-dev libxrandr2 libxrandr-dev libglew libglew1.8  
<br/>
go get github.com/vova616/GarageEngine	

## To-Do list
Clean project:<br/>
Name changing Engine -> engine etc...<br/>
Function changing -> SetWorldPositionf -> SetWorldPosition2d etc...<br>
<br/>
Atlas,Font - Clean whatever we can.<br/>
Material - Think of design that does not require lots of work when creating custom shaders.<br/>
Physics - Think of a better design for arbiter and clean & polish stuff.<br/>
Scene - Get scene by name.<br/>
Tree Behaviours - Clean & polish & new features.<br/>
Camera - support multiple cameras.<br/>
Rendering - support auto-batching.<br/>
Coroutine - try to fix the bug that you cannot access to textures in Coroutines.<br/>
Readme - explain Tree Behaviours.<br/>
Learn from - https://github.com/runningwild/haunts .<br/>
Comments - lacks tons of it.<br/>



## Dependencies
github.com/go-gl/gl<br/>
github.com/go-gl/glfw<br/>
github.com/vova616/chipmunk<br/>
"github.com/vova616/freetype-go/freetype"

## Coroutines(they might be deprecated):
The useage is same as unity coroutines.<br/>
Use Behaviour Trees, its better and faster.

## Behaviour Trees:
Example in SpaceCookies/game/EnemeyAI.go

## SpaceCookies
Mini game to test the engine, it will host server on port 123 then you connect to it.
Make sure your executable file is in the same folder with the data folder.

## Videos:
http://www.youtube.com/watch?v=iMMbf6SRb9Q<br/>
http://www.youtube.com/watch?v=BMRlY9dFVLg
	
## Coroutines Example:
	func (sp *PlayerController) Start() {
		as := StartCoroutine(func() { sp.AutoShoot() })
		
		StartCoroutine(func() {
			CoSleep(3)
			YieldCoroutine(as) //wait for as to finish
			for i := 0; i < 10; i++ {
				CoCoYieldSkip()
				CoYieldSkip()
				CoYieldSkip()
				sp.Shoot()
			}
		})
	}

	func (sp *PlayerController) AutoShoot() {
		for i := 0; i < 3; i++ {
			CoSleep(3)
			sp.Shoot()
		}
	}

	func (sp *PlayerController) AutoShoot2() {
		for i := 0; i < 3; i++ {
			for i:=0;i<3*60;i++ {
				CoYieldSkip() //Frame skip
			}
			sp.Shoot()
		}
	}

	func (sp *PlayerController) AutoShoot3() {
		for i := 0; i < 3; i++ {
			Signal := NewSignal()
			go func() {
				<-time.After(time.Second * 3)
				Signal.SendEnd()
			}() 
			Yield(Signal)
			sp.Shoot()
		}
	} 

