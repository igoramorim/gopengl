package model

import (
	"fmt"
	"image"
	"image/draw"
	"os"
	"strings"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/igoramorim/gopengl/pkg/shader"

	"github.com/bloeys/assimp-go/asig"
	"github.com/go-gl/gl/v4.1-core/gl"
)

func New(path string) (*Model, error) {
	model := &Model{}

	err := model.load(path)
	if err != nil {
		return nil, err
	}

	return model, nil
}

type Model struct {
	texturesLoaded  []Texture
	meshes          []Mesh
	directory       string
	gammaCorrection bool
}

func (m *Model) load(path string) error {
	fsys := os.DirFS(".")
	scene, release, err := asig.ImportFileEx(path, asig.PostProcessTriangulate|asig.PostProcessGenSmoothNormals|asig.PostProcessCalcTangentSpace, fsys)
	if err != nil {
		return err
	}
	defer release()

	// fmt.Printf("path received:%s\n\n", path)
	// fmt.Printf("root node:\n%+v\n\n", scene.RootNode)

	m.directory = path[:strings.LastIndex(path, "/")]

	if err := m.processNode(scene.RootNode, scene); err != nil {
		return err
	}

	return nil
}

func (m *Model) processNode(aiNode *asig.Node, aiScene *asig.Scene) error {
	for _, i := range aiNode.MeshIndicies {
		aiMesh := aiScene.Meshes[i]

		mesh, err := m.processMesh(aiMesh, aiScene)
		if err != nil {
			return err
		}
		m.meshes = append(m.meshes, mesh)
		// fmt.Printf("mesh: %+v\n", mesh)
	}

	for i := range aiNode.Children {
		if err := m.processNode(aiNode.Children[i], aiScene); err != nil {
			return err
		}
	}

	return nil
}

func (m *Model) processMesh(aiMesh *asig.Mesh, aiScene *asig.Scene) (Mesh, error) {
	var vertices []Vertex
	var indices []uint32
	var textures []Texture

	// Walk through each of the mesh's vertices
	for i := range aiMesh.Vertices {
		var vertex Vertex

		// Position
		vec3 := mgl32.Vec3{
			aiMesh.Vertices[i].X(),
			aiMesh.Vertices[i].Y(),
			aiMesh.Vertices[i].Z(),
		}
		vertex.Position = vec3

		// Normals
		if len(aiMesh.Normals) > 0 {
			vec3 = mgl32.Vec3{
				aiMesh.Normals[i].X(),
				aiMesh.Normals[i].Y(),
				aiMesh.Normals[i].Z(),
			}
			vertex.Normal = vec3
		}

		// Texture Coordinates
		if len(aiMesh.TexCoords[0]) > 0 {
			// A vertex can contain up to 8 different texture coordinates. We thus make the assumption
			// that we won't use models where a vertex can have multiple texture coordinates
			// so we always take the first set (0)
			vec2 := mgl32.Vec2{
				aiMesh.TexCoords[0][i].X(),
				aiMesh.TexCoords[0][i].Y(),
			}
			vertex.TexCoords = vec2

			// Tangent
			vec3 = mgl32.Vec3{
				aiMesh.Tangents[i].X(),
				aiMesh.Tangents[i].Y(),
				aiMesh.Tangents[i].Z(),
			}
			vertex.Tangent = vec3

			// Bitangent
			vec3 = mgl32.Vec3{
				aiMesh.BitTangents[i].X(),
				aiMesh.BitTangents[i].Y(),
				aiMesh.BitTangents[i].Z(),
			}
			vertex.Bitangent = vec3
		} else {
			vertex.TexCoords = mgl32.Vec2{0.0, 0.0}
		}
		vertices = append(vertices, vertex)
	}
	// fmt.Printf("vertices: %+v\n\n", vertices)

	// Walk through each of the mesh's faces (a face is a mesh's triangle) and retrieve the
	// corresponding vertex indices
	for i := range aiMesh.Faces {
		aiFace := aiMesh.Faces[i]
		for j := range aiFace.Indices {
			indices = append(indices, uint32(aiFace.Indices[j]))
		}
	}
	// fmt.Printf("indices: %+v\n\n", indices)

	// Materials
	aiMaterial := aiScene.Materials[aiMesh.MaterialIndex]
	// There is a convention assumed for sampler names in the shaders.
	// Each diffuse texture should be named as 'texture_diffuseN' where N is a sequential number
	// ranging from 1 to MAX_SAMPLER_NUMBER.
	// Same applies to other textures as the following list:
	// diffuse:  texture_diffuseN
	// specular: texture_specularN
	// normal:   texture_normalN

	// Diffuse
	diffuseMaps, err := m.loadMaterialTextures(aiMaterial, asig.TextureTypeDiffuse, texDiffuse)
	if err != nil {
		return Mesh{}, err
	}
	textures = append(textures, diffuseMaps...)
	// fmt.Printf("diffuse textures: %+v\n\n", len(diffuseMaps))

	// Specular
	specularMaps, err := m.loadMaterialTextures(aiMaterial, asig.TextureTypeSpecular, texSpecular)
	if err != nil {
		return Mesh{}, err
	}
	textures = append(textures, specularMaps...)
	// fmt.Printf("specular textures: %+v\n\n", len(specularMaps))

	// Normals
	normalMaps, err := m.loadMaterialTextures(aiMaterial, asig.TextureTypeNormal, texNormal)
	if err != nil {
		return Mesh{}, err
	}
	textures = append(textures, normalMaps...)
	// fmt.Printf("normals textures: %+v\n\n", len(normalMaps))

	// Height
	heightMaps, err := m.loadMaterialTextures(aiMaterial, asig.TextureTypeHeight, texHeight)
	if err != nil {
		return Mesh{}, err
	}
	textures = append(textures, heightMaps...)
	// fmt.Printf("height textures: %+v\n\n", len(heightMaps))

	return NewMesh(vertices, indices, textures), nil
}

