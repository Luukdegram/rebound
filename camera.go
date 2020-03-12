package rebound

import "github.com/go-gl/mathgl/mgl32"

//Camera handles the camera of the scene
type Camera struct {
	Pos        [3]float32
	ProjectMat [16]float32
	Pitch,
	Yaw,
	Roll,
	FOV,
	NearPlane,
	FarPlane float32
}

//Move moves the camera around the 3D world given the input
func (c *Camera) Move(x, y, z float32) {
	var pos mgl32.Vec3 = c.Pos
	c.Pos = pos.Add([3]float32{x, y, z})
}

//MoveTo moves the camera to the 3D point
func (c *Camera) MoveTo(x, y, z float32) {
	c.Pos = [3]float32{x, y, z}
}
