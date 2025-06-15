package scenes

import (
	"fmt"
	"log"
	"math"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/igoramorim/gopengl/pkg/camera"
	"github.com/igoramorim/gopengl/pkg/shader"
)

func NewBasicLight() BasicLight {
	return BasicLight{
		camera:     camera.New(),
		firstMouse: true,
		lastX:      float64(width) / 2,
		lastY:      float64(height) / 2,
		deltaTime:  0.0,
		lastFrame:  0.0,
	}
}

type BasicLight struct {
	camera     *camera.Camera
	firstMouse bool
	lastX      float64
	lastY      float64
	deltaTime  float64 // Time between current frame and last frame
	lastFrame  float64
}

func (s BasicLight) Name() string {
	return "basic_light"
}

func (s BasicLight) Width() int {
	return width
}

func (s BasicLight) Height() int {
	return height
}

func (s BasicLight) Show() {
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

	lightingShader, err := shader.New("internal/assets/shaders/basic_light.vert", "internal/assets/shaders/basic_light.frag")
	if err != nil {
		panic(err)
	}

	lightCubeShader, err := shader.New("internal/assets/shaders/light_colors_cube.vert", "internal/assets/shaders/light_colors_cube.frag")
	if err != nil {
		panic(err)
	}

	var vertices = []float32{
		// x y z          normals
		-0.5, -0.5, -0.5, 0.0, 0.0, -1.0,
		0.5, -0.5, -0.5, 0.0, 0.0, -1.0,
		0.5, 0.5, -0.5, 0.0, 0.0, -1.0,
		0.5, 0.5, -0.5, 0.0, 0.0, -1.0,
		-0.5, 0.5, -0.5, 0.0, 0.0, -1.0,
		-0.5, -0.5, -0.5, 0.0, 0.0, -1.0,

		-0.5, -0.5, 0.5, 0.0, 0.0, 1.0,
		0.5, -0.5, 0.5, 0.0, 0.0, 1.0,
		0.5, 0.5, 0.5, 0.0, 0.0, 1.0,
		0.5, 0.5, 0.5, 0.0, 0.0, 1.0,
		-0.5, 0.5, 0.5, 0.0, 0.0, 1.0,
		-0.5, -0.5, 0.5, 0.0, 0.0, 1.0,

		-0.5, 0.5, 0.5, -1.0, 0.0, 0.0,
		-0.5, 0.5, -0.5, -1.0, 0.0, 0.0,
		-0.5, -0.5, -0.5, -1.0, 0.0, 0.0,
		-0.5, -0.5, -0.5, -1.0, 0.0, 0.0,
		-0.5, -0.5, 0.5, -1.0, 0.0, 0.0,
		-0.5, 0.5, 0.5, -1.0, 0.0, 0.0,

		0.5, 0.5, 0.5, 1.0, 0.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 0.0, 0.0,
		0.5, -0.5, -0.5, 1.0, 0.0, 0.0,
		0.5, -0.5, -0.5, 1.0, 0.0, 0.0,
		0.5, -0.5, 0.5, 1.0, 0.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 0.0, 0.0,

		-0.5, -0.5, -0.5, 0.0, -1.0, 0.0,
		0.5, -0.5, -0.5, 0.0, -1.0, 0.0,
		0.5, -0.5, 0.5, 0.0, -1.0, 0.0,
		0.5, -0.5, 0.5, 0.0, -1.0, 0.0,
		-0.5, -0.5, 0.5, 0.0, -1.0, 0.0,
		-0.5, -0.5, -0.5, 0.0, -1.0, 0.0,

		-0.5, 0.5, -0.5, 0.0, 1.0, 0.0,
		0.5, 0.5, -0.5, 0.0, 1.0, 0.0,
		0.5, 0.5, 0.5, 0.0, 1.0, 0.0,
		0.5, 0.5, 0.5, 0.0, 1.0, 0.0,
		-0.5, 0.5, 0.5, 0.0, 1.0, 0.0,
		-0.5, 0.5, -0.5, 0.0, 1.0, 0.0,
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
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 6*floatSize, nil)
	gl.EnableVertexAttribArray(0)
	// Normal attribute
	gl.VertexAttribPointerWithOffset(1, 3, gl.FLOAT, false, 6*floatSize, 3*floatSize)
	gl.EnableVertexAttribArray(1)

	// Second, configure the light's VAO
	// (VBO stays the same. The vertices are the same for the light object wich is also a 3D cube)
	var lightCubeVAO uint32
	gl.GenVertexArrays(1, &lightCubeVAO)
	gl.BindVertexArray(lightCubeVAO)

	// We only need to bind to the VBO (to link it with glVertexAttribPointer), no need to fill it;
	// The VBO's data already contains all we need (it's already bound, but we do it again for educational purposes)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)

	// Position attribute
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 6*floatSize, nil)
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

		gl.ClearColor(0.1, 0.1, 0.1, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		currentFrame := glfw.GetTime()
		s.deltaTime = currentFrame - s.lastFrame
		s.lastFrame = currentFrame

		lightPos := mgl32.Vec3{
			float32(math.Sin(glfw.GetTime())),
			1.0,
			float32(math.Cos(glfw.GetTime())),
		}

		lightColor := mgl32.Vec3{
			float32(math.Sin(glfw.GetTime())*0.5 + 0.5),
			1.0,
			float32(math.Cos(glfw.GetTime())*0.5 + 0.5),
			// 1.0, 1.0, 1.0,
		}

		lightingShader.Use()
		lightingShader.SetVec3f("objectColor", 1.0, 0.0, 1.0)
		lightingShader.SetVec3("lightColor", lightColor)
		lightingShader.SetVec3("lightPos", lightPos)
		lightingShader.SetVec3("viewPos", s.camera.Position)

		viewMatrix := s.camera.ViewMatrix()

		projectionMatrix := mgl32.Ident4()
		projectionMatrix = mgl32.Perspective(mgl32.DegToRad(float32(s.camera.Fov)), width/height, 0.1, 100.0)

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
		lightCubeShader.SetVec3("lightColor", lightColor)

		modelMatrix = mgl32.Ident4()
		translate := mgl32.Translate3D(lightPos.X(), lightPos.Y(), lightPos.Z())
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

func (s *BasicLight) processInput(w *glfw.Window) {
	processInput(w, s)
	processCameraKeyboardInput(w, s.camera, s.deltaTime)
}

func (s *BasicLight) mouseCallback(w *glfw.Window, xpos, ypos float64) {
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

func (s *BasicLight) mouseScrollCallback(w *glfw.Window, xoff, yoff float64) {
	s.camera.ProcessMouseScroll(yoff)
}
