package model

import "github.com/go-gl/mathgl/mgl32"

const (
	sizefloat32      = 4
	sizeInt          = 4
	maxBoneInfluence = 4
)

type Vertex struct {
	Position  mgl32.Vec3
	Normal    mgl32.Vec3
	TexCoords mgl32.Vec2
	Tangent   mgl32.Vec3
	Bitangent mgl32.Vec3
	BoneIDs   [maxBoneInfluence]int
	Weights   [maxBoneInfluence]float32
}
