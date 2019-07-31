package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/luukdegram/rebound"
	"github.com/luukdegram/rebound/display"
	"github.com/luukdegram/rebound/importers"
	"github.com/luukdegram/rebound/models"
	"github.com/luukdegram/rebound/shaders"
)

const (
	height = 480
	width  = 640
)

var (
	triangle = []float32{
		-0.5, 0.5, -0.5,
		-0.5, -0.5, -0.5,
		0.5, -0.5, -0.5,
		0.5, 0.5, -0.5,

		-0.5, 0.5, 0.5,
		-0.5, -0.5, 0.5,
		0.5, -0.5, 0.5,
		0.5, 0.5, 0.5,

		0.5, 0.5, -0.5,
		0.5, -0.5, -0.5,
		0.5, -0.5, 0.5,
		0.5, 0.5, 0.5,

		-0.5, 0.5, -0.5,
		-0.5, -0.5, -0.5,
		-0.5, -0.5, 0.5,
		-0.5, 0.5, 0.5,

		-0.5, 0.5, 0.5,
		-0.5, 0.5, -0.5,
		0.5, 0.5, -0.5,
		0.5, 0.5, 0.5,

		-0.5, -0.5, 0.5,
		-0.5, -0.5, -0.5,
		0.5, -0.5, -0.5,
		0.5, -0.5, 0.5,
	}
	indices = []uint32{
		0, 1, 3,
		3, 1, 2,
		4, 5, 7,
		7, 5, 6,
		8, 9, 11,
		11, 9, 10,
		12, 13, 15,
		15, 13, 14,
		16, 17, 19,
		19, 17, 18,
		20, 21, 23,
		23, 21, 22,
	}
	textureCoords = []float32{
		0, 0,
		0, 1,
		1, 1,
		1, 0,
		0, 0,
		0, 1,
		1, 1,
		1, 0,
		0, 0,
		0, 1,
		1, 1,
		1, 0,
		0, 0,
		0, 1,
		1, 1,
		1, 0,
		0, 0,
		0, 1,
		1, 1,
		1, 0,
		0, 0,
		0, 1,
		1, 1,
		1, 0,
	}
)

func init() {
	err := os.Chdir("assets")
	if err != nil {
		log.Panicln("os.Chdir:", err)
	}
}

func main() {
	window := display.Manager(display.Default())

	err := window.Init(width, height, "Rebound Engine")
	if err != nil {
		panic(err)
	}
	defer window.Close()

	meshes := importers.LoadGltfModel("gltf_objects/triangle.gltf")
	fmt.Println(len(meshes))

	shader, err := shaders.NewShaderProgram("shaders/vertexShader.vert", "shaders/fragmentShader.frag", "position", "textureCoords")
	if err != nil {
		panic(err)
	}
	renderer := rebound.NewRenderer()
	model := rebound.LoadToVAO(triangle, textureCoords, indices)
	texture, err := rebound.LoadTexture("textures/square.png")
	if err != nil {
		panic(err)
	}
	modelTexture := models.NewModelTexture(texture)
	texturedModel := models.NewTexturedModel(model, modelTexture)

	entity := rebound.NewEntity(texturedModel)

	camera := rebound.NewCamera()
	camera.Pos[2] = 3
	projection := rebound.NewProjectionMatrix(renderer.FOV, float32(width/height), renderer.NearPlane, renderer.FarPlane)

	window.RegisterKeyboardHandler(display.KeyP, func() {
		renderer.TogglePolygons()
	})

	for !window.ShouldClose() {
		entity.Rotate(mgl32.Vec3{0.1, 0.1, 0.1})

		transform := rebound.NewTransformationMatrix(entity.Position, entity.Rotation, entity.Scale)

		renderer.Prepare()
		shader.Start()
		shader.LoadMat(shader.GetUniformLocation("projectionMatrix"), projection)
		shader.LoadMat(shader.GetUniformLocation("viewMatrix"), rebound.NewViewMatrix(*camera))
		shader.LoadMat(shader.GetUniformLocation("transformMatrix"), transform)
		renderer.Render(entity)
		shader.Stop()

		window.Update()
	}
	shader.CleanUp()
	rebound.CleanUp()
}
