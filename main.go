package main

import (
	"flag"
	"fmt"
	"github.com/vova616/GarageEngine/engine"
	"github.com/vova616/GarageEngine/engine/input"
	"github.com/vova616/GarageEngine/networkOnline"
	"github.com/vova616/GarageEngine/spaceCookies/game"
	"github.com/vova616/GarageEngine/spaceCookies/login"
	"github.com/vova616/GarageEngine/spaceCookies/server"
	"github.com/vova616/GarageEngine/zumbies"
	//"math"
	//"github.com/go-gl/gl"
	"os"
	"runtime"
	"runtime/pprof"

	//"time"
)

var cpuprofile = flag.String("p", "", "write cpu profile to file")
var memprofile = flag.String("m", "", "write mem profile to file")

func main() {
	runtime.GOMAXPROCS(8)
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			fmt.Errorf("%s\n", err)
		}
		pprof.StartCPUProfile(f)

		defer pprof.StopCPUProfile()
	}

	file, _ := os.Create("./log.txt")
	os.Stdout = file
	os.Stderr = file
	os.Stdin = file
	defer file.Close()

	Start()
	engine.Terminated()

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			fmt.Errorf("%s\n", err)
		}
		pprof.WriteHeapProfile(f)
		f.Close()
		return
	}
}

func Start() {
	defer func() {
		if p := recover(); p != nil {
			fmt.Println(p, engine.PanicPath())
		}

		go engine.Terminate()
	}()
	engine.StartEngine()
	_ = game.GameSceneGeneral
	_ = networkOnline.GameSceneGeneral
	_ = login.LoginSceneGeneral
	_ = zumbies.GameSceneGeneral

	/*
		Running local server.
	*/
	go server.StartServer()

	scene := 0

	engine.LoadScene(login.LoginSceneGeneral)
	for engine.MainLoop() {
		if input.KeyPress('`') {
			scene = (scene + 1) % 3
			switch scene {
			case 0:
				engine.LoadScene(login.LoginSceneGeneral)
			case 1:
				engine.LoadScene(networkOnline.GameSceneGeneral)
			case 2:
				engine.LoadScene(zumbies.GameSceneGeneral)
			}

		}
	}
}

/*
	Need to freeze physics which are disabled which is not effective without changing the physics engine? (done need checking)
	Need to destroy children of gameobject (already did this, just check it again)
	When removing object from scene we need to let the physics engine know about it (done need checking)
	or disable completly removal of objects from scene because we already have active/inactive (nope)
	Make depth test scene.  also consider adding Z Buffer.
	Make Arbiter to return the correct shapes when swapped.
	Recheck all the Remove/Add-FromScene stuff and make sure setActiveRecursiveSilent are run only once.
*/
