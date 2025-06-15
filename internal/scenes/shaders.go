package scenes

import (
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type Shaders struct{}

func (s Shaders) Name() string {
	return "shaders_uniforms"
}

func (s Shaders) Width() int {
	return width
}

func (s Shaders) Height() int {
	return height
}

func (s Shaders) Show() {
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

	vertexShader := gl.CreateShader(gl.VERTEX_SHADER)
	vertexShaderCSource, free := gl.Strs(s.vertexShaderSource())
	gl.ShaderSource(vertexShader, 1, vertexShaderCSource, nil)
	free()
	gl.CompileShader(vertexShader)

	var status int32
	gl.GetShaderiv(vertexShader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(vertexShader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(vertexShader, logLength, nil, gl.Str(log))

		panic(fmt.Sprintf("compile shader source %s\n %s\n", s.vertexShaderSource(), log))
	}

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

	shaderProgram := gl.CreateProgram()
	gl.AttachShader(shaderProgram, vertexShader)
	gl.AttachShader(shaderProgram, fragmentShader)
	gl.LinkProgram(shaderProgram)

	gl.GetProgramiv(shaderProgram, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(shaderProgram, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(shaderProgram, logLength, nil, gl.Str(log))

		panic(fmt.Sprintf("linking shader program %v\n", log))
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	var vertices = []float32{
		// x y z
		-0.5, -0.5, 0.0, // left
		0.5, -0.5, 0.0, // right
		0.0, 0.5, 0.0, // top
	}

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*floatSize, gl.Ptr(vertices), gl.STATIC_DRAW)

	// Position attribute
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*floatSize, nil)
	gl.EnableVertexAttribArray(0)

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

		gl.UseProgram(shaderProgram)

		t := glfw.GetTime()
		green := math.Sin(t)/2.0 + 0.5

		// Send the value to the shader uniform variable named 'color'
		uniformLocation := gl.GetUniformLocation(shaderProgram, gl.Str("color\x00"))
		gl.Uniform3f(uniformLocation, 0.0, float32(green), 0.0)

		gl.BindVertexArray(vao)
		gl.DrawArrays(gl.TRIANGLES, 0, 3)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func (d Shaders) vertexShaderSource() string {
	return `
#version 330 core

layout (location = 0) in vec3 position;

void main() {
	gl_Position = vec4(position, 1.0f);
}
` + "\x00"
}

func (d Shaders) fragmentShaderSource() string {
	return `
#version 330 core

uniform vec3 color; // value is set in the OpenGL code

out vec4 FragColor;

void main() {
	FragColor = vec4(color, 1.0f);
}
` + "\x00"
}
