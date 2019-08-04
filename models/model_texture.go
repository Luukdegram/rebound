package models

//ModelTexture is a textured model
type ModelTexture struct {
	ID              uint32
	ShineDamper     float32
	Reflectivity    float32
	transparancy    bool
	useFakeLighting bool
}

//NewModelTexture returns a new model texture
func NewModelTexture(id uint32) *ModelTexture {
	return &ModelTexture{ID: id, ShineDamper: 1000, Reflectivity: 1, transparancy: false, useFakeLighting: false}
}

//HasTransparancy returns true if the texture has transparancy
func (m ModelTexture) HasTransparancy() bool {
	return m.transparancy
}

//SetTransparancy sets the transparancy of a texture
func (m *ModelTexture) SetTransparancy(val bool) {
	m.transparancy = val
}

//UseFakeLighting returns true if texture requires fake lighting to be used
func (m ModelTexture) UseFakeLighting() bool {
	return m.useFakeLighting
}

//SetFakeLighting sets the transparancy of a texture
func (m *ModelTexture) SetFakeLighting(val bool) {
	m.useFakeLighting = val
}
