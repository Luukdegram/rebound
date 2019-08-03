package rebound

import "github.com/go-gl/mathgl/mgl32"

//Camera handles the camera of the scene
type Camera struct {
	Pos              mgl32.Vec3
	Pitch, Yaw, Roll float32
}

//Move moves the camera around the 3D world given the input
func (c *Camera) Move(vec mgl32.Vec3) {
	c.Pos = c.Pos.Add(vec)
}

//MoveTo moves the camera to the 3D point
func (c *Camera) MoveTo(vec mgl32.Vec3) {
	c.Pos = vec
}
