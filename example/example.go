package main

import (
	"log"
	"os"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/luukdegram/rebound"
	"github.com/luukdegram/rebound/display"
	"github.com/luukdegram/rebound/importers"
	"github.com/luukdegram/rebound/shaders"
)

const (
	height = 1080
	width  = 1920
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

	geo, err := importers.LoadGltfModel("gltf_objects/avacado.gltf")
	if err != nil {
		panic(err)
	}

	modelShader, err := shaders.NewShaderProgram("shaders/vertexShader.vert", "shaders/fragmentShader.frag")
	if err != nil {
		panic(err)
	}

	renderer := rebound.NewRenderer()
	renderer.NewCamera(width, height)
	renderer.NewLight(mgl32.Vec3{3000, 2000, 2000})
	renderer.SetSkyColor(0.5, 0.5, 0.5)

	entity := rebound.NewEntity()
	entity.Geometry = geo

	renderer.Camera.Pos[2] = 1.5
	renderer.Camera.Pos[1] = 0.1

	window.RegisterKeyboardHandler(display.KeyP, func() {
		renderer.TogglePolygons()
	})

	window.RegisterScrollWheelHandler(func(x, y float64) {
		renderer.Camera.Move(mgl32.Vec3{0, 0, float32(-y * 0.1)})
	})

	for !window.ShouldClose() {
		entity.Rotate(mgl32.Vec3{0, 1, 0})
		renderer.RegisterEntity(entity)
		renderer.Render(*modelShader)

		window.Update()
	}
	modelShader.CleanUp()
	rebound.CleanUp()
}
