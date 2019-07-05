package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

func render(model *rawModel) {
	gl.BindVertexArray(model.vaoID)
	gl.EnableVertexAttribArray(0)
	gl.DrawElements(gl.TRIANGLES, int32(model.vertextCount), gl.UNSIGNED_INT, gl.Ptr(nil))
	gl.DisableVertexArrayAttrib(model.vaoID, 0)
	gl.BindVertexArray(0)
}
