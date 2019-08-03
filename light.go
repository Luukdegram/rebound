package rebound

import "github.com/go-gl/mathgl/mgl32"

var light *Light

//Light is an object that construct light in a scene
type Light struct {
	Position mgl32.Vec3
	Colour   mgl32.Vec3
}
