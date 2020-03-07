package main

import (
	"fmt"
	"log"
	"os"

	"github.com/luukdegram/rebound"
	"github.com/luukdegram/rebound/ecs"
	"github.com/luukdegram/rebound/importers"
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
	entity, err := importers.LoadGltfModel("gltf_objects/avacado.gltf")
	if err != nil {
		panic(err)
	}

	renderer := rebound.NewRenderSystem()
	rotater := &rotationSystem{ecs.NewBaseSystem()}
	ecs.GetManager().AddSystems(renderer, rotater)

	renderer.AddEntities(entity.Children()...)
	rotater.AddEntities(entity.Children()...)

	renderer.NewCamera(width, height)
	renderer.NewLight([3]float32{3000, 2000, 2000})
	renderer.SetSkyColor(0.5, 0.5, 0.5)

	renderer.Camera.MoveTo(0, 0.1, 1.5)
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

type rotationSystem struct {
	ecs.BaseSystem
}

func (r *rotationSystem) Update(dt float32) {
	for _, e := range r.Entities() {
		rc := e.Component(rebound.RenderComponentName).(*rebound.RenderComponent)
		rc.Rotation[1] += 1.5
	}
}

func (r *rotationSystem) Name() string {
	return "RotationSystem"
}
