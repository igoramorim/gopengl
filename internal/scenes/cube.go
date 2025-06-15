package scenes

import (
	"fmt"
	"log"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/igoramorim/gopengl/pkg/shader"
	"github.com/igoramorim/gopengl/pkg/texture"
)

type Cube struct{}

func (s Cube) Name() string {
	return "cube"
}

func (s Cube) Width() int {
	return width
}

func (s Cube) Height() int {
	return height
}

func (s Cube) Show() {
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

	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version:", version)

	shader, err := shader.New("internal/assets/shaders/coordinate-system.vert", "internal/assets/shaders/coordinate-system.frag")
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

	// OpenGL stores all depth information in the z-buffer (depth buffer). The depth is stored within each fragment.
	// When the fragment wants to output a color, OpenGL compares its depth value with the z-buffer.
	// If the current fragment is behind the other fragment it is discarded, otherwise overwritten.
	// GL_DEPTH_BUFFER_BIT is also needed in the glClear to make it work.

	// Comment the glEnable below and see that some faces of the cube are above others.
	// That happens because OpenGL does not do the depth test, and since OpenGL does not
	// guarantee the order the triangles are rendered, some triangles are drawn on top of each other.
	gl.Enable(gl.DEPTH_TEST)

	// Main loop
	for !window.ShouldClose() {
		processInput(window, s)

		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		time := glfw.GetTime()

		// Transformations to make it 3D

		modelMatrix := mgl32.Ident4()
		angle := float32(time) * mgl32.DegToRad(45.0)
		rotateX := mgl32.HomogRotate3D(angle, mgl32.Vec3{1.0, 0.0, 0.0})
		rotateY := mgl32.HomogRotate3D(angle, mgl32.Vec3{0.0, 1.0, 0.0})
		rotateZ := mgl32.HomogRotate3D(angle, mgl32.Vec3{0.0, 0.0, 1.0})
		modelMatrix = modelMatrix.Mul4(rotateX)
		modelMatrix = modelMatrix.Mul4(rotateY)
		modelMatrix = modelMatrix.Mul4(rotateZ)

		viewMatrix := mgl32.Ident4()
		translate := mgl32.Translate3D(0.0, 0.0, -4.0)
		viewMatrix = viewMatrix.Mul4(translate)

		projectionMatrix := mgl32.Ident4()
		perspective := mgl32.Perspective(mgl32.DegToRad(45.0), width/height, 0.1, 100.0)
		projectionMatrix = projectionMatrix.Mul4(perspective)

		shader.SetMat4("model", modelMatrix)
		shader.SetMat4("view", viewMatrix)
		shader.SetMat4("projection", projectionMatrix)

		gl.BindVertexArray(vao)
		gl.DrawArrays(gl.TRIANGLES, 0, 36) // 36 vertices to make a cube

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
