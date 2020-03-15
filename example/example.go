package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/luukdegram/rebound"
	"github.com/luukdegram/rebound/ecs"
	"github.com/luukdegram/rebound/importers"
	"github.com/luukdegram/rebound/input"
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

type inputSystem struct {
	ecs.BaseSystem
	Camera *rebound.Camera
	Speed  float32
	rs     *rebound.RenderSystem // Solely for demo purposes, you probably wouldn't want this
}

func (is *inputSystem) Update(dt float64) {
	dist := float32(dt) * is.Speed
	if input.KeyW.Down() {
		is.Camera.Move(0, 0, -dist)
	}
	if input.KeyS.Down() {
		is.Camera.Move(0, 0, dist)
	}
	if input.KeyA.Down() {
		is.Camera.Move(-dist, 0, 0)
	}
	if input.KeyD.Down() {
		is.Camera.Move(dist, 0, 0)
	}
	if input.KeyP.Down() {
		is.rs.TogglePolygons()
	}
	if input.KeyQ.Down() {
		is.Camera.Move(0, dist, 0)
	}
	if input.KeyE.Down() {
		is.Camera.Move(0, -dist, 0)
	}
	if input.RightMouseButton.Down() {
		is.Camera.Yaw += (input.Cursor.Xoffset() * dist)
		is.Camera.Pitch += (input.Cursor.Yoffset() * dist)
	}
}

func (is *inputSystem) Name() string {
	return "InputSystem"
}

func setup() {
	gltfImporter := importers.GLTFImporter{}
	scene, err := gltfImporter.Import("gltf_objects/SciFiHelmet/glTF/SciFiHelmet.gltf")
	if err != nil {
		panic(err)
	}

	skybox, err := rebound.NewSkybox([6]string{
		"skybox/right.png",
		"skybox/left.png",
		"skybox/top.png",
		"skybox/bottom.png",
		"skybox/back.png",
		"skybox/front.png",
	})
	if err != nil {
		panic(err)
	}

	renderer, err := rebound.NewRenderSystem()
	if err != nil {
		panic(err)
	}
	renderer.AddEntities(scene)
	renderer.NewCamera(width, height)
	renderer.Camera.MoveTo(0, 0, 12.5)
	renderer.Skybox = skybox

	// Let's create some point lights for in the scene
	bs := renderer.Shader.(*rebound.BasicShader)
	size := 4
	counter := 0
	for counter < size {
		x := rand.Float32()*10 - 5
		y := rand.Float32()*10 - 5
		z := rand.Float32()*10 - 5
		bs.PointLights = append(bs.PointLights, rebound.PointLight{
			Light: rebound.Light{
				Position: [3]float32{x, y, z},
				Ambient:  [3]float32{0.05, 0.05, 0.05},
				Diffuse:  [3]float32{0.8, 0.8, 0.8},
				Specular: [3]float32{1, 1, 1},
			},
			Constant:  1,
			Linear:    0.09,
			Quadratic: 0.032,
		})
		counter++
	}

	inputSystem := &inputSystem{ecs.NewBaseSystem(), renderer.Camera, 150, renderer}
	ecs.GetManager().AddSystems(renderer, inputSystem)
}

func main() {
	options := rebound.RunOptions{
		Height: height,
		Width:  width,
		Title:  "Rebound Engine",
	}
	if err := rebound.Run(options, setup); err != nil {
		fmt.Println(err)
	}
}
