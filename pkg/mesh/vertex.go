package mesh

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

func sizeofVertex() int {
	return (3 * sizefloat32) + // position
		(3 * sizefloat32) + // normal
		(2 * sizefloat32) + // tex coord
		(3 * sizefloat32) + // tangent
		(3 * sizefloat32) + // bittangent
		(maxBoneInfluence * sizeInt) + // bone ids
		(maxBoneInfluence * sizefloat32) // weights
}

func offsetNormals() int {
	return 3 * sizefloat32
}

func offsetTextCoords() int {
	return 6 * sizefloat32
}

func offsetTangent() int {
	return 9 * sizefloat32
}

func offsetBitangent() int {
	return 12 * sizefloat32
}

func offsetBoneIDs() int {
	return 15 * sizeInt
}

func offsetWeights() int {
	return 19 * sizefloat32
}
