package main

import (
	"flag"
	"fmt"
	"github.com/vova616/garageEngine/engine"
	"github.com/vova616/garageEngine/networkOnline"
	"github.com/vova616/garageEngine/spaceCookies/game"
	"github.com/vova616/garageEngine/spaceCookies/login"
	"github.com/vova616/garageEngine/spaceCookies/server"
	"github.com/vova616/garageEngine/zumbies"
	//"math"
	//"github.com/vova616/gl"
	"os"
	//"runtime"
	"runtime/pprof"
	//"time"
)

var cpuprofile = flag.String("p", "", "write cpu profile to file")
var memprofile = flag.String("m", "", "write mem profile to file")

func main() {
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

	go Start()
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

		engine.Terminate()
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

	engine.LoadScene(login.LoginSceneGeneral)
	for engine.MainLoop() {

	}
}

/*

	Need to freeze physics which are disabled which is not effective without changing the physics engine? (done need checking)
	Need to destroy children of gameobject (already did this, just check it again)
	When removing object from scene we need to let the physics engine know about it (done need checking)
	or disable completly removal of objects from scene because we already have active/inactive (nope)
	Make depth test scene.  also consider adding Z Buffer.
	Make Arbiter to return the correct shapes when swapped.
*/
