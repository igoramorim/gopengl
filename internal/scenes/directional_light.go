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

func NewDirectionalLight() DirectionalLight {
	return DirectionalLight{
		camera:     camera.New(),
		firstMouse: true,
		lastX:      float64(width) / 2,
		lastY:      float64(height) / 2,
		deltaTime:  0.0,
		lastFrame:  0.0,
	}
}

type DirectionalLight struct {
	camera     *camera.Camera
	firstMouse bool
	lastX      float64
	lastY      float64
	deltaTime  float64 // Time between current frame and last frame
	lastFrame  float64
}

func (s DirectionalLight) Name() string {
	return "directional_light"
}

func (s DirectionalLight) Width() int {
	return width
}

func (s DirectionalLight) Height() int {
	return height
}

func (s DirectionalLight) Show() {
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

	lightingShader, err := shader.New("internal/assets/shaders/directional_light.vert", "internal/assets/shaders/directional_light.frag")
	if err != nil {
		panic(err)
	}

	var vertices = []float32{
		// x y z          normals         texture coords
		-0.5, -0.5, -0.5, 0.0, 0.0, -1.0, 0.0, 0.0,
		0.5, -0.5, -0.5, 0.0, 0.0, -1.0, 1.0, 0.0,
		0.5, 0.5, -0.5, 0.0, 0.0, -1.0, 1.0, 1.0,
		0.5, 0.5, -0.5, 0.0, 0.0, -1.0, 1.0, 1.0,
		-0.5, 0.5, -0.5, 0.0, 0.0, -1.0, 0.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 0.0, -1.0, 0.0, 0.0,

		-0.5, -0.5, 0.5, 0.0, 0.0, 1.0, 0.0, 0.0,
		0.5, -0.5, 0.5, 0.0, 0.0, 1.0, 1.0, 0.0,
		0.5, 0.5, 0.5, 0.0, 0.0, 1.0, 1.0, 1.0,
		0.5, 0.5, 0.5, 0.0, 0.0, 1.0, 1.0, 1.0,
		-0.5, 0.5, 0.5, 0.0, 0.0, 1.0, 0.0, 1.0,
		-0.5, -0.5, 0.5, 0.0, 0.0, 1.0, 0.0, 0.0,

		-0.5, 0.5, 0.5, -1.0, 0.0, 0.0, 1.0, 0.0,
		-0.5, 0.5, -0.5, -1.0, 0.0, 0.0, 1.0, 1.0,
		-0.5, -0.5, -0.5, -1.0, 0.0, 0.0, 0.0, 1.0,
		-0.5, -0.5, -0.5, -1.0, 0.0, 0.0, 0.0, 1.0,
		-0.5, -0.5, 0.5, -1.0, 0.0, 0.0, 0.0, 0.0,
		-0.5, 0.5, 0.5, -1.0, 0.0, 0.0, 1.0, 0.0,

		0.5, 0.5, 0.5, 1.0, 0.0, 0.0, 1.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 0.0, 0.0, 1.0, 1.0,
		0.5, -0.5, -0.5, 1.0, 0.0, 0.0, 0.0, 1.0,
		0.5, -0.5, -0.5, 1.0, 0.0, 0.0, 0.0, 1.0,
		0.5, -0.5, 0.5, 1.0, 0.0, 0.0, 0.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 0.0, 0.0, 1.0, 0.0,

		-0.5, -0.5, -0.5, 0.0, -1.0, 0.0, 0.0, 1.0,
		0.5, -0.5, -0.5, 0.0, -1.0, 0.0, 1.0, 1.0,
		0.5, -0.5, 0.5, 0.0, -1.0, 0.0, 1.0, 0.0,
		0.5, -0.5, 0.5, 0.0, -1.0, 0.0, 1.0, 0.0,
		-0.5, -0.5, 0.5, 0.0, -1.0, 0.0, 0.0, 0.0,
		-0.5, -0.5, -0.5, 0.0, -1.0, 0.0, 0.0, 1.0,

		-0.5, 0.5, -0.5, 0.0, 1.0, 0.0, 0.0, 1.0,
		0.5, 0.5, -0.5, 0.0, 1.0, 0.0, 1.0, 1.0,
		0.5, 0.5, 0.5, 0.0, 1.0, 0.0, 1.0, 0.0,
		0.5, 0.5, 0.5, 0.0, 1.0, 0.0, 1.0, 0.0,
		-0.5, 0.5, 0.5, 0.0, 1.0, 0.0, 0.0, 0.0,
		-0.5, 0.5, -0.5, 0.0, 1.0, 0.0, 0.0, 1.0,
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
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 8*floatSize, nil)
	gl.EnableVertexAttribArray(0)
	// Normal attribute
	gl.VertexAttribPointerWithOffset(1, 3, gl.FLOAT, false, 8*floatSize, 3*floatSize)
	gl.EnableVertexAttribArray(1)
	// Texture Coord attribute
	gl.VertexAttribPointerWithOffset(2, 2, gl.FLOAT, false, 8*floatSize, 6*floatSize)
	gl.EnableVertexAttribArray(2)

	// Second, configure the light's VAO
	// (VBO stays the same. The vertices are the same for the light object wich is also a 3D cube)
	var lightCubeVAO uint32
	gl.GenVertexArrays(1, &lightCubeVAO)
	gl.BindVertexArray(lightCubeVAO)

	// We only need to bind to the VBO (to link it with glVertexAttribPointer), no need to fill it;
	// The VBO's data already contains all we need (it's already bound, but we do it again for educational purposes)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)

	// Position attribute
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 8*floatSize, nil)
	gl.EnableVertexAttribArray(0)

	diffuseMapTex, err := texture.New("internal/assets/textures/woodbox.png", gl.TEXTURE_2D, gl.TEXTURE0, gl.RGBA, gl.RGBA, gl.UNSIGNED_INT)
	if err != nil {
		panic(err)
	}

	specularMapTex, err := texture.New("internal/assets/textures/woodbox_specular.png", gl.TEXTURE_2D, gl.TEXTURE1, gl.RGBA, gl.RGBA, gl.UNSIGNED_INT)
	if err != nil {
		panic(err)
	}

	// Clean up all resources
	defer func() {
		gl.DeleteVertexArrays(1, &cubeVAO)
		gl.DeleteBuffers(1, &vbo)
		lightingShader.Delete()
		diffuseMapTex.Delete()
		specularMapTex.Delete()
	}()

	gl.Enable(gl.DEPTH_TEST)

	cubePositions := []mgl32.Vec3{
		mgl32.Vec3{0.0, 0.0, 0.0},
		mgl32.Vec3{2.0, 5.0, -15.0},
		mgl32.Vec3{-1.5, -2.2, -2.5},
		mgl32.Vec3{-3.8, -2.0, -12.3},
		mgl32.Vec3{2.4, -0.4, -3.5},
		mgl32.Vec3{-1.7, 3.0, -7.5},
		mgl32.Vec3{1.3, -2.0, -2.5},
		mgl32.Vec3{1.5, 2.0, -2.5},
		mgl32.Vec3{1.5, 0.2, -1.5},
		mgl32.Vec3{-1.3, 1.0, -1.5},
	}

	// Main loop
	for !window.ShouldClose() {
		s.processInput(window)

		gl.ClearColor(0.1, 0.1, 0.1, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		currentFrame := glfw.GetTime()
		s.deltaTime = currentFrame - s.lastFrame
		s.lastFrame = currentFrame

		viewMatrix := s.camera.ViewMatrix()
		projectionMatrix := mgl32.Ident4()
		projectionMatrix = mgl32.Perspective(mgl32.DegToRad(float32(s.camera.Fov)), width/height, 0.1, 100.0)

		diffuseMapTex.ActiveAndBind()
		specularMapTex.ActiveAndBind()

		lightingShader.Use()
		lightingShader.SetVec3f("light.direction", -0.2, -1.0, -0.3)
		lightingShader.SetVec3("viewPos", s.camera.Position)

		lightingShader.SetVec3f("light.ambient", 0.2, 0.2, 0.2)
		lightingShader.SetVec3f("light.diffuse", 0.5, 0.5, 0.5)
		lightingShader.SetVec3f("light.specular", 1.0, 1.0, 1.0)

		lightingShader.SetInt("material.diffuse", 0)
		lightingShader.SetInt("material.specular", 1)
		lightingShader.SetFloat("material.shininess", 32.0)

		lightingShader.SetMat4("view", viewMatrix)
		lightingShader.SetMat4("projection", projectionMatrix)

		// Render the cubes
		gl.BindVertexArray(cubeVAO)
		for i, pos := range cubePositions {
			modelMatrix := mgl32.Ident4()

			translate := mgl32.Translate3D(pos.X(), pos.Y(), pos.Z())
			angle := mgl32.DegToRad(20.0 * float32(i))
			rotateX := mgl32.HomogRotate3D(angle, mgl32.Vec3{1.0, 0.0, 0.0})
			rotateY := mgl32.HomogRotate3D(angle, mgl32.Vec3{0.0, 1.0, 0.0})
			rotateZ := mgl32.HomogRotate3D(angle, mgl32.Vec3{0.0, 0.0, 1.0})

			modelMatrix = modelMatrix.Mul4(translate)
			modelMatrix = modelMatrix.Mul4(rotateX)
			modelMatrix = modelMatrix.Mul4(rotateY)
			modelMatrix = modelMatrix.Mul4(rotateZ)

			lightingShader.SetMat4("model", modelMatrix)
			gl.DrawArrays(gl.TRIANGLES, 0, 36)
		}

		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func (s *DirectionalLight) processInput(w *glfw.Window) {
	processInput(w, s)
	processCameraKeyboardInput(w, s.camera, s.deltaTime)
}

func (s *DirectionalLight) mouseCallback(w *glfw.Window, xpos, ypos float64) {
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

func (s *DirectionalLight) mouseScrollCallback(w *glfw.Window, xoff, yoff float64) {
	s.camera.ProcessMouseScroll(yoff)
}
