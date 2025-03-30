package scenes

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/igoramorim/gopengl/internal/sshot"
	"github.com/igoramorim/gopengl/pkg/camera"
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

func processCameraKeyboardInput(w *glfw.Window, c *camera.Camera, deltaTime float64) {
	if w.GetKey(glfw.KeyW) == glfw.Press {
		c.ProcessKeyboard(camera.Forward, deltaTime)
	}

	if w.GetKey(glfw.KeyS) == glfw.Press {
		c.ProcessKeyboard(camera.Backward, deltaTime)
	}

	if w.GetKey(glfw.KeyA) == glfw.Press {
		c.ProcessKeyboard(camera.Left, deltaTime)
	}

	if w.GetKey(glfw.KeyD) == glfw.Press {
		c.ProcessKeyboard(camera.Right, deltaTime)
	}

	if w.GetKey(glfw.KeyQ) == glfw.Press {
		c.ProcessKeyboard(camera.Down, deltaTime)
	}

	if w.GetKey(glfw.KeyE) == glfw.Press {
		c.ProcessKeyboard(camera.Up, deltaTime)
	}

	const rotate = 5.0
	if w.GetKey(glfw.KeyLeft) == glfw.Press {
		c.ProcessMouseMovement(-rotate, 0.0, true)
	}

	if w.GetKey(glfw.KeyRight) == glfw.Press {
		c.ProcessMouseMovement(rotate, 0.0, true)
	}

	if w.GetKey(glfw.KeyUp) == glfw.Press {
		c.ProcessMouseMovement(0.0, rotate, true)
	}

	if w.GetKey(glfw.KeyDown) == glfw.Press {
		c.ProcessMouseMovement(0.0, -rotate, true)
	}
}

func frameBufferSizeCallback(w *glfw.Window, width, height int) {
	gl.Viewport(0, 0, int32(width), int32(height))
}
