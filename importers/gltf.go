package importers

import (
	"path"
	"unsafe"

	"github.com/luukdegram/rebound"
	"github.com/luukdegram/rebound/ecs"

	"github.com/qmuntal/gltf"
)

// GLTFImporter loads a gltf file into Rebound to create a visual object.
type GLTFImporter struct {
	dir string
	doc *gltf.Document
}

// Import loads a GLTF file into a Scene.
func (l *GLTFImporter) Import(file string) (*ecs.Entity, error) {
	doc, err := gltf.Open(file)
	if err != nil {
		return nil, err
	}
	l.doc = doc

	dir := path.Dir(file)

	scene := ecs.NewEntity()
	l.dir = dir

	var index int
	if doc.Scene == nil {
		index = 0
	} else {
		index = int(*doc.Scene)
	}

	for _, rootNodeIndex := range doc.Scenes[index].Nodes {
		var node *ecs.Entity
		if node, err = l.buildNode(doc.Nodes[rootNodeIndex]); err != nil {
			return nil, err
		}
		scene.AddChild(node)
	}

	return scene, nil
}

func (l *GLTFImporter) buildNode(n *gltf.Node) (*ecs.Entity, error) {
	var err error
	// Create a node with defailt values if no values exist
	node := ecs.NewEntity()

	// Build a mesh
	if n.Mesh != nil {
		var mesh *rebound.Mesh
		if mesh, err = l.buildMesh(l.doc.Meshes[*n.Mesh]); err != nil {
			return nil, err
		}
		node.AddComponent(&rebound.RenderComponent{
			Mesh:     mesh,
			Position: [3]float32{0, 0, 0},
			Rotation: [3]float32{0, 0, 0},
			Scale:    [3]float32{1, 1, 1},
		})
	}

	// Create children recursively
	if len(n.Children) > 0 {
		for _, child := range n.Children {
			var childNode *ecs.Entity
			if childNode, err = l.buildNode(l.doc.Nodes[child]); err != nil {
				return nil, err
			}
			node.AddChild(childNode)
		}
	}

	return node, nil
}

func (l *GLTFImporter) buildMesh(m *gltf.Mesh) (*rebound.Mesh, error) {
	var err error
	mesh := &rebound.Mesh{
		Attributes: make([]rebound.Attribute, 0),
		Indices:    make([]uint32, 0),
	}

	for _, primitive := range m.Primitives {
		// Load the indices into the mesh
		if primitive.Indices != nil {
			mesh.Indices = l.loadAccessorU32(int(*primitive.Indices))
		}

		// Load the attributes into the mesh such as Normals, texturecoords, etc
		if len(primitive.Attributes) > 0 {
			for name, index := range primitive.Attributes {
				accessor := l.doc.Accessors[index]
				attribute := rebound.Attribute{Type: attTypes[name], Data: l.loadAccessorF32(int(index)), Size: typeSizes[accessor.Type]}
				mesh.Attributes = append(mesh.Attributes, attribute)
			}
		}

		// Set the material of a mesh
		if primitive.Material != nil {
			mat := l.doc.Materials[*primitive.Material]
			if mesh.Material, err = l.buildMaterial(mat); err != nil {
				return nil, err
			}
		}
	}

	// Load the mesh into the GPU
	rebound.LoadMesh(mesh)

	return mesh, nil
}

func (l *GLTFImporter) buildMaterial(m *gltf.Material) (*rebound.Material, error) {
	material := &rebound.Material{
		Transparent: m.DoubleSided,
	}
	if m.PBRMetallicRoughness.BaseColorTexture != nil {
		texture := l.doc.Textures[m.PBRMetallicRoughness.BaseColorTexture.Index]
		texID, err := rebound.LoadTexture(l.dir + "/" + l.doc.Images[*texture.Source].URI)
		if err != nil {
			return nil, err
		}
		material.BaseColorTexture = &texID

		/*
			if texture.Sampler != nil {
				s := l.Doc.Samplers[*texture.Sampler]
				sampler := &rebound.Sampler{
					MagFilter: int32(s.MagFilter),
					MinFilter: uint16(s.MinFilter),
					WrapS:     uint16(s.WrapS),
					WrapT:     uint16(s.WrapT),
				}
			}
		*/
	}
	rgba := m.PBRMetallicRoughness.BaseColorFactor
	material.BaseColor = [4]float32{float32(rgba.A), float32(rgba.R), float32(rgba.G), float32(rgba.B)}

	return material, nil
}

