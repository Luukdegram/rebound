package rebound

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/luukdegram/rebound/ecs"
	"github.com/luukdegram/rebound/internal/thread"
	"github.com/luukdegram/rebound/shaders"
)

const (
	// RenderComponentName is the name of a RenderComponent
	RenderComponentName = "RenderComponent"
)

const (
	//POSITION is a shader attribute used for positional coordinates
	POSITION AttributeType = iota
	//TEXCOORDS0 is a shader attribute used for the coordinates of a base texture
	TEXCOORDS0
	//TEXCOORDS1 is a shader attribute used for the coordinates of a second base texture
	TEXCOORDS1
	//NORMALS holds the coordinates of a normal texture
	NORMALS
	//TANGENTS holds the tangets data in a shader
	TANGENTS
	//COLOR holds the color data of a base texture
	COLOR
	//JOINTS holds the joints data of a skinned mesh
	JOINTS
	//WEIGHTS holds the weights data of a skinned mesh
	WEIGHTS
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
	vaoID      uint32
	attributes []Attribute
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
		FOV:         105,
		NearPlane:   0.1,
		FarPlane:    100,
		drawPolygon: false,
		skyColor:    mgl32.Vec3{0, 0, 0},
	}
	thread.Call(enableCulling)
	return rs
}

//Check returns true if the given Entity contains both a ShaderComponent and RenderComponent
func (rs *RenderSystem) Check(e *ecs.Entity) bool {
	return e.HasComponent(SceneComponentName)
}

//Update draws all entities within a RendererSystem
func (rs *RenderSystem) Update(dt float64) {
	thread.Call(func() {
		rs.prepare()
		for _, e := range rs.BaseSystem.Entities() {
			scene := e.Component(SceneComponentName).(*SceneComponent)
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

			for _, node := range scene.Nodes {
				renderNode(node, sc)
			}
			shaders.Stop()
		}
	})
}

/* TODO
// AddEntities adds entities to the Render System, this differs from the base addEntities function as this setups the entities to be batch rendered
func (rs *RenderSystem) AddEntities(entities ...*ecs.Entity) {

}
*/

// Name returns the name of the rendering system
func (rs *RenderSystem) Name() string {
	return "RenderSystem"
}

// Name returns the RenderComponent name
func (rc *RenderComponent) Name() string {
	return RenderComponentName
}

//NewCamera creates a new camera and attaches it to the renderer
func (rs *RenderSystem) NewCamera(width int, height int) {
	if rs.Camera == nil {
		rs.Camera = &Camera{
			Pos:   mgl32.Vec3{0, 0, 0},
			Yaw:   0,
			Pitch: 0,
			Roll:  0,
		}
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

func renderNode(n *Node, s shaders.ShaderComponent) {
	if n.Mesh != nil {
		prepareNode(n.Mesh, s)
		shaders.LoadMat(s, "transformMatrix", n.Transformation)
		gl.DrawElements(gl.TRIANGLES, int32(n.Mesh.VertexCount()), gl.UNSIGNED_INT, gl.Ptr(nil))
		unbindNode(n.Mesh)
	}

	if len(n.Children) > 0 {
		for _, node := range n.Children {
			renderNode(node, s)
		}
	}
}

func prepareNode(m *Mesh, s shaders.ShaderComponent) {
	gl.BindVertexArray(m.ID)
	for _, a := range m.Attributes {
		gl.EnableVertexAttribArray(uint32(a.Type))
	}

	if m.Material != nil {
		if m.Material.Transparent {
			disableCulling()
		}

		if m.Material.BaseColorTexture != nil {
			tc := m.Material.BaseColorTexture
			shaders.LoadBool(s, "useFakeLighting", tc.Transparant)
			shaders.LoadFloat(s, "shineDamper", tc.ShineDamper)
			shaders.LoadFloat(s, "reflectivity", tc.Reflectivity)
			gl.BindTexture(gl.TEXTURE_2D, tc.id)
		}
	}

}

func unbindNode(m *Mesh) {
	enableCulling()
	for _, a := range m.Attributes {
		gl.DisableVertexAttribArray(uint32(a.Type))
	}
	gl.BindVertexArray(0)
}
