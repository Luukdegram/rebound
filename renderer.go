package main

import (
	"github.com/go-gl/gl/v4.6-core/gl"
)

func render(model *rawModel) {
	gl.BindVertexArray(model.vaoID)
	gl.EnableVertexAttribArray(0)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(model.vertextCount))
	gl.DisableVertexArrayAttrib(model.vaoID, 0)
	gl.BindVertexArray(0)
}