// loadAccessorF32 loads the float32 values from the buffer
func (l *GLTFImporter) loadAccessorF32(index int) []float32 {
	accessor := l.doc.Accessors[index]
	data := l.loadAccessorData(accessor)
	count := int(accessor.Count) * typeSizes[accessor.Type]
	out := make([]float32, count, count)
	switch accessor.ComponentType {
	case gltf.ComponentByte:
	case gltf.ComponentUbyte:
		for i := 0; i < count; i++ {
			out[i] = float32(data[i])
		}
		break
	case gltf.ComponentShort:
	case gltf.ComponentUshort:
		for i := 0; i < count; i++ {
			out[i] = float32(data[i*2]) + float32(data[i*2+1])*256
		}
		break
	default:
		out = (*[1 << 30]float32)(unsafe.Pointer(&data[0]))[:count]
		break
	}
	return out
}

// loadAccessorU32 loads the uint32 values from the buffer
func (l *GLTFImporter) loadAccessorU32(index int) []uint32 {
	accessor := l.doc.Accessors[index]
	data := l.loadAccessorData(accessor)
	count := int(accessor.Count) * typeSizes[accessor.Type]
	out := make([]uint32, count, count)
	switch accessor.ComponentType {
	case gltf.ComponentByte:
	case gltf.ComponentUbyte:
		for i := 0; i < count; i++ {
			out[i] = uint32(data[i])
		}
		break
	case gltf.ComponentShort:
	case gltf.ComponentUshort:
		for i := 0; i < count; i++ {
			out[i] = uint32(data[i*2]) + uint32(data[i*2+1])*256
		}
		break
	default:
		out = (*[1 << 30]uint32)(unsafe.Pointer(&data[0]))[:count]
		break
	}
	return out
}

// loadAccessorData loads data from an accessor inside the buffer view
func (l *GLTFImporter) loadAccessorData(accessor *gltf.Accessor) []uint8 {
	bv := l.doc.BufferViews[*accessor.BufferView]
	buffer := l.doc.Buffers[bv.Buffer]
	data := buffer.Data[bv.ByteOffset : bv.ByteOffset+bv.ByteLength]
	if accessor.ByteOffset != 0 {
		return data[accessor.ByteOffset:]
	}
	return data
}

// typeSizes returns the size of each type
var typeSizes = map[gltf.AccessorType]int{
	gltf.AccessorScalar: 1,
	gltf.AccessorVec2:   2,
	gltf.AccessorVec3:   3,
	gltf.AccessorVec4:   4,
	gltf.AccessorMat2:   4,
	gltf.AccessorMat3:   9,
	gltf.AccessorMat4:   16,
}

// attTypes returns the Rebound attribute type based on the GLTF attribute string
var attTypes = map[string]rebound.AttributeType{
	"TEXCOORD_0": rebound.TEXCOORDS0,
	"TEXCOORD_1": rebound.TEXCOORDS1,
	"COLOR_0":    rebound.COLOR,
	"JOINTS_0":   rebound.JOINTS,
	"WEIGHTS_0":  rebound.WEIGHTS,
	"NORMAL":     rebound.NORMALS,
	"TANGENT":    rebound.TANGENTS,
	"POSITION":   rebound.POSITION,
}

var emptyMatrix = [16]float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

func toFloat32Array(input [16]float64) [16]float32 {
	var result [16]float32

	for index, f64 := range input {
		result[index] = float32(f64)
	}

	return result
}
