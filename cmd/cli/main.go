package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/igoramorim/gopengl/internal/scenes"
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

	scene()
}

var allScenes = map[string]func(){
	"triangle": scenes.Triangle{}.Show,
	"shaders":  scenes.Shaders{}.Show,
	"textures": scenes.Textures{}.Show,
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
