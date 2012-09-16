## Install:
Windows:
To avoid installing mingw and downloading libraries in windows I have added the .a files.

go get github.com/vova616/GarageEngine
go get github.com/vova616/chipmunk
go get github.com/vova616/gl 
go get github.com/jteeuwen/glfw
(just to make sure you got all the sources, ignore all the erroes)

go to GarageEngine source folder and copy the pkg folder to your golang folder. (override)
now you can try to compile GarageEngine.

Other:
go get github.com/vova616/GarageEngine

you need to download glew and glfw libs.
for glfw look here github.com/jteeuwen/glfw.
for gl just google glew and download them.


## Goroutines:
Its a managed goroutines the useage is same as unity coroutines.

## Example:
	func (sp *PlayerController) Start() {
		StartGoroutine(func() { sp.AutoShoot() })
	}

	func (sp *PlayerController) AutoShoot() {
		for {
			Wait(3)
			sp.Shoot()
		}
	}

	func (sp *PlayerController) AutoShoot2() {
		for {
			for i:=0;i<3*60;i++ {
				YieldSkip() //Frame skip
			}
			sp.Shoot()
		}
	}

	func (sp *PlayerController) AutoShoot3() {
		for {
			Signal := NewSignal()
			go func() {
				<-time.After(time.Second * 3)
				Signal.SendEnd()
			}() 
			Yield(Signal)
			sp.Shoot()
		}
	} 

