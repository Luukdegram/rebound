package main

import (
	"rebound/models"

	"github.com/go-gl/gl/v4.1-core/gl"
)

//Renderer can render models and setups the canvas
type Renderer struct{}

//Prepare cleans the screen for the next draw
func (r *Renderer) Prepare() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

//Render renders a model
func (r *Renderer) Render(texturedModel *models.TexturedModel) {
	model := texturedModel.Model
	gl.BindVertexArray(model.VaoID)
	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texturedModel.Texture.ID)
	gl.DrawElements(gl.TRIANGLES, int32(model.VertextCount), gl.UNSIGNED_INT, gl.Ptr(nil))
	gl.DisableVertexAttribArray(0)
	gl.DisableVertexAttribArray(1)
	gl.BindVertexArray(0)
}
