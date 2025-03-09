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
)

func NewLightColors() LightColors {
	return LightColors{
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

// TODO: Create Camera struct
type LightColors struct {
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

func (s LightColors) Name() string {
	return "light_colors"
}

func (s LightColors) Width() int {
	return width
}

func (s LightColors) Height() int {
	return height
}

func (s LightColors) Show() {
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

	lightingShader, err := shader.New("internal/assets/shaders/light_colors.vert", "internal/assets/shaders/light_colors.frag")
	if err != nil {
		panic(err)
	}

	lightCubeShader, err := shader.New("internal/assets/shaders/light_colors_cube.vert", "internal/assets/shaders/light_colors_cube.frag")
	if err != nil {
		panic(err)
	}

	var vertices = []float32{
		// x y z
		-0.5, -0.5, -0.5,
		0.5, -0.5, -0.5,
		0.5, 0.5, -0.5,
		0.5, 0.5, -0.5,
		-0.5, 0.5, -0.5,
		-0.5, -0.5, -0.5,

		-0.5, -0.5, 0.5,
		0.5, -0.5, 0.5,
		0.5, 0.5, 0.5,
		0.5, 0.5, 0.5,
		-0.5, 0.5, 0.5,
		-0.5, -0.5, 0.5,

		-0.5, 0.5, 0.5,
		-0.5, 0.5, -0.5,
		-0.5, -0.5, -0.5,
		-0.5, -0.5, -0.5,
		-0.5, -0.5, 0.5,
		-0.5, 0.5, 0.5,

		0.5, 0.5, 0.5,
		0.5, 0.5, -0.5,
		0.5, -0.5, -0.5,
		0.5, -0.5, -0.5,
		0.5, -0.5, 0.5,
		0.5, 0.5, 0.5,

		-0.5, -0.5, -0.5,
		0.5, -0.5, -0.5,
		0.5, -0.5, 0.5,
		0.5, -0.5, 0.5,
		-0.5, -0.5, 0.5,
		-0.5, -0.5, -0.5,

		-0.5, 0.5, -0.5,
		0.5, 0.5, -0.5,
		0.5, 0.5, 0.5,
		0.5, 0.5, 0.5,
		-0.5, 0.5, 0.5,
		-0.5, 0.5, -0.5,
	}

	// First, configure the cubes's VAO and VBO
	var cubeVAO uint32
	gl.GenVertexArrays(1, &cubeVAO)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*floatSize, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.BindVertexArray(cubeVAO)

	// Position attribute
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*floatSize, nil)
	gl.EnableVertexAttribArray(0)

	// Second, configure the light's VAO
	// (VBO stays the same. The vertices are the same for the light object wich is also a 3D cube)
	var lightCubeVAO uint32
	gl.GenVertexArrays(1, &lightCubeVAO)
	gl.BindVertexArray(lightCubeVAO)

	// We only need to bind to the VBO (to link it with glVertexAttribPointer), no need to fill it;
	// The VBO's data already contains all we need (it's already bound, but we do it again for educational purposes)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)

	// Position attribute
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*floatSize, nil)
	gl.EnableVertexAttribArray(0)

	// Clean up all resources
	defer func() {
		gl.DeleteVertexArrays(1, &cubeVAO)
		gl.DeleteBuffers(1, &vbo)
		lightingShader.Delete()
		lightCubeShader.Delete()
	}()

	gl.Enable(gl.DEPTH_TEST)

	// Main loop
	for !window.ShouldClose() {
		s.processInput(window)

		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		currentFrame := glfw.GetTime()
		s.deltaTime = currentFrame - s.lastFrame
		s.lastFrame = currentFrame

		// Be sure to activate shader when settings uniforms/drawing objects
		lightingShader.Use()
		lightingShader.SetVec3f("objectColor", 1.0, 0.5, 0.31)
		lightingShader.SetVec3f("lightColor", 1.0, 1.0, 1.0)

		viewMatrix := mgl32.Ident4()
		viewMatrix = mgl32.LookAtV(s.cameraPos, s.cameraPos.Add(s.cameraFront), s.cameraUp)

		projectionMatrix := mgl32.Ident4()
		projectionMatrix = mgl32.Perspective(mgl32.DegToRad(float32(s.fov)), width/height, 0.1, 100.0)

		lightingShader.SetMat4("view", viewMatrix)
		lightingShader.SetMat4("projection", projectionMatrix)

		modelMatrix := mgl32.Ident4()
		lightingShader.SetMat4("model", modelMatrix)

		// Render the cube
		gl.BindVertexArray(cubeVAO)
		gl.DrawArrays(gl.TRIANGLES, 0, 36)

		// Now draw the cube "lamp"
		lightCubeShader.Use()
		lightCubeShader.SetMat4("projection", projectionMatrix)
		lightCubeShader.SetMat4("view", viewMatrix)

		modelMatrix = mgl32.Ident4()
		translate := mgl32.Translate3D(1.2, 1.2, 0.5)
		modelMatrix = modelMatrix.Mul4(translate)
		scale := mgl32.Scale3D(0.2, 0.2, 0.2)
		modelMatrix = modelMatrix.Mul4(scale)
		lightCubeShader.SetMat4("model", modelMatrix)

		gl.BindVertexArray(lightCubeVAO)
		gl.DrawArrays(gl.TRIANGLES, 0, 36)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func (s *LightColors) processInput(w *glfw.Window) {
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

func (s *LightColors) mouseCallback(w *glfw.Window, xpos, ypos float64) {
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

func (s *LightColors) mouseScrollCallback(w *glfw.Window, xoff, yoff float64) {
	s.fov -= yoff

	if s.fov < 1.0 {
		s.fov = 1.0
	}

	if s.fov > 45.0 {
		s.fov = 45.0
	}
}