func (m *Model) loadMaterialTextures(aiMaterial *asig.Material, aiTexType asig.TextureType,
	texType texType) ([]Texture, error) {

	var textures []Texture

	for i := 0; i < asig.GetMaterialTextureCount(aiMaterial, aiTexType); i++ {
		var skip bool
		matInfo, err := asig.GetMaterialTexture(aiMaterial, aiTexType, uint(i))
		if err != nil {
			return nil, err
		}
		// fmt.Printf("material texture: %+v\n\n", matInfo.Path)

		for j := 0; j < len(m.texturesLoaded); j++ {
			// fmt.Printf("texture already loaded: %+v == material texture: %+v\n\n",
			// m.texturesLoaded[j].path, matInfo.Path)
			if m.texturesLoaded[j].path == matInfo.Path {
				// Texture with the same filepath already been loaded. Continue to next one
				textures = append(textures, m.texturesLoaded[j])
				skip = true
				// fmt.Printf("material texture already loaded: %+v\n\n", matInfo.Path)
				break
			}
		}

		// If texture hasn't been loaded already, load it now
		if !skip {
			id, err := textureFromFile(matInfo.Path, m.directory)
			if err != nil {
				return nil, err
			}

			texture := Texture{
				id:    id,
				xtype: texType,
				path:  matInfo.Path,
			}
			textures = append(textures, texture)
			// Store it as texture loaded for entire model to ensure
			// we won't unnecessary load duplicated textures
			m.texturesLoaded = append(m.texturesLoaded, texture)
			// fmt.Printf("texture: %+v\n\n", texture)
		}
	}

	return textures, nil
}

func textureFromFile(path, directory string) (uint32, error) {
	fullpath := directory + "/" + path
	// fmt.Printf("load texture from file: %+v\n\n", fullpath)

	imageData, err := loadImage(fullpath)
	if err != nil {
		return 0, err
	}
	// fmt.Printf("texture data: %+v bytes\n\n", len(imageData.Pix))

	var id uint32
	gl.GenTextures(1, &id)
	gl.BindTexture(gl.TEXTURE_2D, id)

	// fmt.Printf("width: %d height: %d\n", imageData.Rect.Size().X, int32(imageData.Rect.Size().Y))

	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(imageData.Rect.Size().X),
		int32(imageData.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(imageData.Pix),
	)
	// gl.GenerateTextureMipmap(id) // FIXME: It is panicing

	// gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)

	return id, nil
}

func loadImage(path string) (*image.RGBA, error) {
	imgFile, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("texture %q not found on disk: %v", path, err)
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

func (m *Model) Draw(shader *shader.Shader) {
	for i := range m.meshes {
		m.meshes[i].Draw(shader)
	}
}
