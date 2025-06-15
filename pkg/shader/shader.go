package shader

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

func New(vertexPath, fragPath string) (*Shader, error) {
	vertexCode, err := readFile(vertexPath)
	if err != nil {
		return nil, err
	}

	vertexShader, err := buildShader(gl.VERTEX_SHADER, vertexCode)
	if err != nil {
		return nil, err
	}
	defer gl.DeleteShader(vertexShader)

	fragCode, err := readFile(fragPath)
	if err != nil {
		return nil, err
	}

	fragShader, err := buildShader(gl.FRAGMENT_SHADER, fragCode)
	if err != nil {
		return nil, err
	}
	defer gl.DeleteShader(fragShader)

	id := gl.CreateProgram()
	gl.AttachShader(id, vertexShader)
	gl.AttachShader(id, fragShader)
	gl.LinkProgram(id)

	if err := checkCompileErr(id, "PROGRAM"); err != nil {
		return nil, err
	}

	return &Shader{ID: id}, nil
}

func readFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func buildShader(xtype uint32, sourceCode []byte) (uint32, error) {
	shader := gl.CreateShader(xtype)

	csrc, free := gl.Strs(string(sourceCode) + "\x00")
	defer free()

	gl.ShaderSource(shader, 1, csrc, nil)
	gl.CompileShader(shader)

	var typestr string
	switch xtype {
	case gl.VERTEX_SHADER:
		typestr = "VERTEX"
	case gl.FRAGMENT_SHADER:
		typestr = "FRAGMENT"
	default:
		typestr = "UNKNOWN"
	}

	if err := checkCompileErr(shader, typestr); err != nil {
		return 0, err
	}

	return shader, nil
}

func checkCompileErr(shader uint32, xtype string) error {
	var status int32

	if xtype != "PROGRAM" {
		gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	} else {
		gl.GetProgramiv(shader, gl.LINK_STATUS, &status)
	}

	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		msg := fmt.Sprintf("compile %s shader\n %s\n", xtype, log)
		return errors.New(msg)
	}

	return nil
}

type Shader struct {
	ID uint32
}

func (s *Shader) Use() {
	gl.UseProgram(s.ID)
}

func (s *Shader) Delete() {
	gl.DeleteProgram(s.ID)
}

// TODO: Add uniform location cache to be used in every Set*Uniform* method below

func (s *Shader) SetInt(name string, value int32) {
	uniform := gl.GetUniformLocation(s.ID, gl.Str(name+"\x00"))
	gl.Uniform1i(uniform, value)
}

func (s *Shader) SetFloat(name string, value float32) {
	uniform := gl.GetUniformLocation(s.ID, gl.Str(name+"\x00"))
	gl.Uniform1f(uniform, value)
}

func (s *Shader) SetMat4(name string, value mgl32.Mat4) {
	uniform := gl.GetUniformLocation(s.ID, gl.Str(name+"\x00"))
	gl.UniformMatrix4fv(uniform, 1, false, &value[0])
}

func (s *Shader) SetVec3f(name string, x, y, z float32) {
	uniform := gl.GetUniformLocation(s.ID, gl.Str(name+"\x00"))
	gl.Uniform3f(uniform, x, y, z)
}

func (s *Shader) SetVec3(name string, value mgl32.Vec3) {
	uniform := gl.GetUniformLocation(s.ID, gl.Str(name+"\x00"))
	gl.Uniform3fv(uniform, 1, &value[0])
}
