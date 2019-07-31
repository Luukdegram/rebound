package models

//RawModel is a basic model without any textures
type RawModel struct {
	VaoID       uint32
	VertexCount int
}

//Mesh is the base type that contains all data that belongs to a 3D model
type Mesh struct {
	RawModel RawModel
	Textures *ModelTexture
	Childs   []Mesh
}

//IsTextured returns true if the Mesh has a modeltexture
func (m Mesh) IsTextured() (out bool) {
	out = m.Textures != nil
}
