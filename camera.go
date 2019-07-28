package rebound

import "github.com/go-gl/mathgl/mgl32"

//Camera handles the camera of the scene
type Camera struct {
	Pos              mgl32.Vec3
	Pitch, Yaw, Roll float32
}

//NewCamera creates a new camera object
func NewCamera() *Camera {
	return &Camera{}
}

//Move moves the camera around the 3D world given the input
func (c *Camera) Move(vec mgl32.Vec3) {
	c.Pos = c.Pos.Add(vec)
}
