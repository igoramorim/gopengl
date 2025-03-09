package scenes

import (
	"fmt"
	"log"
	"math"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/igoramorim/gopengl/pkg/shader"
	"github.com/igoramorim/gopengl/pkg/texture"
)

func NewCamera() Camera {
	return Camera{
		cameraPos:   mgl32.Vec3{0.0, 0.0, 3.0},
		cameraFront: mgl32.Vec3{0.0, 0.0, -1.0},
		cameraUp:    mgl32.Vec3{0.0, 1.0, 0.0},
		firstMouse:  true,
		lastX:       float64(width) / 2,
		lastY:       float64(height) / 2,
		yaw:         -90.0,
		pitch:       0.0,
		fov:         45.0,
		deltaTime:   0.0,
		lastFrame:   0.0,
	}
}

type Camera struct {
	cameraPos   mgl32.Vec3
	cameraFront mgl32.Vec3
	cameraUp    mgl32.Vec3
	firstMouse  bool
	lastX       float64
	lastY       float64
	yaw         float64
	pitch       float64
	fov         float64
	deltaTime   float64 // Time between current frame and last frame
	lastFrame   float64
}

func (s Camera) Name() string {
	return "camera"
}

func (s Camera) Width() int {
	return width
}

func (s Camera) Height() int {
	return height
}

