package rebound

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/luukdegram/rebound/ecs"
	"github.com/luukdegram/rebound/internal/thread"
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
	drawPolygon bool
	Camera      *Camera
	Shader      Shader
	BaseColour  Colour
	Skybox      *Skybox
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
	*Mesh
	Rotation [3]float32
	Position [3]float32
	Scale    [3]float32
}

//NewRenderSystem returns a new RendererSystem with default settings
func NewRenderSystem() (*RenderSystem, error) {
	rs := &RenderSystem{
		BaseSystem:  ecs.NewBaseSystem(),
		drawPolygon: false,
		BaseColour:  Colour{0.1, 0.1, 0.1, 1},
	}

	shader, err := NewBasicShader()
	if err != nil {
		return nil, err
	}

	rs.Shader = shader
	return rs, nil
}

//Update draws all entities within a RendererSystem
func (rs *RenderSystem) Update(dt float64) {
	thread.Call(func() {
		rs.prepare()
		startShader(rs.Shader)
		rs.Shader.Setup(*rs.Camera)
		for _, e := range rs.BaseSystem.Entities() {
			rc := e.Component(RenderComponentName).(*RenderComponent)
			rs.Shader.Render(*rc)
			render(*rc)
		}
		stopShader()

		//as last, render our skybox
		rs.renderSkybox()
	})
}

// AddEntities adds entities to the Render System.
// TODO: This differs from the base addEntities function as this setups the entities to be batch rendered
func (rs *RenderSystem) AddEntities(entities ...*ecs.Entity) {
	for _, e := range entities {
		if e.HasComponent(RenderComponentName) {
			rs.BaseSystem.AddEntities(e)
		}

		if len(e.Children()) > 0 {
			rs.AddEntities(e.Children()...)
		}
	}
}

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
	var camera *Camera
	if rs.Camera == nil {
		camera = &Camera{
			Position:  [3]float32{0, 0, 0},
			Yaw:       0,
			Pitch:     0,
			Roll:      0,
			FOV:       105,
			NearPlane: 0.1,
			FarPlane:  100,
		}
	} else {
		camera = rs.Camera
	}

	camera.Projection = NewProjectionMatrix(camera.FOV, float32(width/height), camera.NearPlane, camera.FarPlane)
	rs.Camera = camera
}

//Prepare cleans the screen for the next draw
func (rs *RenderSystem) prepare() {
	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)
	gl.Enable(gl.DEPTH_TEST)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.ClearColor(rs.BaseColour.R, rs.BaseColour.G, rs.BaseColour.B, rs.BaseColour.A)
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

func render(rc RenderComponent) {
	// Bind the VAO
	gl.BindVertexArray(rc.ID)

	// Enable all the vertex attributes (position, texcoords, colors, normals, etc)
	for _, a := range rc.Attributes {
		gl.EnableVertexAttribArray(uint32(a.Type))
	}

	// If transparent, disable culling
	if rc.Material.Transparent {
		gl.Disable(gl.CULL_FACE)
	}

	// if texture exists, bind texture
	if rc.Material.BaseColorTexture != nil {
		gl.BindTexture(gl.TEXTURE_2D, *rc.Material.BaseColorTexture)
	}

	// Finally, draw the model
	gl.DrawElements(gl.TRIANGLES, int32(rc.VertexCount()), gl.UNSIGNED_INT, gl.Ptr(nil))

	// Cleanup, disable attributes and unbind vao
	for _, a := range rc.Attributes {
		gl.DisableVertexAttribArray(uint32(a.Type))
	}
	gl.BindVertexArray(0)
}

// renderSkybox renders a Skybox into the scene
func (rs *RenderSystem) renderSkybox() {
	if rs.Skybox == nil {
		return
	}
	sb := rs.Skybox
	gl.DepthFunc(gl.LEQUAL)
	startShader(sb.shader)
	sb.shader.Setup(*rs.Camera)
	gl.BindVertexArray(sb.mesh.ID)
	gl.EnableVertexAttribArray(0)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, sb.texture)
	gl.DrawArrays(gl.TRIANGLES, 0, 36)
	gl.DisableVertexAttribArray(0)
	gl.BindVertexArray(0)
	gl.DepthFunc(gl.LESS)
	stopShader()
}
