package rebound

import (
	"github.com/luukdegram/rebound/ecs"
)

// Light represents a light within the scene, this can be positional or based on a direction
type Light struct {
	ecs.Entity
	Position  [3]float32
	Direction [3]float32
	Colour    [3]float32
	Specular  [3]float32
	Diffuse   [3]float32
	Ambient   [3]float32
}

// PointLight represents a point light in the scene
type PointLight struct {
	Light
	Constant  float32
	Linear    float32
	Quadratic float32
}
