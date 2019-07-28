package rebound

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/luukdegram/rebound/models"
)

//Entity is a generic game object with a textured model, position, rotation and a scale
type Entity struct {
	model    *models.TexturedModel
	Position mgl32.Vec3
	Rotation mgl32.Vec3
	Scale    float32
}

//GetModel returns the model of an entity
func (e *Entity) GetModel() *models.TexturedModel {
	return e.model
}

//Rotate rotates an Entity with each coordinate being the amount of degrees
func (e *Entity) Rotate(rot mgl32.Vec3) {
	e.Rotation = e.Rotation.Add(rot)
}

//Trans translates an Entity (moves it accross the 3d world)
func (e *Entity) Trans(trans mgl32.Vec3) {
	e.Position = e.Position.Add(trans)
}

//NewEntity creates a new entity object
func NewEntity(model *models.TexturedModel) *Entity {
	e := Entity{model: model}
	e.Position = mgl32.Vec3{0, 0, 0}
	e.Rotation = mgl32.Vec3{0, 0, 0}
	e.Scale = 1.0
	return &e
}
