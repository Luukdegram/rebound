package importers

import (
	"unsafe"

	"github.com/luukdegram/rebound"
	"github.com/luukdegram/rebound/models"

	"github.com/qmuntal/gltf"
)

//LoadGltfModel imports a GLTF file into a model
func LoadGltfModel(file string) []models.RawModel {
	doc, err := gltf.Open(file)
	if err != nil {
		panic(err)
	}
	models := make([]models.RawModel, 0)

	for _, mesh := range doc.Meshes {
		indices := make([]uint32, 0, 0)
		attributes := make([]rebound.Attribute, 0, 0)

		for _, primitive := range mesh.Primitives {
			if primitive.Indices != nil {
				indices = append(indices, loadIndices(doc, int(*primitive.Indices))...)
			}

			if len(primitive.Attributes) > 0 {
				for name, index := range primitive.Attributes {
					accessor := doc.Accessors[index]
					attributes = append(attributes, rebound.Attribute{Name: name, Data: loadAccessorF32(doc, int(index)), Size: typeSizes[accessor.Type]})
				}
			}

		}

		models = append(models, *rebound.LoadToVAO2(indices, attributes))
	}
	return models
}

func loadIndices(doc *gltf.Document, index int) []uint32 {
	return loadAccessorU32(doc, index)
}

func loadAccessorF32(doc *gltf.Document, index int) []float32 {
	accessor := doc.Accessors[index]
	data := loadAccessorData(doc, accessor)
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

func loadAccessorU32(doc *gltf.Document, index int) []uint32 {
	accessor := doc.Accessors[index]
	data := loadAccessorData(doc, accessor)
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

func loadAccessorData(doc *gltf.Document, accessor gltf.Accessor) []uint8 {
	bv := doc.BufferViews[*accessor.BufferView]
	buffer := doc.Buffers[bv.Buffer]
	return buffer.Data[bv.ByteOffset : bv.ByteOffset+bv.ByteLength]
}

var typeSizes = map[gltf.AccessorType]int{
	gltf.Scalar: 1,
	gltf.Vec2:   2,
	gltf.Vec3:   3,
	gltf.Vec4:   4,
	gltf.Mat2:   4,
	gltf.Mat3:   9,
	gltf.Mat4:   16,
}
