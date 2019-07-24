package main

import (
	"fmt"
	"image"
	"image/draw"
	"os"
	"rebound/models"

	_ "image/png"

	"github.com/go-gl/gl/v4.1-core/gl"
)

var (
	vaos     []uint32
	vbos     []uint32
	textures []uint32
)

//LoadToVAO loads a model to the gpu
func LoadToVAO(points []float32, textureCoords []float32, indices []uint32) *models.RawModel {
	id := createVAO()
	bindIndicesBuffer(indices)
	storeDataInAttributeList(0, 3, points)
	storeDataInAttributeList(1, 2, textureCoords)
	unbindVAO()
	return &models.RawModel{VaoID: id, VertextCount: len(indices)}
}

func loadTexture(fileName string) (uint32, error) {
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
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
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
	vbos = append(vaos, vao)
	return vao
}

func storeDataInAttributeList(index uint32, coordinateSize int32, data []float32) {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	vbos = append(vbos, vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(data), gl.Ptr(data), gl.STATIC_DRAW)
	gl.VertexAttribPointer(index, coordinateSize, gl.FLOAT, false, 0, nil)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
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

func cleanUp() {
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
