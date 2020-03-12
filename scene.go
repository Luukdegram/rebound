package rebound

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
	NormalTexture    *uint32
	OcclusionTexture *uint32
	EmmisiveTexture  *uint32
	PBRMetallicRoughness
}

// PBRMetallicRoughness holds all data related to PBR such as roughness, basecolor and metallicness
type PBRMetallicRoughness struct {
	BaseColor                [4]float32
	BaseColorTexture         *uint32
	MetallicFactor           int
	RoughnessFactor          int
	MetallicRoughnessTexture *uint32
}
