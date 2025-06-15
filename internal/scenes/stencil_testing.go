package scenes

import (
	"fmt"
	"log"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/igoramorim/gopengl/pkg/camera"
	"github.com/igoramorim/gopengl/pkg/shader"
	"github.com/igoramorim/gopengl/pkg/texture"
)

func NewStencilTesting() StencilTesting {
	return StencilTesting{
		camera:     camera.New(),
		firstMouse: true,
		lastX:      float64(width) / 2,
		lastY:      float64(height) / 2,
		deltaTime:  0.0,
		lastFrame:  0.0,
	}
}

type StencilTesting struct {
	camera     *camera.Camera
	firstMouse bool
	lastX      float64
	lastY      float64
	deltaTime  float64 // Time between current frame and last frame
	lastFrame  float64
}

func (s StencilTesting) Name() string {
	return "stencil_testing"
}

func (s StencilTesting) Width() int {
	return width
}

func (s StencilTesting) Height() int {
	return height
}

func (s StencilTesting) Show() {
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

	shaderObject, err := shader.New("internal/assets/shaders/stencil_testing.vert", "internal/assets/shaders/stencil_testing.frag")
	if err != nil {
		panic(err)
	}

	shaderBorder, err := shader.New("internal/assets/shaders/stencil_testing.vert", "internal/assets/shaders/stencil_testing_border.frag")
	if err != nil {
		panic(err)
	}

	var cubeVertices = []float32{
		// x y z u v (tex coord)
		-0.5, -0.5, -0.5, 0.0, 0.0,
		0.5, -0.5, -0.5, 1.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		-0.5, 0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 0.0,

		-0.5, -0.5, 0.5, 0.0, 0.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 1.0,
		0.5, 0.5, 0.5, 1.0, 1.0,
		-0.5, 0.5, 0.5, 0.0, 1.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,

		-0.5, 0.5, 0.5, 1.0, 0.0,
		-0.5, 0.5, -0.5, 1.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,
		-0.5, 0.5, 0.5, 1.0, 0.0,

		0.5, 0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, 0.5, 0.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 0.0,

		-0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, -0.5, 1.0, 1.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,

		-0.5, 0.5, -0.5, 0.0, 1.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, 0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 0.0,
		-0.5, 0.5, 0.5, 0.0, 0.0,
		-0.5, 0.5, -0.5, 0.0, 1.0,
	}

	var planeVertices = []float32{
		// x y z  uv (tex coords) (note we set these higher than 1 (together with GL_REPEAT as texture wrapping mode). this will cause the floor texture to repeat)
		5.0, -0.5, 5.0, 2.0, 0.0,
		-5.0, -0.5, 5.0, 0.0, 0.0,
		-5.0, -0.5, -5.0, 0.0, 2.0,
		5.0, -0.5, 5.0, 2.0, 0.0,
		-5.0, -0.5, -5.0, 0.0, 2.0,
		5.0, -0.5, -5.0, 2.0, 2.0,
	}

	// Cube
	var cubeVAO, cubeVBO uint32
	gl.GenVertexArrays(1, &cubeVAO)
	gl.GenBuffers(1, &cubeVBO)
	gl.BindVertexArray(cubeVAO)
	gl.BindBuffer(gl.ARRAY_BUFFER, cubeVBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(cubeVertices)*floatSize, gl.Ptr(cubeVertices), gl.STATIC_DRAW)
	// Position attribute
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*floatSize, nil)
	gl.EnableVertexAttribArray(0)
	// Texture Coord attribute
	gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, 5*floatSize, 3*floatSize)
	gl.EnableVertexAttribArray(1)
	gl.BindVertexArray(0)

	// Plane
	var planeVAO, planeVBO uint32
	gl.GenVertexArrays(1, &planeVAO)
	gl.GenBuffers(1, &planeVBO)
	gl.BindVertexArray(planeVAO)
	gl.BindBuffer(gl.ARRAY_BUFFER, planeVBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(planeVertices)*floatSize, gl.Ptr(planeVertices), gl.STATIC_DRAW)
	// Position attribute
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*floatSize, nil)
	gl.EnableVertexAttribArray(0)
	// Texture Coord attribute
	gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, 5*floatSize, 3*floatSize)
	gl.EnableVertexAttribArray(1)
	gl.BindVertexArray(0)

	// Textures
	cubeTexture, err := texture.New("internal/assets/textures/marble.jpg", gl.TEXTURE_2D, gl.TEXTURE0, gl.RGBA, gl.RGBA, gl.UNSIGNED_INT)
	if err != nil {
		panic(err)
	}

	floorTexture, err := texture.New("internal/assets/textures/metal.png", gl.TEXTURE_2D, gl.TEXTURE1, gl.RGBA, gl.RGBA, gl.UNSIGNED_INT)
	if err != nil {
		panic(err)
	}

	// Clean up all resources
	defer func() {
		shaderObject.Delete()
		shaderBorder.Delete()
		cubeTexture.Delete()
		floorTexture.Delete()
	}()

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.Enable(gl.STENCIL_TEST)
	gl.StencilFunc(gl.NOTEQUAL, 1, 0xFF)
	gl.StencilOp(gl.KEEP, gl.KEEP, gl.REPLACE)

	// Main loop
	for !window.ShouldClose() {
		s.processInput(window)

		gl.ClearColor(0.1, 0.1, 0.1, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT | gl.STENCIL_BUFFER_BIT)

		currentFrame := glfw.GetTime()
		s.deltaTime = currentFrame - s.lastFrame
		s.lastFrame = currentFrame

		viewMatrix := s.camera.ViewMatrix()
		projectionMatrix := mgl32.Ident4()
		projectionMatrix = mgl32.Perspective(mgl32.DegToRad(float32(s.camera.Fov)), width/height, 0.1, 100.0)

		shaderBorder.Use()
		shaderBorder.SetMat4("view", viewMatrix)
		shaderBorder.SetMat4("projection", projectionMatrix)

		shaderObject.Use()
		shaderObject.SetMat4("view", viewMatrix)
		shaderObject.SetMat4("projection", projectionMatrix)

		// Floor
		// Draw the floor but do not write to the stencil buffer by setting its mask to 0x00
		gl.StencilMask(0x00)
		floorTexture.ActiveAndBind()
		shaderObject.SetInt("texture0", 1)
		gl.BindVertexArray(planeVAO)
		shaderObject.SetMat4("model", mgl32.Ident4())
		gl.DrawArrays(gl.TRIANGLES, 0, 6)
		gl.BindVertexArray(0)

		// 1st render pass
		// Draw objects as normal, writing to the stencil buffer
		gl.StencilFunc(gl.ALWAYS, 1, 0xFF)
		gl.StencilMask(0xFF)

		// Cube 1
		cubeTexture.ActiveAndBind()
		shaderObject.SetInt("texture0", 0)
		gl.BindVertexArray(cubeVAO)
		modelMatrix := mgl32.Ident4()
		translate := mgl32.Translate3D(-1.0, 0.0, -1.0)
		modelMatrix = modelMatrix.Mul4(translate)
		shaderObject.SetMat4("model", modelMatrix)
		gl.DrawArrays(gl.TRIANGLES, 0, 36)

		// Cube 2
		gl.BindVertexArray(cubeVAO)
		modelMatrix = mgl32.Ident4()
		translate = mgl32.Translate3D(2.0, 0.0, 0.0)
		modelMatrix = modelMatrix.Mul4(translate)
		shaderObject.SetMat4("model", modelMatrix)
		gl.DrawArrays(gl.TRIANGLES, 0, 36)

		// 2nd render pass
		// Draw scaled version of the objects, this time without writing to the stencil buffer
		// Because the stencil buffer is now filled with 1s, the parts of the buffer that are
		// 1 are not draw, thus only drawing the objects size differences, making it look like a border
		gl.StencilFunc(gl.NOTEQUAL, 1, 0xFF)
		gl.StencilMask(0x00)
		gl.Disable(gl.DEPTH_TEST)
		shaderBorder.Use()

		var scale float32
		scale = 1.1
		// Cube 1
		cubeTexture.ActiveAndBind()
		shaderBorder.SetInt("texture0", 0)
		gl.BindVertexArray(cubeVAO)
		modelMatrix = mgl32.Ident4()
		modelMatrix = modelMatrix.Mul4(mgl32.Translate3D(-1.0, 0.0, -1.0))
		modelMatrix = modelMatrix.Mul4(mgl32.Scale3D(scale, scale, scale))
		shaderBorder.SetMat4("model", modelMatrix)
		gl.DrawArrays(gl.TRIANGLES, 0, 36)

		// Cube 2
		gl.BindVertexArray(cubeVAO)
		modelMatrix = mgl32.Ident4()
		modelMatrix = modelMatrix.Mul4(mgl32.Translate3D(2.0, 0.0, 0.0))
		modelMatrix = modelMatrix.Mul4(mgl32.Scale3D(scale, scale, scale))
		shaderBorder.SetMat4("model", modelMatrix)
		gl.DrawArrays(gl.TRIANGLES, 0, 36)

		gl.BindVertexArray(0)
		gl.StencilMask(0xFF)
		gl.StencilFunc(gl.ALWAYS, 0, 0xFF)
		gl.Enable(gl.DEPTH_TEST)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func (s *StencilTesting) processInput(w *glfw.Window) {
	processInput(w, s)
	processCameraKeyboardInput(w, s.camera, s.deltaTime)
}

func (s *StencilTesting) mouseCallback(w *glfw.Window, xpos, ypos float64) {
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

func (s *StencilTesting) mouseScrollCallback(w *glfw.Window, xoff, yoff float64) {
	s.camera.ProcessMouseScroll(yoff)
}
