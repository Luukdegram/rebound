package rebound

import (
	"fmt"
	"image"
	"image/draw"
	"os"

	_ "image/png" //Import png package to be able to decode png files

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/luukdegram/rebound/models"
)

var (
	vaos     []uint32
	vbos     []uint32
	textures []uint32
)

//LoadToVAO test
func LoadToVAO(indices []uint32, attributes []Attribute) *models.RawModel {
	id := createVAO()
	bindIndicesBuffer(indices)
	for _, attribute := range attributes {
		storeDataInAttributeList(int(attribute.Type), attribute.Size, attribute.Data)
	}
	unbindVAO()
	return &models.RawModel{VaoID: id, VertexCount: len(indices)}
}

//LoadTexture loads a texture into the GPU
func LoadTexture(fileName string) (uint32, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return 0, err
	}
	img, _, err := image.Decode(file)
	if err != nil {
		return 0, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return 0, fmt.Errorf("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	var texture uint32
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

	textures = append(textures, texture)
	return texture, nil
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
func CleanUp() {
	for _, id := range vaos {
		gl.DeleteVertexArrays(1, &id)
	}

	for _, id := range vbos {
		gl.DeleteBuffers(1, &id)
	}

	for _, id := range textures {
		gl.DeleteTextures(1, &id)
	}
}
