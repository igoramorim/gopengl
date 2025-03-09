package scenes

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/igoramorim/gopengl/internal/sshot"
)

const (
	width  = 800
	height = 600

	floatSize  = 4
	uint32Size = 4
)

func processInput(w *glfw.Window, scene Scene) {
	if w.GetKey(glfw.KeyEscape) == glfw.Press {
		// Closes window
		w.SetShouldClose(true)
	}

	if w.GetKey(glfw.KeyL) == glfw.Press {
		// Enables wireframe drawing
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	}

	if w.GetKey(glfw.KeyF) == glfw.Press {
		// Disables wireframe drawing
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
	}

	if w.GetKey(glfw.KeyLeftControl) == glfw.Press && w.GetKey(glfw.KeyP) == glfw.Press {
		// Takes a screen shot
		sshoter := sshot.NewScreenShoter(scene.Name(), scene.Width(), scene.Height())
		sshoter.TakeOne()
	}
}

func frameBufferSizeCallback(w *glfw.Window, width, height int) {
	gl.Viewport(0, 0, int32(width), int32(height))
}
