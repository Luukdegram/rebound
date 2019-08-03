package rebound

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/luukdegram/rebound/shaders"
)

//Renderer can render models and setups the canvas
type Renderer struct {
	FOV         float32
	NearPlane   float32
	FarPlane    float32
	drawPolygon bool
	Camera      *Camera
	Light       *Light
	pm          *mgl32.Mat4
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

//NewCamera creates a new camera and attaches it to the renderer
func (r *Renderer) NewCamera(width int, height int) {
	if r.Camera == nil {
		r.Camera = new(Camera)
	}

	pm := NewProjectionMatrix(r.FOV, float32(width/height), r.NearPlane, r.FarPlane)
	r.pm = &pm
}

//NewLight Adds Light to the renderer
func (r *Renderer) NewLight() {
	if r.Light == nil {
		r.Light = &Light{Position: mgl32.Vec3{0, 0, 0}, Colour: mgl32.Vec3{1, 1, 1}}
	}
}

//Prepare cleans the screen for the next draw
func (r Renderer) Prepare() {
	gl.Enable(gl.DEPTH_TEST)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	if r.drawPolygon {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	} else {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
	}
}

//TogglePolygons enables/disables drawing the polygons
func (r *Renderer) TogglePolygons() {
	r.drawPolygon = !r.drawPolygon
}

//Render draws a 3D model into the screen
func (r Renderer) Render(entity Entity, shader shaders.ShaderProgram) {
	transform := NewTransformationMatrix(entity.Position, entity.Rotation, entity.Scale)
	shader.LoadMat(shader.GetUniformLocation("transformMatrix"), transform)

	if r.Camera != nil {
		shader.LoadMat(shader.GetUniformLocation("projectionMatrix"), *r.pm)
		shader.LoadMat(shader.GetUniformLocation("viewMatrix"), NewViewMatrix(*r.Camera))
	}

	for _, mesh := range entity.Geometry.Meshes {
		gl.BindVertexArray(mesh.RawModel.VaoID)
		for _, attr := range mesh.attributes {
			gl.EnableVertexAttribArray(uint32(attr.Type))
		}
		if mesh.IsTextured() {
			gl.BindTexture(gl.TEXTURE_2D, mesh.Texture.ID)
			shader.LoadFloat(shader.GetUniformLocation("shineDamper"), mesh.Texture.ShineDamper)
			shader.LoadFloat(shader.GetUniformLocation("reflectivity"), mesh.Texture.Reflectivity)
		}

		gl.DrawElements(gl.TRIANGLES, int32(mesh.RawModel.VertexCount), gl.UNSIGNED_INT, gl.Ptr(nil))
		for _, attr := range mesh.attributes {
			gl.DisableVertexAttribArray(uint32(attr.Type))
		}
		gl.BindVertexArray(0)
	}
}
