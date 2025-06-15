package scenes

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/igoramorim/gopengl/pkg/shader"
)

type Textures struct{}

func (s Textures) Name() string {
	return "textures"
}

func (s Textures) Width() int {
	return width
}

func (s Textures) Height() int {
	return height
}

func (s Textures) Show() {
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

	shader, err := shader.New("internal/assets/shaders/texture.vert", "internal/assets/shaders/texture.frag")
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

	// Load first image
	imageData0, err := s.loadImage("internal/assets/textures/container.jpg")
	if err != nil {
		panic(err)
	}

	// Generate the first texture
	var texture0 uint32
	// Generate one texture
	gl.GenTextures(1, &texture0)
	// Bind the texture BEFORE setting the gl.TEXTURE_2D configurations
	gl.BindTexture(gl.TEXTURE_2D, texture0)
	// Set the texture filtering mode when downscaling
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	// Set the texture filtering mode when upscaling
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	// Set the texture wrapping mode horizontally
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	// Set the texture wrapping mode vertically
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,                   // Speficies the texture target that was bind with gl.BindTexture(gl.TEXTURE_2D) before
		0,                               // Mipmap level if you want to set manually. 0 is the base level
		gl.RGBA,                         // The format we want OpenGL to store our texture
		int32(imageData0.Rect.Size().X), // Width of the texture
		int32(imageData0.Rect.Size().Y), // Height of the texture
		0,                               // Border?
		gl.RGBA,                         // Format of the source image
		gl.UNSIGNED_BYTE,                // Datatype wich the image was loaded in
		gl.Ptr(imageData0.Pix),          // The image data
	)
	// gl.GenerateTextureMipmap(texture) // FIXME: Está gerando panic

	// Load second image
	imageData1, err := s.loadImage("internal/assets/textures/awesomeface.png")
	if err != nil {
		panic(err)
	}

	// Generate the second texture
	var texture1 uint32
	gl.GenTextures(1, &texture1)
	gl.BindTexture(gl.TEXTURE_2D, texture1)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(imageData1.Rect.Size().X),
		int32(imageData1.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(imageData1.Pix),
	)
	// gl.GenerateTextureMipmap(texture) // FIXME: Está gerando panic

	// Clean up all resources
	defer func() {
		gl.DeleteVertexArrays(1, &vao)
		gl.DeleteBuffers(1, &vbo)
		gl.DeleteBuffers(1, &ebo)
		shader.Delete()
	}()

	// Main loop
	for !window.ShouldClose() {
		processInput(window, s)

		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		// Need to activate the shader before setting the texture uniform
		shader.Use()

		// Activate, bind and set the first texture uniform in the shader
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture0)
		shader.SetInt("texture0", 0)

		// Activate, bind and set the second texture uniform in the shader
		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, texture1)
		shader.SetInt("texture1", 1)

		gl.BindVertexArray(vao)
		// Draw the rectangle using the ebo
		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func (d Textures) loadImage(file string) (*image.RGBA, error) {
	imgFile, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("texture %q not found on disk: %v", file, err)
	}
	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return nil, fmt.Errorf("unsupported stride")
	}

	// FIXME: Images are being loaded upside-down
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	return rgba, nil
}
