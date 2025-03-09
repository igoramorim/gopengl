package scenes

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type Triangle struct{}

func (s Triangle) Name() string {
	return "triangle"
}

func (s Triangle) Width() int {
	return width
}

func (s Triangle) Height() int {
	return height
}

func (s Triangle) Show() {
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

	// Handles window resize. Calls frameBufferSizeCallback func whenever the window changes in size
	window.SetFramebufferSizeCallback(frameBufferSizeCallback)

	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version:", version)

	// Vertex shader
	vertexShader := gl.CreateShader(gl.VERTEX_SHADER)
	vertexShaderCSource, free := gl.Strs(s.vertexShaderSource())
	// Attach the shader source code to the shader object
	gl.ShaderSource(vertexShader, 1, vertexShaderCSource, nil)
	free()
	gl.CompileShader(vertexShader)

	// Checks if the shader compiled ok
	var status int32
	gl.GetShaderiv(vertexShader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(vertexShader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(vertexShader, logLength, nil, gl.Str(log))

		panic(fmt.Sprintf("compile shader source %s\n %s\n", s.vertexShaderSource(), log))
	}

	// Fragment shader
	fragmentShader := gl.CreateShader(gl.FRAGMENT_SHADER)
	fragmentShaderCSource, free := gl.Strs(s.fragmentShaderSource())
	gl.ShaderSource(fragmentShader, 1, fragmentShaderCSource, nil)
	free()
	gl.CompileShader(fragmentShader)

	gl.GetShaderiv(fragmentShader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(fragmentShader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(fragmentShader, logLength, nil, gl.Str(log))

		panic(fmt.Sprintf("compile shader source %s\n %s\n", s.fragmentShaderSource(), log))
	}

	// Shader program. Link vertex and fragment shaders into one obeject
	shaderProgram := gl.CreateProgram()
	gl.AttachShader(shaderProgram, vertexShader)
	gl.AttachShader(shaderProgram, fragmentShader)
	gl.LinkProgram(shaderProgram)

	// Checks if linking failed
	gl.GetProgramiv(shaderProgram, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(shaderProgram, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(shaderProgram, logLength, nil, gl.Str(log))

		panic(fmt.Sprintf("linking shader program %v\n", log))
	}

	// Once the linking is done, we do not need the shader objects anymore
	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	// Vertex input data
	var vertices = []float32{
		// x y z r g b
		-0.5, -0.5, 0.0, 1.0, 0.0, 0.0, // left red
		0.5, -0.5, 0.0, 0.0, 1.0, 0.0, // right green
		0.0, 0.5, 0.0, 0.0, 0.0, 1.0, // top blue
	}

	// Vertex Array Object. Used to make it easy to switch between vertex buffers / attributes
	var vao uint32
	// Generate a vertex array ID
	gl.GenVertexArrays(1, &vao)
	// Bind the vertex array before the vertex buffer(s)
	gl.BindVertexArray(vao)

	// Vertex Buffer Object. Used to store the vertices in the GPU's memory
	var vbo uint32
	// Generate a vertex buffer ID
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	// Copy user data (vertices) into the currently bound buffer (VBO wich was binded to GL_ARRAY_BUFFER)
	// Now we have vertex data stored in the GPU memory managed by a vertex buffer object (VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*floatSize, gl.Ptr(vertices), gl.STATIC_DRAW)

	// We need to tell OpenGL how to read the vertex input data. Attribute position
	gl.VertexAttribPointer(
		0,           // Wich vertex we want to configure. Vertex shader location = 0
		3,           // Size of the vertex attribute. The input is a vec3, so it is composed of 3 values
		gl.FLOAT,    // Type of the data. vec3 is a vector of 3 floats
		false,       // Data should be normalized?
		6*floatSize, // Stride: Space between consecutive vertex attribute. Our input data have 6 float values, wich will be 3 packs of floats
		nil,         // Offset where the position data begings in the buffer. The position is the first 3 values, so no need a offset
	)
	// Enables attribute ID 0 wich was just set
	gl.EnableVertexAttribArray(0)

	// Tells OpenGL how to read the Color attribute
	// Note that now we do have a offset. The color data is the second group of floats
	gl.VertexAttribPointerWithOffset(1, 3, gl.FLOAT, false, 6*floatSize, 3*floatSize)
	gl.EnableVertexAttribArray(1)

	// Clean up all resources
	defer func() {
		gl.DeleteVertexArrays(1, &vao)
		gl.DeleteBuffers(1, &vbo)
		gl.DeleteProgram(shaderProgram)
	}()

	// Main loop
	for !window.ShouldClose() {
		processInput(window, s)

		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		// Sets what shader program the render calls will use
		gl.UseProgram(shaderProgram)
		// Draws the triangle using the data from the vao
		gl.BindVertexArray(vao)
		gl.DrawArrays(
			gl.TRIANGLES, // Mode we want to draw
			0,            // Start index of the vertex array we want to draw
			3,            // How many vertices we want to draw
		)

		// Swap the front buffer and the back buffer
		window.SwapBuffers()
		// Checks if any events are triggered (like keyboard)
		glfw.PollEvents()
	}
}

func (s Triangle) vertexShaderSource() string {
	return `
#version 330 core

layout (location = 0) in vec3 position;
layout (location = 1) in vec3 color;

out vec3 vertexColor; // specify a color ouptput to the fragment shader

void main() {
	gl_Position = vec4(position, 1.0f);
	vertexColor = color;
}
` + "\x00"
}

func (s Triangle) fragmentShaderSource() string {
	return `
#version 330 core

in vec3 vertexColor; // the input variable from the vertex shader (must be same name and type)
out vec4 FragColor;

void main() {
	FragColor = vec4(vertexColor, 1.0f);
}
` + "\x00"
}
