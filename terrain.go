package rebound

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/luukdegram/rebound/models"
)

const (
	terrainName = "Terrain"
)

//terrain is a struct that contains all data needed to generate a 3D terrain for our world
type terrain struct {
	size        float32
	vertexCount int
	x           float32
	z           float32
}

//NewTerrain creates a new terrain with the given x and z axis
func NewTerrain(x, z int, texture *models.ModelTexture) *Entity {
	t := new(terrain)
	t.size = 800
	t.vertexCount = 128
	t.x = float32(x) * t.size
	t.z = t.size * float32(z)

	e := NewEntity()
	e.Position = mgl32.Vec3{t.x, 0, t.z}
	e.Geometry = generateTerrain(*t, texture)
	return e
}

func generateTerrain(t terrain, texture *models.ModelTexture) *Geometry {
	count := t.vertexCount * t.vertexCount
	vertices := make([]float32, count*3)
	normals := make([]float32, count*3)
	texCoords := make([]float32, count*2)
	indices := make([]uint32, 6*(t.vertexCount-1)*(t.vertexCount-1))
	vertexPointer := 0
	for i := 0; i < t.vertexCount; i++ {
		for j := 0; j < t.vertexCount; j++ {
			vertices[vertexPointer*3] = float32(j) / float32(t.vertexCount-1) * t.size
			vertices[vertexPointer*3+1] = 0
			vertices[vertexPointer*3+2] = float32(i) / float32(t.vertexCount-1) * t.size
			normals[vertexPointer*3] = 0
			normals[vertexPointer*3+1] = 1
			normals[vertexPointer*3+2] = 0
			texCoords[vertexPointer*2] = float32(j) / float32(t.vertexCount-1)
			texCoords[vertexPointer*2+1] = float32(i) / float32(t.vertexCount-1)
			vertexPointer++
		}
	}

	pointer := 0
	for gz := 0; gz < t.vertexCount-1; gz++ {
		for gx := 0; gx < t.vertexCount-1; gx++ {
			topLeft := gz*t.vertexCount + gx
			topRight := topLeft + 1
			bottomLeft := (gz+1)*t.vertexCount + gx
			bottomRight := bottomLeft + 1
			indices[pointer] = uint32(topLeft)
			pointer++
			indices[pointer] = uint32(bottomLeft)
			pointer++
			indices[pointer] = uint32(topRight)
			pointer++
			indices[pointer] = uint32(topRight)
			pointer++
			indices[pointer] = uint32(bottomLeft)
			pointer++
			indices[pointer] = uint32(bottomRight)
			pointer++
		}
	}

	attributes := []Attribute{
		Attribute{Type: Pos, Data: vertices, Size: 3},
		Attribute{Type: Normals, Data: normals, Size: 3},
		Attribute{Type: TexCoords, Data: vertices, Size: 2},
	}

	//model := LoadToVAO(indices, attributes)
	mesh := NewMesh(terrainName)
	//mesh.RawModel = model
	mesh.Texture = texture
	for _, attribute := range attributes {
		mesh.AddAttribute(attribute)
	}
	return NewGeometry(*mesh)
}
