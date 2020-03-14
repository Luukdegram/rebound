package rebound

import (
	"fmt"
	"image"
	"image/draw"
	"os"

	_ "image/jpeg" //Import jpg package to be able to decode jpg files
	_ "image/png"  //Import png package to be able to decode png files

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/luukdegram/rebound/internal/thread"
)

var (
	vaos     []uint32
	vbos     []uint32
	textures map[string]uint32 = make(map[string]uint32)
)

//LoadMesh creates a new vao and stores the mesh data inside its buffer
func LoadMesh(m *Mesh) {
	var id uint32
	thread.Call(func() {
		id = createVAO()
		if len(m.Indices) > 0 {
			bindIndicesBuffer(m.Indices)
		}
		for _, attribute := range m.Attributes {
			storeDataInAttributeList(int(attribute.Type), attribute.Size, attribute.Data)
		}
		unbindVAO()
	})

	m.ID = id
}

//LoadTexture loads a texture into the GPU
func LoadTexture(fileName string) (uint32, error) {
	// Return the texture if we already loaded it before. This increases performance as loading textures is quite intensive.
	if val, exists := textures[fileName]; exists {
		return val, nil
	}

	rgba, err := loadTextureData(fileName)
	if err != nil {
		return 0, err
	}
	var texture uint32
	thread.Call(func() {
		gl.GenTextures(1, &texture)
		gl.BindTexture(gl.TEXTURE_2D, texture)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		gl.TexImage2D(
			gl.TEXTURE_2D,
			0,
			gl.RGBA,
			int32(rgba.Rect.Size().X),
			int32(rgba.Rect.Size().Y),
			0,
			gl.RGBA,
			gl.UNSIGNED_BYTE,
			gl.Ptr(rgba.Pix))
		gl.GenerateMipmap(gl.TEXTURE_2D)
	})

	// Save the new texture into the texture map
	textures[fileName] = texture

	return texture, nil
}

// LoadCubeMap loads a cubemap into a GPU texture, returns the index of the texture as an unsigned 32bit integer.
func LoadCubeMap(faces [6]string) (uint32, error) {
	data := make([]*image.RGBA, len(faces))
	for index := 0; index < len(faces); index++ {
		rgba, err := loadTextureData(faces[index])
		if err != nil {
			return 0, err
		}
		data[index] = rgba
	}

	var texture uint32
	thread.Call(func() {
		gl.GenTextures(1, &texture)
		gl.BindTexture(gl.TEXTURE_CUBE_MAP, texture)
		for index := 0; index < len(data); index++ {
			gl.TexImage2D(gl.TEXTURE_CUBE_MAP_POSITIVE_X+uint32(index),
				0,
				gl.RGBA,
				int32(data[index].Rect.Size().X),
				int32(data[index].Rect.Size().Y),
				0,
				gl.RGBA,
				gl.UNSIGNED_BYTE,
				gl.Ptr(data[index].Pix))
		}

		gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)
	})

	return texture, nil

}

func loadTextureData(fileName string) (*image.RGBA, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return nil, fmt.Errorf("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	return rgba, nil
}

func createVAO() uint32 {
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	vaos = append(vaos, vao)
	return vao
}

func storeDataInAttributeList(index int, coordinateSize int, data []float32) {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	vbos = append(vbos, vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(data), gl.Ptr(data), gl.STATIC_DRAW)
	gl.VertexAttribPointer(uint32(index), int32(coordinateSize), gl.FLOAT, false, 0, nil)
}

func bindIndicesBuffer(indices []uint32) {
	var ebo uint32
	gl.GenBuffers(1, &ebo)
	vbos = append(vbos, ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(indices), gl.Ptr(indices), gl.STATIC_DRAW)
}

func unbindVAO() {
	gl.BindVertexArray(0)
}

//CleanUp removes all loaded data from GPU to free up space.
//As this removes all data, only run this when shutting down.
func CleanUp() {
	thread.Call(func() {
		gl.DeleteVertexArrays(int32(len(vaos)), &vaos[0])
		gl.DeleteBuffers(int32(len(vbos)), &vbos[0])
		for _, id := range textures {
			gl.DeleteTextures(1, &id)
		}

		vaos = []uint32{}
		vbos = []uint32{}
		textures = make(map[string]uint32)
	})
}

// Sampler describes how to render a texture
type Sampler struct {
	MagFilter int32
	MinFilter uint16
	WrapS     uint16
	WrapT     uint16
}
