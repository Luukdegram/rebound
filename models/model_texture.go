package models

//ModelTexture is a textured model
type ModelTexture struct {
	ID           uint32
	ShineDamper  float32
	Reflectivity float32
}

//NewModelTexture returns a new model texture
func NewModelTexture(id uint32) *ModelTexture {
	return &ModelTexture{ID: id, ShineDamper: 1, Reflectivity: 1}
}
