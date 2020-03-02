package rebound

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/luukdegram/rebound/ecs"
)

var light *Light

//Light is an object that construct light in a scene
type Light struct {
	ecs.Entity
	Position mgl32.Vec3
	Colour   mgl32.Vec3
}

