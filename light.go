package rebound

import (
	"github.com/luukdegram/rebound/ecs"
)

var light *Light

//Light is an object that construct light in a scene
type Light struct {
	ecs.Entity
	Position [3]float32
	Colour   [3]float32
}
