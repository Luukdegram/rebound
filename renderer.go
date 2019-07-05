package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/luukdegram/rebound/models"
)

//Renderer can render models and setups the canvas
type Renderer struct{}

//Prepare cleans the screen for the next draw
func (r *Renderer) Prepare() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

//Render renders a model
func (r *Renderer) Render(TexturedModel *models.TexturedModel) {
	model := TexturedModel.Model
	
	gl.BindVertexArray(model.VaoID)
	gl.EnableVertexAttribArray(0)
	gl.DrawElements(gl.TRIANGLES, int32(model.VertextCount), gl.UNSIGNED_INT, gl.Ptr(nil))
	gl.DisableVertexArrayAttrib(model.VaoID, 0)
	gl.BindVertexArray(0)
}
