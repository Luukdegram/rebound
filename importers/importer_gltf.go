package importers

import (
	"path"
	"unsafe"

	"github.com/luukdegram/rebound"

	"github.com/qmuntal/gltf"
)

// GLTFImporter loads a gltf file into Rebound to create a visual object.
type GLTFImporter struct {
	Dir string
	Doc *gltf.Document
}

// Import loads a GLTF file into a Scene.
func (l *GLTFImporter) Import(file string) (*rebound.SceneComponent, error) {
	doc, err := gltf.Open(file)
	if err != nil {
		return nil, err
	}
	l.Doc = doc

	dir := path.Dir(file)
	scene := &rebound.SceneComponent{
		Nodes: make([]*rebound.Node, 0),
	}
	l.Dir = dir

	for _, rootNodeIndex := range doc.Scenes[*doc.Scene].Nodes {
		var node *rebound.Node
		if node, err = l.buildNode(doc.Nodes[rootNodeIndex]); err != nil {
			return nil, err
		}

		scene.Nodes = append(scene.Nodes, node)
	}

	return scene, nil
}

func (l *GLTFImporter) buildNode(n gltf.Node) (*rebound.Node, error) {
	var err error
	// Create a node with defailt values if no values exist
	node := &rebound.Node{
		Name: n.Name,
	}

	if n.Matrix == emptyMatrix {
		node.Transformation = rebound.NewTransformationMatrix(
			n.TranslationOrDefault(),
			n.RotationOrDefault(),
			n.ScaleOrDefault(),
		)
	} else {
		node.Transformation = toFloat32Array(n.Matrix)
	}

	// Build a mesh
	if n.Mesh != nil {
		if node.Mesh, err = l.buildMesh(l.Doc.Meshes[*n.Mesh]); err != nil {
			return nil, err
		}
	}

	// Create children recursively
	if len(n.Children) > 0 {
		for _, child := range n.Children {
			var childNode *rebound.Node
			if childNode, err = l.buildNode(l.Doc.Nodes[child]); err != nil {
				return nil, err
			}
			node.Children = append(node.Children, childNode)
		}
	}

	return node, nil
}

func (l *GLTFImporter) buildMesh(m gltf.Mesh) (*rebound.Mesh, error) {
	mesh := &rebound.Mesh{
		Attributes: make([]rebound.Attribute, 0),
		Indices:    make([]uint32, 0),
	}

	for _, primitive := range m.Primitives {
		// Load the indices into the mesh
		if primitive.Indices != nil {
			mesh.Indices = append(mesh.Indices, l.loadAccessorU32(int(*primitive.Indices))...)
		}

		// Load the attributes into the mesh such as Normals, texturecoords, etc
		if len(primitive.Attributes) > 0 {
			for name, index := range primitive.Attributes {
				accessor := l.Doc.Accessors[index]
				attribute := rebound.Attribute{Type: attTypes[name], Data: l.loadAccessorF32(int(index)), Size: typeSizes[accessor.Type]}
				mesh.Attributes = append(mesh.Attributes, attribute)
			}
		}

		// Set the material of a mesh
		if primitive.Material != nil {
			material := &rebound.Material{}
			if l.Doc.Materials[*primitive.Material].PBRMetallicRoughness.BaseColorTexture != nil {
				textureSource := l.Doc.Textures[(l.Doc.Materials[*primitive.Material].PBRMetallicRoughness.BaseColorTexture).Index].Source
				texID, err := rebound.LoadTexture(l.Dir + "/" + l.Doc.Images[*textureSource].URI)
				if err != nil {
					return nil, err
				}
				material.BaseColorTexture = rebound.NewTextureComponent(texID)
			}

			rgba := l.Doc.Materials[*primitive.Material].PBRMetallicRoughness.BaseColorFactor
			material.BaseColor = [4]float32{float32(rgba.A), float32(rgba.R), float32(rgba.G), float32(rgba.B)}
			mesh.Material = material
		}
	}

	// Load the mesh into the GPU
	rebound.LoadToVAO(mesh)

	return mesh, nil
}

// loadAccessorF32 loads the float32 values from the buffer
func (l *GLTFImporter) loadAccessorF32(index int) []float32 {
	accessor := l.Doc.Accessors[index]
	data := l.loadAccessorData(accessor)
	count := int(accessor.Count) * typeSizes[accessor.Type]
	out := make([]float32, count, count)
	switch accessor.ComponentType {
	case gltf.UnsignedByte:
		for i := 0; i < int(accessor.Count); i++ {
			out[i] = float32(data[i])
		}
		break
	case gltf.UnsignedShort:
		for i := 0; i < int(accessor.Count); i++ {
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
	accessor := l.Doc.Accessors[index]
	data := l.loadAccessorData(accessor)
	count := int(accessor.Count) * typeSizes[accessor.Type]
	out := make([]uint32, count, count)
	switch accessor.ComponentType {
	case gltf.UnsignedByte:
		for i := 0; i < int(accessor.Count); i++ {
			out[i] = uint32(data[i])
		}
		break
	case gltf.UnsignedShort:
		for i := 0; i < int(accessor.Count); i++ {
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
func (l *GLTFImporter) loadAccessorData(accessor gltf.Accessor) []uint8 {
	bv := l.Doc.BufferViews[*accessor.BufferView]
	buffer := l.Doc.Buffers[bv.Buffer]
	return buffer.Data[bv.ByteOffset : bv.ByteOffset+bv.ByteLength]
}

// typeSizes returns the size of each type
var typeSizes = map[gltf.AccessorType]int{
	gltf.Scalar: 1,
	gltf.Vec2:   2,
	gltf.Vec3:   3,
	gltf.Vec4:   4,
	gltf.Mat2:   4,
	gltf.Mat3:   9,
	gltf.Mat4:   16,
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
