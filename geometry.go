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

//Attribute is vbo that stores data such as texture coordinates
type Attribute struct {
	Type AttributeType
	Data []float32
	Size int
}

//AttributeType can be used to link external attribute names to Rebound's
type AttributeType int

const (
	//Position is a shader attribute used for positional coordinates
	Position AttributeType = iota
	//TexCoords is a shader attribute used for the coordinates of a texture
	TexCoords
	//Normals holds the coordinates of a normal texture
	Normals
	//Tangents holds the tangets data in a shader
	Tangents
)

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
