package main

import (
	"log"
	"os"

	"github.com/luukdegram/rebound"
	"github.com/luukdegram/rebound/display"
	"github.com/luukdegram/rebound/ecs"
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
	defer rebound.CleanUp()
	defer shaders.CleanUp()

	manager := ecs.GetManager()
	renderer := rebound.NewRenderSystem()
	rotater := &rotationSystem{ecs.NewBaseSystem()}
	manager.AddSystems(renderer, rotater)

	entity, err := importers.LoadGltfModel("gltf_objects/avacado.gltf")
	if err != nil {
		panic(err)
	}

	renderer.AddEntities(entity.Children()...)
	rotater.AddEntities(entity.Children()...)

	renderer.NewCamera(width, height)
	renderer.NewLight([3]float32{3000, 2000, 2000})
	renderer.SetSkyColor(0.5, 0.5, 0.5)

	renderer.Camera.MoveTo(0, 0.1, 1.5)

	window.RegisterKeyboardHandler(display.KeyP, func() {
		renderer.TogglePolygons()
	})

	window.RegisterScrollWheelHandler(func(x, y float64) {
		renderer.Camera.Move(0, 0, float32(-y*0.1))
	})

	for !window.ShouldClose() {
		manager.Update(1)
		window.Update()
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
