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

func NewMaterials() Materials {
	return Materials{
		camera:     camera.New(),
		firstMouse: true,
		lastX:      float64(width) / 2,
		lastY:      float64(height) / 2,
		deltaTime:  0.0,
		lastFrame:  0.0,
	}
}

type Materials struct {
	camera     *camera.Camera
	firstMouse bool
	lastX      float64
	lastY      float64
	deltaTime  float64 // Time between current frame and last frame
	lastFrame  float64
}

func (s Materials) Name() string {
	return "materials"
}

func (s Materials) Width() int {
	return width
}

func (s Materials) Height() int {
	return height
}

func (s Materials) Show() {
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

	lightingShader, err := shader.New("internal/assets/shaders/materials.vert", "internal/assets/shaders/materials.frag")
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

	type Material struct {
		ambientColor  mgl32.Vec3
		diffuseColor  mgl32.Vec3
		specularColor mgl32.Vec3
		shininess     float32
	}

	// http://devernay.free.fr/cours/opengl/materials.html
	materials := []Material{
		{mgl32.Vec3{0.0215, 0.1745, 0.0215}, mgl32.Vec3{0.07568, 0.61424, 0.07568}, mgl32.Vec3{0.633, 0.727811, 0.633}, 0.6 * 128.0},                       //emerald
		{mgl32.Vec3{0.135, 0.2225, 0.1575}, mgl32.Vec3{0.54, 0.89, 0.63}, mgl32.Vec3{0.316228, 0.316228, 0.316228}, 0.1 * 128.0},                           //jade
		{mgl32.Vec3{0.05375, 0.05, 0.06625}, mgl32.Vec3{0.18275, 0.17, 0.22525}, mgl32.Vec3{0.332741, 0.328634, 0.346435}, 0.3 * 128.0},                    //obsidian
		{mgl32.Vec3{0.25, 0.20725, 0.20725}, mgl32.Vec3{1, 0.829, 0.829}, mgl32.Vec3{0.296648, 0.296648, 0.296648}, 0.088 * 128.0},                         //pearl
		{mgl32.Vec3{0.1745, 0.01175, 0.01175}, mgl32.Vec3{0.61424, 0.04136, 0.04136}, mgl32.Vec3{0.727811, 0.626959, 0.626959}, 0.6 * 128.0},               //ruby
		{mgl32.Vec3{0.1, 0.18725, 0.1745}, mgl32.Vec3{0.396, 0.74151, 0.69102}, mgl32.Vec3{0.297254, 0.30829, 0.306678}, 0.1 * 128.0},                      //turquoise
		{mgl32.Vec3{0.329412, 0.223529, 0.027451}, mgl32.Vec3{0.780392, 0.568627, 0.113725}, mgl32.Vec3{0.992157, 0.941176, 0.807843}, 0.21794872 * 128.0}, //brass
		{mgl32.Vec3{0.2125, 0.1275, 0.054}, mgl32.Vec3{0.714, 0.4284, 0.18144}, mgl32.Vec3{0.393548, 0.271906, 0.166721}, 0.2 * 128.0},                     //bronze
		{mgl32.Vec3{0.25, 0.25, 0.25}, mgl32.Vec3{0.4, 0.4, 0.4}, mgl32.Vec3{0.774597, 0.774597, 0.774597}, 0.6 * 128.0},                                   //chrome
		{mgl32.Vec3{0.19125, 0.0735, 0.0225}, mgl32.Vec3{0.7038, 0.27048, 0.0828}, mgl32.Vec3{0.256777, 0.137622, 0.086014}, 0.1 * 128.0},                  //copper
		{mgl32.Vec3{0.24725, 0.1995, 0.0745}, mgl32.Vec3{0.75164, 0.60648, 0.22648}, mgl32.Vec3{0.628281, 0.555802, 0.366065}, 0.4 * 128.0},                //gold
		{mgl32.Vec3{0.19225, 0.19225, 0.19225}, mgl32.Vec3{0.50754, 0.50754, 0.50754}, mgl32.Vec3{0.508273, 0.508273, 0.508273}, 0.4 * 128.0},              //silver
		{mgl32.Vec3{0.0, 0.0, 0.0}, mgl32.Vec3{0.01, 0.01, 0.01}, mgl32.Vec3{0.50, 0.50, 0.50}, .25 * 128.0},                                               //black plastic
		{mgl32.Vec3{0.0, 0.1, 0.06}, mgl32.Vec3{0.0, 0.50980392, 0.50980392}, mgl32.Vec3{0.50196078, 0.50196078, 0.50196078}, .25 * 128.0},                 //cyan plastic
		{mgl32.Vec3{0.0, 0.0, 0.0}, mgl32.Vec3{0.1, 0.35, 0.1}, mgl32.Vec3{0.45, 0.55, 0.45}, .25 * 128.0},                                                 //green plastic
		{mgl32.Vec3{0.0, 0.0, 0.0}, mgl32.Vec3{0.5, 0.0, 0.0}, mgl32.Vec3{0.7, 0.6, 0.6}, .25 * 128.0},                                                     //red plastic
		{mgl32.Vec3{0.0, 0.0, 0.0}, mgl32.Vec3{0.55, 0.55, 0.55}, mgl32.Vec3{0.70, 0.70, 0.70}, .25 * 128.0},                                               //white plastic
		{mgl32.Vec3{0.0, 0.0, 0.0}, mgl32.Vec3{0.5, 0.5, 0.0}, mgl32.Vec3{0.60, 0.60, 0.50}, .25 * 128.0},                                                  //yellow plastic
		{mgl32.Vec3{0.02, 0.02, 0.02}, mgl32.Vec3{0.01, 0.01, 0.01}, mgl32.Vec3{0.4, 0.4, 0.4}, .078125 * 128.0},                                           //black rubber
		{mgl32.Vec3{0.0, 0.05, 0.05}, mgl32.Vec3{0.4, 0.5, 0.5}, mgl32.Vec3{0.04, 0.7, 0.7}, .078125 * 128.0},                                              //cyan rubber
		{mgl32.Vec3{0.0, 0.05, 0.0}, mgl32.Vec3{0.4, 0.5, 0.4}, mgl32.Vec3{0.04, 0.7, 0.04}, .078125 * 128.0},                                              //green rubber
		{mgl32.Vec3{0.05, 0.0, 0.0}, mgl32.Vec3{0.5, 0.4, 0.4}, mgl32.Vec3{0.7, 0.04, 0.04}, .078125 * 128.0},                                              //red rubber
		{mgl32.Vec3{0.05, 0.05, 0.05}, mgl32.Vec3{0.5, 0.5, 0.5}, mgl32.Vec3{0.7, 0.7, 0.7}, .078125 * 128.0},                                              //white rubber
		{mgl32.Vec3{0.05, 0.05, 0.0}, mgl32.Vec3{0.5, 0.5, 0.4}, mgl32.Vec3{0.7, 0.7, 0.04}, .078125 * 128.0},                                              //yellow rubber
	}

	// Main loop
	for !window.ShouldClose() {
		s.processInput(window)

		gl.ClearColor(0.1, 0.1, 0.1, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		currentFrame := glfw.GetTime()
		s.deltaTime = currentFrame - s.lastFrame
		s.lastFrame = currentFrame

		lightPos := mgl32.Vec3{
			float32(2.0 + math.Sin(glfw.GetTime()*2.0)*2.0),
			float32(-1.5 + math.Sin(glfw.GetTime())*1.5),
			1.0,
		}
		// fmt.Println(lightPos)

		lightColor := mgl32.Vec3{
			// float32(math.Sin(glfw.GetTime()*2.0)*0.5 + 0.5),
			// float32(math.Sin(glfw.GetTime()*0.7)*0.5 + 0.5),
			// float32(math.Sin(glfw.GetTime()*1.3)*0.5 + 0.5),
			1.0, 1.0, 1.0,
		}
		// diffuseColor := lightColor.Mul(0.5)
		// ambientColor := diffuseColor.Mul(0.2)

		viewMatrix := s.camera.ViewMatrix()
		modelMatrix := mgl32.Ident4()
		projectionMatrix := mgl32.Ident4()
		projectionMatrix = mgl32.Perspective(mgl32.DegToRad(float32(s.camera.Fov)), width/height, 0.1, 100.0)

		for i := 0; i < len(materials); i++ {
			modelMatrix = mgl32.Ident4()

			lightingShader.Use()
			lightingShader.SetVec3("light.position", lightPos)
			lightingShader.SetVec3("viewPos", s.camera.Position)

			lightingShader.SetVec3f("light.ambient", 1.0, 1.0, 1.0)
			lightingShader.SetVec3f("light.diffuse", 1.0, 1.0, 1.0)
			lightingShader.SetVec3f("light.specular", 1.0, 1.0, 1.0)

			lightingShader.SetVec3("material.ambient", materials[i].ambientColor)
			lightingShader.SetVec3("material.diffuse", materials[i].diffuseColor)
			lightingShader.SetVec3("material.specular", materials[i].specularColor)
			lightingShader.SetFloat("material.shininess", materials[i].shininess)

			lightingShader.SetMat4("view", viewMatrix)
			lightingShader.SetMat4("projection", projectionMatrix)

			var x, y float32
			x = float32(i) * 0.5

			size := 5
			if i >= size {
				x = float32(i-size) * 0.5
				y = y - 0.5
			}
			if i >= (size * 2) {
				x = float32(i-size*2) * 0.5
				y = y - 0.5
			}
			if i >= (size * 3) {
				x = float32(i-size*3) * 0.5
				y = y - 0.5
			}
			if i >= (size * 4) {
				x = float32(i-size*4) * 0.5
				y = y - 0.5
			}

			cubePosition := mgl32.Vec3{x, y, 0.0}
			translate := mgl32.Translate3D(cubePosition.X(), cubePosition.Y(), cubePosition.Z())
			modelMatrix = modelMatrix.Mul4(translate)
			lightingShader.SetMat4("model", modelMatrix)

			// Render the cube
			gl.BindVertexArray(cubeVAO)
			gl.DrawArrays(gl.TRIANGLES, 0, 36)
		}

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

func (s *Materials) processInput(w *glfw.Window) {
	processInput(w, s)
	processCameraKeyboardInput(w, s.camera, s.deltaTime)
}

func (s *Materials) mouseCallback(w *glfw.Window, xpos, ypos float64) {
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

func (s *Materials) mouseScrollCallback(w *glfw.Window, xoff, yoff float64) {
	s.camera.ProcessMouseScroll(yoff)
}
