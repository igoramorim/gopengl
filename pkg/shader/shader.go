package shader

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
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

	return &Shader{id: id}, nil
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
	id uint32
}

func (s *Shader) Use() {
	gl.UseProgram(s.id)
}

func (s *Shader) Delete() {
	gl.DeleteProgram(s.id)
}

func (s *Shader) SetInt(name string, value int32) {
	uniform := gl.GetUniformLocation(s.id, gl.Str(name+"\x00"))
	gl.Uniform1i(uniform, value)
}
