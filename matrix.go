package rebound

import (
	"github.com/go-gl/mathgl/mgl32"
)

//NewTransformationMatrix returns a new transformation matrix, it translates, rotates and scales.
func NewTransformationMatrix(trans mgl32.Vec3, rot mgl32.Vec3, scale float32) mgl32.Mat4 {
	mat := mgl32.Ident4()
	translation := mgl32.Translate3D(trans[0], trans[1], trans[2])
	rotX := mgl32.HomogRotate3DX(mgl32.DegToRad(rot[0]))
	rotY := mgl32.HomogRotate3DY(mgl32.DegToRad(rot[1]))
	rotZ := mgl32.HomogRotate3DZ(mgl32.DegToRad(rot[2]))
	scaleMatrix := mgl32.Scale3D(scale, scale, scale)

	return mat.Add(translation).Mul4(rotX).Mul4(rotY).Mul4(rotZ).Mul4(scaleMatrix)
}

//NewProjectionMatrix returns a new projection matrix
func NewProjectionMatrix(angle, aspect, nearPlane, farPlane float32) mgl32.Mat4 {
	return mgl32.Perspective(mgl32.DegToRad(angle), aspect, nearPlane, farPlane)
}

//NewViewMatrix returns a new view matrix
func NewViewMatrix(camera Camera) mgl32.Mat4 {
	mat := mgl32.Ident4()
	return mgl32.LookAtV(camera.Pos, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	rotX := mgl32.HomogRotate3DX(mgl32.DegToRad(camera.Pitch))
	rotY := mgl32.HomogRotate3DY(mgl32.DegToRad(camera.Yaw))
	translation := mgl32.Translate3D(-camera.Pos[0], -camera.Pos[1], -camera.Pos[2])
	return mat.Mul4(rotX).Mul4(rotY).Add(translation)
}
