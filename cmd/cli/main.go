package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/igoramorim/gopengl/internal/scenes"
	"github.com/igoramorim/gopengl/internal/sshot"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main() {
	if len(os.Args) <= 1 {
		help()
		os.Exit(1)
	}

	arg := os.Args[1]
	scene, ok := allScenes[arg]
	if !ok {
		help()
		os.Exit(1)
	}

	defer func() {
		err := sshot.MakeGIF()
		if err != nil {
			fmt.Println(err.Error())
		}
	}()

	scene()
}

var allScenes = map[string]func(){
	scenes.Triangle{}.Name():         scenes.Triangle{}.Show,
	scenes.Shaders{}.Name():          scenes.Shaders{}.Show,
	scenes.Textures{}.Name():         scenes.Textures{}.Show,
	scenes.Transformations{}.Name():  scenes.Transformations{}.Show,
	scenes.CoordinateSystem{}.Name(): scenes.CoordinateSystem{}.Show,
	scenes.Cube{}.Name():             scenes.Cube{}.Show,
	scenes.Camera{}.Name():           scenes.NewCamera().Show,
}

func help() {
	fmt.Printf("scene name is required\n")
	fmt.Printf("possible values are: %q\n", possibleScenes())
}

func possibleScenes() []string {
	var names []string
	for k, _ := range allScenes {
		names = append(names, k)
	}
	return names
}
