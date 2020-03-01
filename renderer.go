package rebound

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/luukdegram/rebound/ecs"
	"github.com/luukdegram/rebound/shaders"
)

const (
	// RenderComponentName is the name of a RenderComponent
	RenderComponentName = "RenderComponent"
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
	skyColor    mgl32.Vec3
}

//RenderSystem handles the rendering of all entities
type RenderSystem struct {
	ecs.BaseSystem
	FOV         float32
	NearPlane   float32
	FarPlane    float32
	drawPolygon bool
	Camera      *Camera
	Light       *Light
	pm          *mgl32.Mat4
	skyColor    mgl32.Vec3
}

//NewRenderSystem returns a new RendererSystem with default settings
func NewRenderSystem() *RenderSystem {
	rs := &RenderSystem{
		BaseSystem:  ecs.NewBaseSystem(),
		FOV:         45,
		NearPlane:   0.1,
		FarPlane:    100,
		drawPolygon: false,
		skyColor:    mgl32.Vec3{0, 0, 0},
	}
	enableCulling()
	return rs
}

//Check returns true if the given Entity contains both a ShaderComponent and RenderComponent
func (rs *RenderSystem) Check(e *ecs.Entity) bool {
	return e.HasComponent(shaders.ShaderComponentName) && e.HasComponent(RenderComponentName)
}

//Update draws all entities within a RendererSystem
func (rs *RenderSystem) Update(dt float32) {
	rs.prepare()
	for _, e := range rs.Entities() {
		if e.HasComponent(shaders.ShaderComponentName) {
			sc := *e.Component(shaders.ShaderComponentName).(*shaders.ShaderComponent)
			shaders.Start(sc)
			shaders.LoadVec3(sc, "lightPos", rs.Light.Position)
			shaders.LoadVec3(sc, "lightColour", rs.Light.Colour)
			shaders.LoadVec3(sc, "skyColour", rs.skyColor)

			if rs.Camera != nil {
				shaders.LoadMat(sc, "projectionMatrix", *rs.pm)
				shaders.LoadMat(sc, "viewMatrix", NewViewMatrix(*rs.Camera))
			}

			if rc := e.Component(RenderComponentName).(*RenderComponent); rc != nil {
				rc.Renderable.Render(sc, *rs.Camera)
			}
		}
	}
}

//RenderComponent holds the data to render an entity
type RenderComponent struct {
	Renderable
}

//Name returns the RenderComponent name
func (rc *RenderComponent) Name() string {
	return RenderComponentName
}

//Renderable is an object that can be rendered
type Renderable interface {
	//Render is a function that is called by the RenderSystem which can be used to modify the parameters in the shader
	Render(shaders.ShaderComponent, Camera)
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
	r.FOV = 45
	r.NearPlane = 0.1
	r.FarPlane = 100
	r.drawPolygon = false
	r.registry = new(registry)
	r.registry.entries = make(map[string]*entry)
	r.skyColor = mgl32.Vec3{0, 0, 0}
	enableCulling()
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
func (r *Renderer) NewLight(position mgl32.Vec3) {
	if r.Light == nil {
		r.Light = &Light{Position: position, Colour: mgl32.Vec3{1, 1, 1}}
	}
}

//SetSkyColor sets the color of the sky
func (r *Renderer) SetSkyColor(red, green, blue float32) {
	r.skyColor = mgl32.Vec3{red, green, blue}
}

//Prepare cleans the screen for the next draw
func (rs *RenderSystem) prepare() {
	gl.Enable(gl.DEPTH_TEST)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.ClearColor(rs.skyColor[0], rs.skyColor[1], rs.skyColor[2], 1.0)
	if rs.drawPolygon {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	} else {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
	}
}

//TogglePolygons enables/disables drawing the polygons
func (rs *RenderSystem) TogglePolygons() {
	rs.drawPolygon = !rs.drawPolygon
}

//Render draws a 3D model into the screen
func (r Renderer) Render(shader shaders.ShaderProgram) {
	shader.Start()
	shader.LoadVec3("lightPos", r.Light.Position)
	shader.LoadVec3("lightColour", r.Light.Colour)
	shader.LoadVec3("skyColour", r.skyColor)

	if r.Camera != nil {
		shader.LoadMat("projectionMatrix", *r.pm)
		shader.LoadMat("viewMatrix", NewViewMatrix(*r.Camera))
	}

	for _, entry := range r.registry.entries {
		prepareMesh(*entry.mesh, shader)
		for _, entity := range entry.entities {
			prepareInstance(*entity, shader) // Set transformation matrix
			gl.DrawElements(gl.TRIANGLES, int32(entry.mesh.RawModel.VertexCount), gl.UNSIGNED_INT, gl.Ptr(nil))
		}
		unbindMesh(*entry.mesh)
	}

	shader.Stop()
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

func enableCulling() {
	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)
}

func disableCulling() {
	gl.Disable(gl.CULL_FACE)
}

func prepareMesh(mesh Mesh, shader shaders.ShaderProgram) {
	gl.BindVertexArray(mesh.RawModel.VaoID)
	for _, attr := range mesh.attributes {
		gl.EnableVertexAttribArray(uint32(attr.Type))
	}

	if mesh.IsTextured() {
		if mesh.Texture.HasTransparancy() {
			disableCulling()
		}
		shader.LoadBool("useFakeLighting", mesh.Texture.UseFakeLighting())
		shader.LoadFloat("shineDamper", mesh.Texture.ShineDamper)
		shader.LoadFloat("reflectivity", mesh.Texture.Reflectivity)
		gl.BindTexture(gl.TEXTURE_2D, mesh.Texture.ID)
	}
}

func unbindMesh(mesh Mesh) {
	enableCulling()
	for _, attr := range mesh.attributes {
		gl.DisableVertexAttribArray(uint32(attr.Type))
	}
	gl.BindVertexArray(0)
}

func prepareInstance(entity Entity, shader shaders.ShaderProgram) {
	transform := NewTransformationMatrix(entity.Position, entity.Rotation, entity.Scale)
	shader.LoadMat("transformMatrix", transform)
}
