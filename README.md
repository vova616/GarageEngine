## Install:
Windows:
To avoid installing mingw and downloading libraries in windows I have added the .a files.
<br/>
go get github.com/vova616/GarageEngine<br/>
go get github.com/vova616/chipmunk<br/>
go get github.com/vova616/gl <br/>
go get github.com/jteeuwen/glfw<br/>
(just to make sure you got all the sources, ignore all the erroes)<br/>

go to GarageEngine source folder and copy the pkg folder to your golang folder. (override)
now you can try to compile GarageEngine.

Other:
You need to download glfw/gl/glew libraries.<br/>
sudo apt-get install binutils-gold freeglut3 freeglut3-dev libglew1.5 libglew1.5-dev libglfw-dev<br/>
go get github.com/vova616/GarageEngine

## Dependencies
github.com/vova616/gl<br/>
github.com/vova616/chipmunk<br/>
github.com/jteeuwen/glfw

## Coroutines(they might be deprecated):
The useage is same as unity coroutines.<br/>
Use Behaviour Trees, its better and faster.

## Behaviour Trees:
Example in SpaceCookies/Game/EnemeyAI.go


## Videos:
http://www.youtube.com/watch?v=iMMbf6SRb9Q<br/>
http://www.youtube.com/watch?v=BMRlY9dFVLg

## Coroutines Example:
	func (sp *PlayerController) Start() {
		as := StartCoroutine(func() { sp.AutoShoot() })
		
		StartCoroutine(func() {
			Wait(3)
			YieldCoroutine(as) //wait for as to finish
			for i := 0; i < 10; i++ {
				YieldSkip()
				YieldSkip()
				YieldSkip()
				sp.Shoot()
			}
		})
	}

	func (sp *PlayerController) AutoShoot() {
		for i := 0; i < 3; i++ {
			Wait(3)
			sp.Shoot()
		}
	}

	func (sp *PlayerController) AutoShoot2() {
		for i := 0; i < 3; i++ {
			for i:=0;i<3*60;i++ {
				YieldSkip() //Frame skip
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

