package scenes

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	width  = 800
	height = 600

	floatSize  = 4
	uint32Size = 4
)

func processInput(w *glfw.Window) {
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
}

func frameBufferSizeCallback(w *glfw.Window, width, height int) {
	gl.Viewport(0, 0, int32(width), int32(height))
}
