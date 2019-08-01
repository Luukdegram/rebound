package rebound

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

//Renderer can render models and setups the canvas
type Renderer struct {
	FOV         float32
	NearPlane   float32
	FarPlane    float32
	drawPolygon bool
}

//NewRenderer returns a new Renderer object.
func NewRenderer() *Renderer {
	r := new(Renderer)
	r.FOV = 70
	r.NearPlane = 0.1
	r.FarPlane = 1000
	r.drawPolygon = false
	return r
}

//Prepare cleans the screen for the next draw
func (r Renderer) Prepare() {
	gl.Enable(gl.DEPTH_TEST)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

//TogglePolygons enables/disables drawing the polygons
func (r *Renderer) TogglePolygons() {
	r.drawPolygon = !r.drawPolygon
}

//Render draws a 3D model into the screen
func (r Renderer) Render(geometry *Geometry) {
	if r.drawPolygon {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	} else {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
	}

	for _, mesh := range geometry.Meshes {
		gl.BindVertexArray(mesh.RawModel.VaoID)
		for index := range mesh.attributes {
			gl.EnableVertexAttribArray(uint32(index))
		}
		if mesh.IsTextured() {
			gl.BindTexture(gl.TEXTURE_2D, mesh.Texture.ID)
		}
		gl.DrawElements(gl.TRIANGLES, int32(mesh.RawModel.VertexCount), gl.UNSIGNED_INT, gl.Ptr(nil))
		for index := range mesh.attributes {
			gl.DisableVertexAttribArray(uint32(index))
		}
		gl.BindVertexArray(0)
	}
}
