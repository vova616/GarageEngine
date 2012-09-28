package main

import (
	"flag"
	"fmt"
	. "github.com/vova616/GarageEngine/Engine"
	"github.com/vova616/GarageEngine/NetworkOnline"
	"github.com/vova616/GarageEngine/SpaceCookies"
	//"math"
	"os"
	"runtime/pprof"
)

var cpuprofile = flag.String("p", "", "write cpu profile to file")

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
	Terminated()
}

func Start() {
	defer func() {
		if p := recover(); p != nil {
			fmt.Println(p, PanicPath())
		}
		Terminate()
	}()
	StartEngine()
	_ = SpaceCookies.GameSceneGeneral
	_ = NetworkOnline.GameSceneGeneral
	LoadScene(SpaceCookies.GameSceneGeneral)
	for MainLoop() {

	}
}
