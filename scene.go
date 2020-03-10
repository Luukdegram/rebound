package rebound

// SceneComponentName is the name of the SceneComponent
const SceneComponentName = "SceneComponent"

// SceneComponent defines a 3D scene. But is not limited to 1 scene in 1 rendering cycle.
type SceneComponent struct {
	Nodes []*Node
}

// Name returns the name of the SceneComponent
func (sc *SceneComponent) Name() string {
	return SceneComponentName
}

// Node holds information regarding a singular node within a scene.
// It can have a Name, transformation data or children.
type Node struct {
	Name           string
	Transformation [16]float32
	Children       []*Node
	Mesh           *Mesh
	Camera         *Camera
}

// Mesh holds geometry data.
type Mesh struct {
	ID         uint32
	Attributes []Attribute
	Indices    []uint32
	Material   *Material
}

// VertexCount returns the amount of vertices the Mesh contains.
func (m *Mesh) VertexCount() int {
	return len(m.Indices)
}

// Material describes the look of a geometric object
type Material struct {
	Transparent      bool
	NormalTexture    *TextureComponent
	OcclusionTexture *TextureComponent
	EmmisiveTexture  *TextureComponent
	PBRMetallicRoughness
}

// PBRMetallicRoughness holds all data related to PBR such as roughness, basecolor and metallicness
type PBRMetallicRoughness struct {
	BaseColor                [4]float32
	BaseColorTexture         *TextureComponent
	MetallicFactor           int
	RoughnessFactor          int
	MetallicRoughnessTexture *TextureComponent
}
