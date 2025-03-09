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

	scene.Show()
}

var allScenes = map[string]scenes.Scene{
	scenes.Triangle{}.Name():         scenes.Triangle{},
	scenes.Shaders{}.Name():          scenes.Shaders{},
	scenes.Textures{}.Name():         scenes.Textures{},
	scenes.Transformations{}.Name():  scenes.Transformations{},
	scenes.CoordinateSystem{}.Name(): scenes.CoordinateSystem{},
	scenes.Cube{}.Name():             scenes.Cube{},
	scenes.Camera{}.Name():           scenes.NewCamera(),
	scenes.LightColors{}.Name():      scenes.NewLightColors(),
}

func help() {
	fmt.Printf("scene name is required\n")
	fmt.Printf("possible values are: %q\n", possibleScenes())
	// TODO: Add message about possible controls (camera movement, screenshot etc)
}

func possibleScenes() []string {
	var names []string
	for k, _ := range allScenes {
		names = append(names, k)
	}
	return names
}
