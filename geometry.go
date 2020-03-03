package rebound

import "github.com/luukdegram/rebound/models"

//Mesh is a struct that holds all data of a part of a 3D model
type Mesh struct {
	RawModel   *models.RawModel
	Texture    *models.ModelTexture
	Childs     []Mesh
	attributes []Attribute
	Name       string
}

//Geometry is a struct that holds all meshes of a model
type Geometry struct {
	Meshes []Mesh
}

//IsTextured returns true if the Mesh has a modeltexture
func (m Mesh) IsTextured() (out bool) {
	out = m.Texture != nil
	return
}

//AddAttribute adds one attribute to the mesh
func (m *Mesh) AddAttribute(attribute Attribute) {
	m.attributes = append(m.attributes, attribute)
}

//Attributes returns the attributes of a mesh
func (m Mesh) Attributes() []Attribute {
	return m.attributes
}

//NewGeometry creates a new geometry object with optinal meshesh
func NewGeometry(mesh ...Mesh) *Geometry {
	g := new(Geometry)
	g.Meshes = append(g.Meshes, mesh...)
	return g
}

//NewMesh creates a new Mesh object
func NewMesh(name string) *Mesh {
	return &Mesh{Name: name}
}
