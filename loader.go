package main

import (
	"fmt"
	"image"
	"image/draw"
	"os"

	_ "image/png"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/luukdegram/rebound/models"
)

var (
	ids []uint32
)

func loadToVAO(points []float32, textureCoords []float32, indices []uint32) *models.RawModel {
	id := createVAO(points)
	bindIndicesBuffer(indices)
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
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
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

	ids = append(ids, texture)
	return texture, nil
}

func createVAO(points []float32) uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	ids = append(ids, vao)
	ids = append(ids, vbo)

	return vao
}

func bindIndicesBuffer(indices []uint32) {
	var ebo uint32
	gl.GenBuffers(1, &ebo)
	ids = append(ids, ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(indices), gl.Ptr(indices), gl.STATIC_DRAW)
}

func unbindVAO() {
	gl.BindVertexArray(0)
}

func cleanUp() {
	for _, id := range ids {
		gl.DeleteVertexArrays(1, &id)
	}
}
