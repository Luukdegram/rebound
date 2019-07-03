package main

import (
	"fmt"

	"github.com/go-gl/gl/v4.6-core/gl"
)

var (
	vaos []uint32
	vbos []uint32
)

func loadToVAO(positions []float32) *rawModel {
	id := createVAO(positions)
	return &rawModel{vaoID: id, vertextCount: len(positions) / 3}
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
	gl.BindVertexArray(0)

	vaos = append(vaos, vao)
	vbos = append(vbos, vbo)

	return vao
}

func cleanUp() {
	i := 0
	for ; i < len(vaos); i++ {
		gl.DeleteVertexArrays(1, &vaos[i])
	}

	i = 0
	for ; i < len(vbos); i++ {
		gl.DeleteVertexArrays(1, &vbos[i])
	}

	fmt.Println("Clean up of vertex arrays")
}
