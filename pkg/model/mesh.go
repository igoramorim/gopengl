package model

import (
	"fmt"
	"strconv"
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/igoramorim/gopengl/pkg/shader"
)

func NewMesh(vertices []Vertex, indices []uint32, textures []Texture) Mesh {
	mesh := Mesh{
		Vertices: vertices,
		Indices:  indices,
		Textures: textures,
	}

	mesh.setup()

	return mesh
}

type Mesh struct {
	Vertices []Vertex
	Indices  []uint32
	Textures []Texture
	vao      uint32
	vbo      uint32
	ebo      uint32
}

func (m *Mesh) Draw(shader *shader.Shader) {
	diffuseNr := 1
	specularNr := 1
	normalNr := 1
	heightNr := 1

	for i, tex := range m.Textures {
		var number string
		name := tex.xtype

		switch name {
		case texDiffuse:
			number = strconv.Itoa(diffuseNr)
			diffuseNr++
		case texSpecular:
			number = strconv.Itoa(specularNr)
			specularNr++
		case texNormal:
			number = strconv.Itoa(normalNr)
			normalNr++
		case texHeight:
			number = strconv.Itoa(heightNr)
			heightNr++
		}

		gl.ActiveTexture(gl.TEXTURE0 + uint32(i))

		uniform := fmt.Sprintf("%s%s", name, number)
		// fmt.Printf("uniform name: %+v\n", uniform)
		shader.SetInt(uniform, int32(i))
		gl.BindTexture(gl.TEXTURE_2D, tex.id)
	}

	gl.BindVertexArray(m.vao)
	gl.DrawElements(gl.TRIANGLES, int32(len(m.Indices)), gl.UNSIGNED_INT, nil)
	gl.BindVertexArray(0)

	gl.ActiveTexture(gl.TEXTURE0)
}

func (m *Mesh) setup() {
	var dummy Vertex
	vertexSize := int(unsafe.Sizeof(dummy))
	vertexSize32 := int32(vertexSize)

	gl.GenVertexArrays(1, &m.vao)
	gl.GenBuffers(1, &m.vbo)
	gl.GenBuffers(1, &m.ebo)

	gl.BindVertexArray(m.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, m.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(m.Vertices)*vertexSize, gl.Ptr(m.Vertices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, m.ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(m.Indices)*int(unsafe.Sizeof(m.ebo)), gl.Ptr(m.Indices), gl.STATIC_DRAW)

	// Vertex Positions
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, vertexSize32, nil)

	// Vertex Normals
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointerWithOffset(1, 3, gl.FLOAT, false, vertexSize32, unsafe.Offsetof(dummy.Normal))

	// Vertex Texture Coords
	gl.EnableVertexAttribArray(2)
	gl.VertexAttribPointerWithOffset(2, 2, gl.FLOAT, false, vertexSize32, unsafe.Offsetof(dummy.TexCoords))

	// Vertex Tangent
	gl.EnableVertexAttribArray(3)
	gl.VertexAttribPointerWithOffset(3, 3, gl.FLOAT, false, vertexSize32, unsafe.Offsetof(dummy.Tangent))

	// Vertex Bitangent
	gl.EnableVertexAttribArray(4)
	gl.VertexAttribPointerWithOffset(4, 3, gl.FLOAT, false, vertexSize32, unsafe.Offsetof(dummy.Bitangent))

	// Vertex BoneIDs
	gl.EnableVertexAttribArray(5)
	gl.VertexAttribPointerWithOffset(5, 4, gl.INT, false, vertexSize32, unsafe.Offsetof(dummy.BoneIDs))

	// Vertex Weights
	gl.EnableVertexAttribArray(6)
	gl.VertexAttribPointerWithOffset(6, 4, gl.FLOAT, false, vertexSize32, unsafe.Offsetof(dummy.Weights))

	gl.BindVertexArray(0)
}