func (s Camera) Show() {
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

	// Handles mouse position. Calls mouseCallback every time the cursor moves
	window.SetCursorPosCallback(s.mouseCallback)
	// Handles mouse scroll. Calls mouseScrollCallback evert time the scrolling is used
	window.SetScrollCallback(s.mouseScrollCallback)
	// Hides the mouse cursor
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version:", version)

	shader, err := shader.New("internal/assets/shaders/camera.vert", "internal/assets/shaders/camera.frag")
	if err != nil {
		panic(err)
	}

	var vertices = []float32{
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

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*floatSize, gl.Ptr(vertices), gl.STATIC_DRAW)

	// Position attribute
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*floatSize, nil)
	gl.EnableVertexAttribArray(0)

	// Texture Coord attribute
	gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, 5*floatSize, 3*floatSize)
	gl.EnableVertexAttribArray(1)

	gl.BindVertexArray(vao)

	texture0, err := texture.New("internal/assets/textures/container.jpg", gl.TEXTURE_2D, gl.TEXTURE0, gl.RGBA, gl.RGBA, gl.UNSIGNED_INT)
	if err != nil {
		panic(err)
	}

	texture1, err := texture.New("internal/assets/textures/awesomeface.png", gl.TEXTURE_2D, gl.TEXTURE1, gl.RGBA, gl.RGBA, gl.UNSIGNED_INT)
	if err != nil {
		panic(err)
	}

	// Clean up all resources
	defer func() {
		gl.DeleteVertexArrays(1, &vao)
		gl.DeleteBuffers(1, &vbo)
		shader.Delete()
		texture0.Delete()
		texture1.Delete()
	}()

	shader.Use()

	texture0.ActiveAndBind()
	shader.SetInt("texture0", 0)

	texture1.ActiveAndBind()
	shader.SetInt("texture1", 1)

	gl.Enable(gl.DEPTH_TEST)

	cubePositions := []mgl32.Vec3{
		mgl32.Vec3{0.0, 0.0, 0.0},
		mgl32.Vec3{2.0, 5.0, -15.0},
		mgl32.Vec3{-2.0, 5.0, -15.0},
		mgl32.Vec3{2.0, -5.0, -15.0},
		mgl32.Vec3{-2.0, -5.0, -15.0},
		mgl32.Vec3{2.0, 2.5, -5.0},
		mgl32.Vec3{-2.0, 2.5, -5.0},
		mgl32.Vec3{2.0, -2.5, -5.0},
		mgl32.Vec3{-2.0, -2.5, -5.0},
	}

	// Main loop
	for !window.ShouldClose() {
		s.processInput(window)

		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Per frame logic
		currentFrame := glfw.GetTime()
		s.deltaTime = currentFrame - s.lastFrame
		s.lastFrame = currentFrame

		time := glfw.GetTime()

		// Transformations to make it 3D

		viewMatrix := mgl32.Ident4()
		viewMatrix = mgl32.LookAtV(
			s.cameraPos,                    // Position of the camera - 'eye'
			s.cameraPos.Add(s.cameraFront), // Target position
			s.cameraUp,                     // Vector that points UP in the world space
		)

		projectionMatrix := mgl32.Ident4()
		projectionMatrix = mgl32.Perspective(mgl32.DegToRad(float32(s.fov)), width/height, 0.1, 100.0)

		shader.SetMat4("view", viewMatrix)
		shader.SetMat4("projection", projectionMatrix)

		for _, cubePos := range cubePositions {
			modelMatrix := mgl32.Ident4()

			translate := mgl32.Translate3D(cubePos.X(), cubePos.Y(), cubePos.Z())
			angle := float32(time) * mgl32.DegToRad(45.0)
			rotateX := mgl32.HomogRotate3D(angle, mgl32.Vec3{1.0, 0.0, 0.0})
			rotateY := mgl32.HomogRotate3D(angle, mgl32.Vec3{0.0, 1.0, 0.0})
			rotateZ := mgl32.HomogRotate3D(angle, mgl32.Vec3{0.0, 0.0, 1.0})

			modelMatrix = modelMatrix.Mul4(translate)
			modelMatrix = modelMatrix.Mul4(rotateX)
			modelMatrix = modelMatrix.Mul4(rotateY)
			modelMatrix = modelMatrix.Mul4(rotateZ)

			shader.SetMat4("model", modelMatrix)

			gl.DrawArrays(gl.TRIANGLES, 0, 36) // 36 vertices to make a cube
		}

		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func (s *Camera) processInput(w *glfw.Window) {
	processInput(w, s)

	// deltaTime used to make speed consistency among different hardware setups
	cameraSpeed := 2.5 * float32(s.deltaTime)
	if w.GetKey(glfw.KeyW) == glfw.Press {
		s.cameraPos = s.cameraPos.Add(s.cameraFront.Mul(cameraSpeed))
	}

	if w.GetKey(glfw.KeyS) == glfw.Press {
		s.cameraPos = s.cameraPos.Sub(s.cameraFront.Mul(cameraSpeed))
	}

	if w.GetKey(glfw.KeyA) == glfw.Press {
		s.cameraPos = s.cameraPos.Sub(s.cameraFront.Cross(s.cameraUp).Normalize().Mul(cameraSpeed))
	}

	if w.GetKey(glfw.KeyD) == glfw.Press {
		s.cameraPos = s.cameraPos.Add(s.cameraFront.Cross(s.cameraUp).Normalize().Mul(cameraSpeed))
	}
}

func (s *Camera) mouseCallback(w *glfw.Window, xpos, ypos float64) {
	if s.firstMouse {
		s.lastX = xpos
		s.lastY = ypos
		s.firstMouse = false
	}

	xoffset := xpos - s.lastX
	yoffset := s.lastY - ypos
	s.lastX = xpos
	s.lastY = ypos

	sensitivity := 0.1
	xoffset *= sensitivity
	yoffset *= sensitivity

	s.yaw += xoffset
	s.pitch += yoffset

	if s.pitch > 89.0 {
		s.pitch = 89.0
	}

	if s.pitch < -89.0 {
		s.pitch = -89.0
	}

	direction := mgl32.Vec3{
		float32(math.Cos(mgl64.DegToRad(s.yaw)) * math.Cos(mgl64.DegToRad(s.pitch))),
		float32(math.Sin(mgl64.DegToRad(s.pitch))),
		float32(math.Sin(mgl64.DegToRad(s.yaw)) * math.Cos(mgl64.DegToRad(s.pitch))),
	}
	s.cameraFront = direction.Normalize()
}

func (s *Camera) mouseScrollCallback(w *glfw.Window, xoff, yoff float64) {
	s.fov -= yoff

	if s.fov < 1.0 {
		s.fov = 1.0
	}

	if s.fov > 45.0 {
		s.fov = 45.0
	}
}
