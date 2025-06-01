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

func NewDepthTesting() DepthTesting {
	return DepthTesting{
		camera:     camera.New(),
		firstMouse: true,
		lastX:      float64(width) / 2,
		lastY:      float64(height) / 2,
		deltaTime:  0.0,
		lastFrame:  0.0,
	}
}

type DepthTesting struct {
	camera     *camera.Camera
	firstMouse bool
	lastX      float64
	lastY      float64
	deltaTime  float64 // Time between current frame and last frame
	lastFrame  float64
}

func (s DepthTesting) Name() string {
	return "depth_testing"
}

func (s DepthTesting) Width() int {
	return width
}

func (s DepthTesting) Height() int {
	return height
}

func (s DepthTesting) Show() {
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

	shader, err := shader.New("internal/assets/shaders/depth_testing.vert", "internal/assets/shaders/depth_testing.frag")
	if err != nil {
		panic(err)
	}

	// var cubeVertices = []float32{
	// 	// x y z u v (tex coord)
	// 	-0.5, -0.5, -0.5, 0.0, 0.0,
	// 	0.5, -0.5, -0.5, 1.0, 0.0,
	// 	0.5, 0.5, -0.5, 1.0, 1.0,
	// 	0.5, 0.5, -0.5, 1.0, 1.0,
	// 	-0.5, 0.5, -0.5, 0.0, 1.0,
	// 	-0.5, -0.5, -0.5, 0.0, 0.0,
	//
	// 	-0.5, -0.5, 0.5, 0.0, 0.0,
	// 	0.5, -0.5, 0.5, 1.0, 0.0,
	// 	0.5, 0.5, 0.5, 1.0, 1.0,
	// 	0.5, 0.5, 0.5, 1.0, 1.0,
	// 	-0.5, 0.5, 0.5, 0.0, 1.0,
	// 	-0.5, -0.5, 0.5, 0.0, 0.0,
	//
	// 	-0.5, 0.5, 0.5, 1.0, 0.0,
	// 	-0.5, 0.5, -0.5, 1.0, 1.0,
	// 	-0.5, -0.5, -0.5, 0.0, 1.0,
	// 	-0.5, -0.5, -0.5, 0.0, 1.0,
	// 	-0.5, -0.5, 0.5, 0.0, 0.0,
	// 	-0.5, 0.5, 0.5, 1.0, 0.0,
	//
	// 	0.5, 0.5, 0.5, 1.0, 0.0,
	// 	0.5, 0.5, -0.5, 1.0, 1.0,
	// 	0.5, -0.5, -0.5, 0.0, 1.0,
	// 	0.5, -0.5, -0.5, 0.0, 1.0,
	// 	0.5, -0.5, 0.5, 0.0, 0.0,
	// 	0.5, 0.5, 0.5, 1.0, 0.0,
	//
	// 	-0.5, -0.5, -0.5, 0.0, 1.0,
	// 	0.5, -0.5, -0.5, 1.0, 1.0,
	// 	0.5, -0.5, 0.5, 1.0, 0.0,
	// 	0.5, -0.5, 0.5, 1.0, 0.0,
	// 	-0.5, -0.5, 0.5, 0.0, 0.0,
	// 	-0.5, -0.5, -0.5, 0.0, 1.0,
	//
	// 	-0.5, 0.5, -0.5, 0.0, 1.0,
	// 	0.5, 0.5, -0.5, 1.0, 1.0,
	// 	0.5, 0.5, 0.5, 1.0, 0.0,
	// 	0.5, 0.5, 0.5, 1.0, 0.0,
	// 	-0.5, 0.5, 0.5, 0.0, 0.0,
	// 	-0.5, 0.5, -0.5, 0.0, 1.0,
	// }

	// var planeVertices = []float32{
	// 	// x y z  uv (tex coords) (note we set these higher than 1 (together with GL_REPEAT as texture wrapping mode). this will cause the floor texture to repeat)
	// 	5.0, -0.5, 5.0, 2.0, 0.0,
	// 	-5.0, -0.5, 5.0, 0.0, 0.0,
	// 	-5.0, -0.5, -5.0, 0.0, 2.0,
	// 	5.0, -0.5, 5.0, 2.0, 0.0,
	// 	-5.0, -0.5, -5.0, 0.0, 2.0,
	// 	5.0, -0.5, -5.0, 2.0, 2.0,
	// }

	// Cube
	// var cubeVAO, cubeVBO uint32
	// gl.GenVertexArrays(1, &cubeVAO)
	// gl.GenBuffers(1, &cubeVBO)
	// gl.BindVertexArray(cubeVAO)
	// gl.BindBuffer(gl.ARRAY_BUFFER, cubeVBO)
	// gl.BufferData(gl.ARRAY_BUFFER, len(cubeVertices)*floatSize, gl.Ptr(cubeVertices), gl.STATIC_DRAW)
	// // Position attribute
	// gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*floatSize, nil)
	// gl.EnableVertexAttribArray(0)
	// // Texture Coord attribute
	// gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, 5*floatSize, 3*floatSize)
	// gl.EnableVertexAttribArray(1)
	// gl.BindVertexArray(0)

	// Plane
	// var planeVAO, planeVBO uint32
	// gl.GenVertexArrays(1, &planeVAO)
	// gl.GenBuffers(1, &planeVBO)
	// gl.BindVertexArray(planeVAO)
	// gl.BindBuffer(gl.ARRAY_BUFFER, planeVBO)
	// gl.BufferData(gl.ARRAY_BUFFER, len(planeVertices)*floatSize, gl.Ptr(planeVertices), gl.STATIC_DRAW)
	// // Position attribute
	// gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*floatSize, nil)
	// gl.EnableVertexAttribArray(0)
	// // Texture Coord attribute
	// gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, 5*floatSize, 3*floatSize)
	// gl.EnableVertexAttribArray(1)
	// gl.BindVertexArray(0)

	// Textures
	// cubeTexture, err := texture.New("internal/assets/textures/marble.jpg", gl.TEXTURE_2D, gl.TEXTURE0, gl.RGBA, gl.RGBA, gl.UNSIGNED_INT)
	// if err != nil {
	// 	panic(err)
	// }

	// floorTexture, err := texture.New("internal/assets/textures/metal.png", gl.TEXTURE_2D, gl.TEXTURE1, gl.RGBA, gl.RGBA, gl.UNSIGNED_INT)
	// if err != nil {
	// 	panic(err)
	// }

	model3D, err := model.New("internal/assets/models/sponza/sponza.obj")
	if err != nil {
		panic(err)
	}

	// Clean up all resources
	defer func() {
		shader.Delete()
		// cubeTexture.Delete()
		// floorTexture.Delete()
	}()

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS) // gl.ALWAYS

	// Main loop
	for !window.ShouldClose() {
		s.processInput(window)

		gl.ClearColor(0.1, 0.1, 0.1, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		currentFrame := glfw.GetTime()
		s.deltaTime = currentFrame - s.lastFrame
		s.lastFrame = currentFrame

		shader.Use()

		viewMatrix := s.camera.ViewMatrix()
		projectionMatrix := mgl32.Ident4()
		projectionMatrix = mgl32.Perspective(mgl32.DegToRad(float32(s.camera.Fov)), width/height, 0.1, 100.0)
		shader.SetMat4("view", viewMatrix)
		shader.SetMat4("projection", projectionMatrix)

		// // Cube 1
		// cubeTexture.ActiveAndBind()
		// shader.SetInt("texture0", 0)
		// gl.BindVertexArray(cubeVAO)
		// modelMatrix := mgl32.Ident4()
		// translate := mgl32.Translate3D(-1.0, 0.0, -1.0)
		// modelMatrix = modelMatrix.Mul4(translate)
		// shader.SetMat4("model", modelMatrix)
		// gl.DrawArrays(gl.TRIANGLES, 0, 36)
		//
		// // Cube 2
		// // gl.BindVertexArray(cubeVAO)
		// modelMatrix := mgl32.Ident4()
		// translate = mgl32.Translate3D(2.0, 0.0, 0.0)
		// modelMatrix = modelMatrix.Mul4(translate)
		// shader.SetMat4("model", modelMatrix)
		// gl.DrawArrays(gl.TRIANGLES, 0, 36)
		//
		// // Floor
		// floorTexture.ActiveAndBind()
		// shader.SetInt("texture0", 1)
		// gl.BindVertexArray(planeVAO)
		// shader.SetMat4("model", mgl32.Ident4())
		// gl.DrawArrays(gl.TRIANGLES, 0, 6)
		// gl.BindVertexArray(0)

		modelMatrix := mgl32.Ident4()
		translate := mgl32.Translate3D(0, 0, 0)
		modelMatrix = modelMatrix.Mul4(translate)
		// Sponza model is large. Scale it down so we can navigate better
		scale := mgl32.Scale3D(0.01, 0.01, 0.01)
		modelMatrix = modelMatrix.Mul4(scale)
		shader.SetMat4("model", modelMatrix)
		model3D.Draw(shader)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func (s *DepthTesting) processInput(w *glfw.Window) {
	processInput(w, s)
	processCameraKeyboardInput(w, s.camera, s.deltaTime)
}

func (s *DepthTesting) mouseCallback(w *glfw.Window, xpos, ypos float64) {
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

func (s *DepthTesting) mouseScrollCallback(w *glfw.Window, xoff, yoff float64) {
	s.camera.ProcessMouseScroll(yoff)
}
