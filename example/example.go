package main

import (
	"fmt"
	"log"
	"os"

	"github.com/luukdegram/rebound"
	"github.com/luukdegram/rebound/ecs"
	"github.com/luukdegram/rebound/importers"
	"github.com/luukdegram/rebound/input"
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

func setup() {
	gltfImporter := &importers.GLTFImporter{}
	scene, err := gltfImporter.Import("gltf_objects/SciFiHelmet/glTF/SciFiHelmet.gltf")
	if err != nil {
		panic(err)
	}

	shader, err := shaders.NewShaderComponent(shaders.VertexShader, shaders.FragmentShader)
	if err != nil {
		panic(err)
	}

	entity := ecs.NewEntity(scene, shader)

	renderer := rebound.NewRenderSystem()
	renderer.AddEntities(entity)
	renderer.NewCamera(width, height)
	renderer.NewLight([3]float32{300, 200, 200})
	renderer.SetSkyColor(0.5, 0.5, 0.5)
	renderer.Camera.MoveTo(0, 0.1, 10.5)

	rotater := &cameraSystem{ecs.NewBaseSystem(), renderer.Camera, 50}
	ecs.GetManager().AddSystems(renderer, rotater)
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

type cameraSystem struct {
	ecs.BaseSystem
	Camera *rebound.Camera
	Speed  float32
}

func (cs *cameraSystem) Update(dt float64) {
	dist := float32(dt) * cs.Speed
	if input.KeyW.Down() {
		cs.Camera.Move(0, dist, 0)
	}
	if input.KeyS.Down() {
		cs.Camera.Move(0, -dist, 0)
	}
	if input.KeyA.Down() {
		cs.Camera.Move(dist, 0, 0)
	}
	if input.KeyD.Down() {
		cs.Camera.Move(-dist, 0, 0)
	}
}

func (cs *cameraSystem) Name() string {
	return "RotationSystem"
}
