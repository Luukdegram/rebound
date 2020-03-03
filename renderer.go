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
	//Pos is a shader attribute used for positional coordinates
	Pos AttributeType = iota
	//TexCoords is a shader attribute used for the coordinates of a texture
	TexCoords
	//Normals holds the coordinates of a normal texture
	Normals
	//Tangents holds the tangets data in a shader
	Tangents
)

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
	skyColor    [3]float32
}

//Attribute is vbo that stores data such as texture coordinates
type Attribute struct {
	Type AttributeType
	Data []float32
	Size int
}

//AttributeType can be used to link external attribute names to Rebound's
type AttributeType int

// RenderComponent holds the data to render an entity
type RenderComponent struct {
	Renderable
	vaoID       uint32
	vertexCount int
	attributes  []Attribute
	// Rotation holds the rotational data of the render object related to the 3D world
	Rotation [3]float32
	// Position holds the positional data of the render object related to the 3D world
	Position [3]float32
	// Scale allows the object to be scaled
	Scale float32
}

// Renderable is an object that can be rendered
type Renderable interface {
	//Render is a function that is called by the RenderSystem which can be used to modify the parameters in the shader
	Render(ecs.Entity, shaders.ShaderComponent, Camera)
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
	for _, e := range rs.BaseSystem.Entities() {
		if !rs.Check(e) {
			continue
		}

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
			if rc.Renderable != nil {
				rc.Renderable.Render(*e, sc, *rs.Camera)
			}
			prepareRenderable(*e)
			gl.DrawElements(gl.TRIANGLES, int32(rc.vertexCount), gl.UNSIGNED_INT, gl.PtrOffset(0))
			unbindRenderable(*e)
		}
		shaders.Stop()
	}
}

// Name returns the RenderComponent name
func (rc *RenderComponent) Name() string {
	return RenderComponentName
}

//NewCamera creates a new camera and attaches it to the renderer
func (rs *RenderSystem) NewCamera(width int, height int) {
	if rs.Camera == nil {
		rs.Camera = new(Camera)
	}

	pm := NewProjectionMatrix(rs.FOV, float32(width/height), rs.NearPlane, rs.FarPlane)
	rs.pm = &pm
}

//NewLight Adds Light to the renderer
func (rs *RenderSystem) NewLight(position [3]float32) {
	if rs.Light == nil {
		rs.Light = &Light{Position: position, Colour: [3]float32{1, 1, 1}}
	}
}

//SetSkyColor sets the color of the sky
func (rs *RenderSystem) SetSkyColor(red, green, blue float32) {
	rs.skyColor = [3]float32{red, green, blue}
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

func enableCulling() {
	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)
}

func disableCulling() {
	gl.Disable(gl.CULL_FACE)
}

func prepareRenderable(e ecs.Entity) {
	rc := e.Component(RenderComponentName).(*RenderComponent)
	sc := *e.Component(shaders.ShaderComponentName).(*shaders.ShaderComponent)
	gl.BindVertexArray(rc.vaoID)

	for _, a := range rc.attributes {
		gl.EnableVertexAttribArray(uint32(a.Type))
	}

	if e.HasComponent(TextureComponentName) {
		tc := e.Component(TextureComponentName).(*TextureComponent)
		if tc.Transparant {
			disableCulling()
		}

		shaders.LoadBool(sc, "useFakeLighting", tc.Transparant)
		shaders.LoadFloat(sc, "shineDamper", tc.ShineDamper)
		shaders.LoadFloat(sc, "reflectivity", tc.Reflectivity)
		gl.BindTexture(gl.TEXTURE_2D, tc.id)
	}

	transform := NewTransformationMatrix(rc.Position, rc.Rotation, rc.Scale)
	shaders.LoadMat(sc, "transformMatrix", transform)
}

func unbindRenderable(e ecs.Entity) {
	enableCulling()
	rc := e.Component(RenderComponentName).(*RenderComponent)
	for _, a := range rc.attributes {
		gl.DisableVertexAttribArray(uint32(a.Type))
	}
	gl.BindVertexArray(0)
}
