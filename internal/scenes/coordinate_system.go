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

type CoordinateSystem struct{}

func (s CoordinateSystem) Name() string {
	return "coordinate_system"
}

func (s CoordinateSystem) Width() int {
	return width
}

func (s CoordinateSystem) Height() int {
	return height
}

func (s CoordinateSystem) Show() {
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
		0.5, 0.5, 0.0, 1.0, 1.0, // top right
		0.5, -0.5, 0.0, 1.0, 0.0, // bottom right
		-0.5, -0.5, 0.0, 0.0, 0.0, // bottom left
		-0.5, 0.5, 0.0, 0.0, 1.0, // top left
	}

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

	shader.Use()

	texture0.ActiveAndBind()
	shader.SetInt("texture0", 0)

	texture1.ActiveAndBind()
	shader.SetInt("texture1", 1)

	// Transformations to make it 3D

	// Model matrix. Applies transformations to the object's vertices
	modelMatrix := mgl32.Ident4()
	rotate := mgl32.HomogRotate3D(mgl32.DegToRad(-55.0), mgl32.Vec3{1.0, 0.0, 0.0})
	modelMatrix = modelMatrix.Mul4(rotate)

	// View matrix. Applies transformations to the 'eye' but actually we move the entire scene
	// so if you want to move the backwards, you move the entire scene forward
	viewMatrix := mgl32.Ident4()
	translate := mgl32.Translate3D(0.0, 0.0, -2.0)
	viewMatrix = viewMatrix.Mul4(translate)

	// Projection matrix
	projectionMatrix := mgl32.Ident4()
	perspective := mgl32.Perspective(
		mgl32.DegToRad(45.0),           // Field of view (FOV)
		float32(width)/float32(height), // Aspect ratio
		0.1,                            // Near plane. Vertices between the near plane and the 'eye' won't be rendered
		100.0,                          // Far plane. Vertices after it won't be rendered
	)
	projectionMatrix = projectionMatrix.Mul4(perspective)

	// Main loop
	for !window.ShouldClose() {
		processInput(window, s)

		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		shader.SetMat4("model", modelMatrix)
		shader.SetMat4("view", viewMatrix)
		shader.SetMat4("projection", projectionMatrix)

		gl.BindVertexArray(vao)
		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
