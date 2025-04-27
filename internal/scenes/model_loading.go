package scenes

import (
	"fmt"
	"log"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/igoramorim/gopengl/pkg/camera"
	"github.com/igoramorim/gopengl/pkg/model"
	"github.com/igoramorim/gopengl/pkg/shader"
)

func NewModelLoading() ModelLoading {
	return ModelLoading{
		camera:     camera.New(),
		firstMouse: true,
		lastX:      float64(width) / 2,
		lastY:      float64(height) / 2,
		deltaTime:  0.0,
		lastFrame:  0.0,
	}
}

type ModelLoading struct {
	camera     *camera.Camera
	firstMouse bool
	lastX      float64
	lastY      float64
	deltaTime  float64 // Time between current frame and last frame
	lastFrame  float64
}

func (s ModelLoading) Name() string {
	return "model_loading"
}

func (s ModelLoading) Width() int {
	return width
}

func (s ModelLoading) Height() int {
	return height
}

func (s ModelLoading) Show() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(s.Width(), s.Height(), s.Name(), nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	window.SetFramebufferSizeCallback(frameBufferSizeCallback)
	window.SetCursorPosCallback(s.mouseCallback)
	window.SetScrollCallback(s.mouseScrollCallback)
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version:", version)

	modelShader, err := shader.New("internal/assets/shaders/model_loading.vert", "internal/assets/shaders/model_loading.frag")
	if err != nil {
		panic(err)
	}

	model3D, err := model.New("internal/assets/models/backpack/backpack.obj")
	if err != nil {
		panic(err)
	}

	// Clean up all resources
	defer func() {
		modelShader.Delete()
	}()

	gl.Enable(gl.DEPTH_TEST)

	// Main loop
	for !window.ShouldClose() {
		s.processInput(window)

		gl.ClearColor(1.0, 0.05, 1.0, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		currentFrame := glfw.GetTime()
		s.deltaTime = currentFrame - s.lastFrame
		s.lastFrame = currentFrame

		modelShader.Use()

		viewMatrix := s.camera.ViewMatrix()
		projectionMatrix := mgl32.Ident4()
		projectionMatrix = mgl32.Perspective(mgl32.DegToRad(float32(s.camera.Fov)), width/height, 0.1, 100.0)
		modelShader.SetMat4("view", viewMatrix)
		modelShader.SetMat4("projection", projectionMatrix)

		modelMatrix := mgl32.Ident4()
		modelShader.SetMat4("model", modelMatrix)
		model3D.Draw(modelShader)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func (s *ModelLoading) processInput(w *glfw.Window) {
	processInput(w, s)
	processCameraKeyboardInput(w, s.camera, s.deltaTime)
}

func (s *ModelLoading) mouseCallback(w *glfw.Window, xpos, ypos float64) {
	if s.firstMouse {
		s.lastX = xpos
		s.lastY = ypos
		s.firstMouse = false
	}

	xoffset := xpos - s.lastX
	yoffset := s.lastY - ypos
	s.lastX = xpos
	s.lastY = ypos

	s.camera.ProcessMouseMovement(xoffset, yoffset, true)
}

func (s *ModelLoading) mouseScrollCallback(w *glfw.Window, xoff, yoff float64) {
	s.camera.ProcessMouseScroll(yoff)
}
