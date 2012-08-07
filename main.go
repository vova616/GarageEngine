package main

import (
	. "./Engine"
	"./NetworkOnline"
	"flag"
	"fmt"
	"math"
	"os"
	"runtime"
	"runtime/pprof"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	println("asd")
	os.Chdir("./bin/windows_386/")
	os.Chdir("./bin/windows_amd64/")

	Radian := math.Pi / 180
	Degree := 180 / math.Pi
	_ = Degree
	_ = Radian

	fmt.Println(math.Pi)

	fmt.Println(math.Cos(25))
	fmt.Println(math.Sin(25))
	fmt.Println(math.Atan(math.Sin(25*Radian)/math.Cos(25*Radian)) * Degree)

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
	LoadScene(NetworkOnline.GameSceneGeneral)
	for MainLoop() {

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
