package rebound

import "github.com/go-gl/mathgl/mgl32"

//Camera handles the camera of the scene
type Camera struct {
	Pos              mgl32.Vec3
	Pitch, Yaw, Roll float32
}

//Move moves the camera around the 3D world given the input
func (c *Camera) Move(x, y, z float32) {
	c.Pos = c.Pos.Add([3]float32{x, y, z})
}

//MoveTo moves the camera to the 3D point
func (c *Camera) MoveTo(x, y, z float32) {
	c.Pos = [3]float32{x, y, z}
}
