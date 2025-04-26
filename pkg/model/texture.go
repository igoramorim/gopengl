package model

type texType string

const (
	texDiffuse  = "texture_diffuse"
	texSpecular = "texture_specular"
	texNormal   = "texture_normal"
	texHeight   = "texture_height"
)

type Texture struct {
	id    uint32
	xtype texType
	path  string
}
