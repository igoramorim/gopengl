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

type Transformations struct{}

func (s Transformations) Name() string {
	return "transformations"
}

func (s Transformations) Width() int {
	return width
}

func (s Transformations) Height() int {
	return height
}

func (s Transformations) Show() {
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

	shader, err := shader.New("internal/assets/shaders/transformation.vert", "internal/assets/shaders/transformation.frag")
	if err != nil {
		panic(err)
	}

	var vertices = []float32{
		// x y z u v (tex coord)
		0.5, 0.5, 0.0, 1.0, 1.0, // top right
		0.5, -0.5, 0.0, 1.0, 0.0, // bottom right
		-0.5, -0.5, 0.0, 0.0, 0.0, // bottom left
		-0.5, 0.5, 0.0, 0.0, 1.0, // top left
	}

	// We need two triangles to draw a rectangle. To avoid duplicate vertices we can use an array of
	// with the indices requried to draw both triangles. Note that the indices 1 and 3 appears on both triangles.
	var indices = []uint32{
		0, 1, 3, // first triangle
		1, 2, 3, // second triangle
	}

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*floatSize, gl.Ptr(vertices), gl.STATIC_DRAW)

	// Element Buffer Object
	var ebo uint32
	gl.GenBuffers(1, &ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*uint32Size, gl.Ptr(indices), gl.STATIC_DRAW)

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
		gl.DeleteBuffers(1, &ebo)
		shader.Delete()
		texture0.Delete()
		texture1.Delete()
	}()

	// Main loop
	for !window.ShouldClose() {
		processInput(window, s)

		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		shader.Use()

		texture0.ActiveAndBind()
		shader.SetInt("texture0", 0)

		texture1.ActiveAndBind()
		shader.SetInt("texture1", 1)

		time := glfw.GetTime()

		// Transformations
		// The matrix multiplication is applied in reverse (from bottom to top). So the order is:
		// 1. Scale
		// 2. Rotate
		// 3. Translate
		// Try switching the order between translate and rotate.
		// When the translate is applied first this is what happens:
		// Its rotation origin is no longer (0,0,0) making it looks as if its circling around the origin of the scene

		// Identity matrix:
		// 1 0 0 0
		// 0 1 0 0
		// 0 0 1 0
		// 0 0 0 1
		transform := mgl32.Ident4()

		// Transformations
		translate := mgl32.Translate3D(0.5, -0.5, 0.0)
		rotate := mgl32.HomogRotate3D(float32(time), mgl32.Vec3{0.0, 0.0, 1.0})
		scale := mgl32.Scale3D(0.5, 0.5, 0.5)

		// Multiplications
		transform = transform.Mul4(translate)
		transform = transform.Mul4(rotate)
		transform = transform.Mul4(scale)

		shader.SetMat4("transform", transform)

		gl.BindVertexArray(vao)
		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)

		// TODO: Control how many seconds will be taking screen shots
		// TODO: Be able to control how many screen shots will be taken in a second
		// screenShot(s.Name())

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
