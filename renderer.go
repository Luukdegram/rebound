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
	registry    *registry
}

type registry struct {
	entries map[string]*entry
}

type entry struct {
	mesh     *Mesh
	entities []*Entity
}

//NewRenderer returns a new Renderer object.
func NewRenderer() *Renderer {
	r := new(Renderer)
	r.FOV = 70
	r.NearPlane = 0.1
	r.FarPlane = 1000
	r.drawPolygon = false
	r.registry = new(registry)
	r.registry.entries = make(map[string]*entry)
	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)
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
func (r Renderer) prepare() {
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
func (r Renderer) Render(shader shaders.ShaderProgram) {
	r.prepare()
	shader.Start()
	shader.LoadVec3(shader.GetUniformLocation("lightPos"), r.Light.Position)
	shader.LoadVec3(shader.GetUniformLocation("lightColour"), r.Light.Colour)

	if r.Camera != nil {
		shader.LoadMat(shader.GetUniformLocation("projectionMatrix"), *r.pm)
		shader.LoadMat(shader.GetUniformLocation("viewMatrix"), NewViewMatrix(*r.Camera))
	}

	for _, entry := range r.registry.entries {
		prepareMesh(*entry.mesh, shader)
		for _, entity := range entry.entities {
			prepareInstance(*entity, shader)
			gl.DrawElements(gl.TRIANGLES, int32(entry.mesh.RawModel.VertexCount), gl.UNSIGNED_INT, gl.Ptr(nil))
		}
		unbindMesh(*entry.mesh)
	}

	shader.Stop()
	r.registry.entries = make(map[string]*entry)
}

//RegisterEntity registers entities to the renderer
func (r *Renderer) RegisterEntity(entities ...*Entity) {
	for _, entity := range entities {
		for _, mesh := range entity.Geometry.Meshes {
			if val, ok := r.registry.entries[mesh.Name]; ok {
				val.entities = append(val.entities, entity)
			} else {
				r.registry.entries[mesh.Name] = &entry{mesh: &mesh, entities: []*Entity{entity}}
			}
		}
	}
}

func prepareMesh(mesh Mesh, shader shaders.ShaderProgram) {
	gl.BindVertexArray(mesh.RawModel.VaoID)
	for _, attr := range mesh.attributes {
		gl.EnableVertexAttribArray(uint32(attr.Type))
	}

	if mesh.IsTextured() {
		gl.BindTexture(gl.TEXTURE_2D, mesh.Texture.ID)
		shader.LoadFloat(shader.GetUniformLocation("shineDamper"), mesh.Texture.ShineDamper)
		shader.LoadFloat(shader.GetUniformLocation("reflectivity"), mesh.Texture.Reflectivity)
	}
}
func unbindMesh(mesh Mesh) {
	for _, attr := range mesh.attributes {
		gl.DisableVertexAttribArray(uint32(attr.Type))
	}
	gl.BindVertexArray(0)
}
func prepareInstance(entity Entity, shader shaders.ShaderProgram) {
	transform := NewTransformationMatrix(entity.Position, entity.Rotation, entity.Scale)
	shader.LoadMat(shader.GetUniformLocation("transformMatrix"), transform)
}
